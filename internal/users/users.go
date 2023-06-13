package users

import (
	_ "Twitter_like_application/internal/database/pg"
	pg "Twitter_like_application/internal/database/pg"
	"Twitter_like_application/internal/services"
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"net/smtp"
	"net/url"
	"strconv"
	"time"
)

func handleAuthenticatedRequest(w http.ResponseWriter, r *http.Request, next http.Handler) {
	cookie, err := r.Cookie("session")
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)

	}

	sessionID := cookie.Value

	query := "SELECT id FROM users_tweeter WHERE logintoken = $1"
	var userID int
	err = pg.DB.QueryRow(query, sessionID).Scan(&userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

	}

	ctx := context.WithValue(r.Context(), "userID", userID)
	r = r.WithContext(ctx)

	next.ServeHTTP(w, r)
	return
}

func AuthHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handleAuthenticatedRequest(w, r, next)
	})
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	newUser := &Users{}
	err := json.NewDecoder(r.Body).Decode(newUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	query := `SELECT id FROM users_tweeter WHERE email = $1`
	var existingUserID int
	err = pg.DB.QueryRow(query, newUser.Email).Scan(&existingUserID)
	if err == nil {
		http.Error(w, "User with this email already exists", http.StatusBadRequest)
		return
	} else if err != sql.ErrNoRows {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if newUser.Name == "" || newUser.Email == "" || newUser.Password == "" {
		http.Error(w, "Invalid user data", http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	newUser.Password = string(hashedPassword)
	query = `INSERT INTO users_tweeter (name, password, email, nickname, location, bio, birthdate) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`
	err = pg.DB.QueryRow(query, newUser.Name, newUser.Password, newUser.Email, newUser.Nickname, newUser.Location, newUser.Bio, newUser.BirthDate).Scan(&newUser.ID)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			http.Error(w, "This user is already added", http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	userToken := CheckEmail(newUser)
	newUser.EmailToken = userToken

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newUser)
}

func LoginUsers(w http.ResponseWriter, r *http.Request) {
	user := &Users{}
	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	query := "SELECT password FROM users_tweeter WHERE email = $1"
	var savedPassword string
	err = pg.DB.QueryRow(query, user.Email).Scan(&savedPassword)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(savedPassword), []byte(user.Password))
	if err == nil {
		sessionID := uuid.New().String()

		cookie := &http.Cookie{
			Name:     "session",
			Value:    sessionID,
			Expires:  time.Now().AddDate(0, 0, 30),
			HttpOnly: true,
			Path:     "/",
		}
		http.SetCookie(w, cookie)

		updateQuery := "UPDATE users_tweeter SET logintoken = $1 WHERE email = $2"
		_, err = pg.DB.Exec(updateQuery, sessionID, user.Email)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response := map[string]interface{}{
			"status":  "success",
			"message": "Authentication successful",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	} else if err == bcrypt.ErrMismatchedHashAndPassword {
		response := map[string]interface{}{
			"status":  "error",
			"message": "Invalid email or password",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	} else {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func LogoutUser(w http.ResponseWriter, r *http.Request) {
	cookie := &http.Cookie{
		Name:    "session",
		Value:   "",
		Expires: time.Now().AddDate(0, 0, -1),
		Path:    "/",
	}
	http.SetCookie(w, cookie)

	response := map[string]interface{}{
		"status":  "success",
		"message": "Logged out successfully",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func ResetPassword(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)

	var userResPass Users
	err := json.NewDecoder(r.Body).Decode(&userResPass)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	query := "SELECT name, email  FROM users_tweeter WHERE id = $1"
	var user Users
	err = pg.DB.QueryRow(query, userID).Scan(&user.Name, &user.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	userResPass.Email = user.Email
	userResPass.Name = user.Name

	ResetPasswordPlusEmail(&userResPass)
}

func ResetPasswordPlusEmail(user *Users) {

	resetToken := services.GenerateResetToken()
	user.ResetPasswordToken = resetToken
	confirmURL := &url.URL{
		Scheme: "http",
		Host:   "test.com",
		Path:   "/reset-password",
		RawQuery: url.Values{
			"token": {resetToken},
		}.Encode(),
	}
	to := user.Email
	subject := "Reset your password"
	body := fmt.Sprintf("Dear %s,\n\nReset your password: click this link:\n%s", user.Name, confirmURL.String())

	auth := smtp.PlainAuth("", "your-email", "password", "your-site")
	err := smtp.SendMail("your-site:587", auth, "your-email", []string{to}, []byte(fmt.Sprintf("Subject: %s\n\n%s", subject, body)))
	if err != nil {
		return
	}
	return
}
func GetUserProfile(w http.ResponseWriter, r *http.Request) {
	userID := GetCurrentUserID(w, r)

	query := "SELECT id, name, email, nickname FROM users WHERE id = $1"
	row := pg.DB.QueryRow(query, userID)

	var userProfile Users
	err := row.Scan(&userProfile.ID, &userProfile.Name, &userProfile.Email, &userProfile.Nickname)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userProfile)
}

func GetCurrentUserID(w http.ResponseWriter, r *http.Request) int {
	var user Users
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return 0
	}
	return 1
}

func FollowUser(w http.ResponseWriter, r *http.Request) {
	var usersFollow UsersFollow
	err := json.NewDecoder(r.Body).Decode(&usersFollow)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	currentUserID, err := strconv.Atoi(usersFollow.ID1)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	targetUserID, err := strconv.Atoi(usersFollow.ID2)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var count int
	err = pg.DB.QueryRow("SELECT COUNT(*) FROM followers_subscriptions WHERE follower_id = $1 AND subscription_id = $2", currentUserID, targetUserID).Scan(&count)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if count > 0 {
		http.Error(w, "User is already subscribed to the target user", http.StatusBadRequest)
		return
	}

	_, err = pg.DB.Exec("INSERT INTO followers_subscriptions (follower_id, subscription_id) VALUES ($1, $2)", currentUserID, targetUserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Println("Done")
}

func UnfollowUser(w http.ResponseWriter, r *http.Request) {
	var usersFollow UsersFollow
	err := json.NewDecoder(r.Body).Decode(&usersFollow)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	currentUserID, err := strconv.Atoi(usersFollow.ID1)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	targetUserID, err := strconv.Atoi(usersFollow.ID2)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = pg.DB.Exec("DELETE FROM followers_subscriptions WHERE follower_id = $1 AND subscription_id = $2", currentUserID, targetUserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Println("Done")
}
func EditProfile(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int)

	var updatedProfile Users
	err := json.NewDecoder(r.Body).Decode(&updatedProfile)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	query := "UPDATE users_tweeter SET username = $1, nickname = $2, birthdate = $3, email = $4, password = $5 WHERE id = $6"
	values := []interface{}{updatedProfile.Name, updatedProfile.Nickname, updatedProfile.BirthDate, updatedProfile.Email, updatedProfile.Password, userID}

	if updatedProfile.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(updatedProfile.Password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		query += ", password = $5"
		values = append(values, hashedPassword)
	}

	query += " WHERE id = $6"
	values = append(values, userID)

	_, err = pg.DB.Exec(query, values...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"status":  "success",
		"message": "Profile updated successfully",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func GetFollowers(w http.ResponseWriter, r *http.Request) {
	userID := r.FormValue("user_id")
	if userID == "" {
		http.Error(w, "Missing user ID", http.StatusBadRequest)
		return
	}

	query := "SELECT u.id, u.username FROM users u INNER JOIN subscriptions s ON u.id = s.follower_id WHERE s.followee_id = $1"
	rows, err := pg.DB.Query(query, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var followers []Users

	for rows.Next() {
		var follower Users
		err := rows.Scan(&follower.UserID, &follower.Name)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		followers = append(followers, follower)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(followers)
}

func GetFollowing(w http.ResponseWriter, r *http.Request) {
	userID := r.FormValue("user_id")
	if userID == "" {
		http.Error(w, "Missing user ID", http.StatusBadRequest)
		return
	}

	query := "SELECT u.id, u.username FROM users u INNER JOIN subscriptions s ON u.id = s.followee_id WHERE s.follower_id = $1"
	rows, err := pg.DB.Query(query, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var following []Users

	for rows.Next() {
		var followee Users
		err := rows.Scan(&followee.UserID, &followee.Name)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		following = append(following, followee)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(following)
}

func SearchUsers(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query")
	if query == "" {
		http.Error(w, "Missing search query", http.StatusBadRequest)
		return
	}

	searchQuery := "%" + query + "%"
	query = "SELECT id, name, username FROM users WHERE name ILIKE $1 OR username ILIKE $1"

	rows, err := pg.DB.Query(query, searchQuery)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var users []Users

	for rows.Next() {
		var user Users
		err := rows.Scan(&user.ID, &user.Name, &user.Name)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func GetStatistics(w http.ResponseWriter, r *http.Request) {
	userCount, err := services.GetUserCount()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tweetCount, err := services.GetTweetCount()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	statistics := struct {
		UserCount  int `json:"user_count"`
		TweetCount int `json:"tweet_count"`
	}{
		UserCount:  userCount,
		TweetCount: tweetCount,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(statistics)
}
func CheckEmail(newUser *Users) string {
	token := make([]byte, 32)
	_, err := rand.Read(token)
	if err != nil {
		return ""
	}
	confirmToken := base64.URLEncoding.EncodeToString(token)

	confirmURL := &url.URL{
		Scheme: "http",
		Host:   "test.com",
		Path:   "/confirm-email",
		RawQuery: url.Values{
			"token": {confirmToken},
		}.Encode(),
	}
	to := newUser.Email
	subject := "Confirment your email"
	body := fmt.Sprintf("Confirment email: click this link:\n%s", confirmURL.String())

	auth := smtp.PlainAuth("", "your email", "password", "your site/token")

	err = smtp.SendMail("your email:587", auth, "your site/token", []string{to}, []byte(fmt.Sprintf("Subject: %s\n\n%s", subject, body)))
	if err != nil {
		return ""
	}

	return confirmToken
}

func GetCurrentProfile(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int)

	query := "SELECT id, name, email, birthdate,bio,location,nicname  FROM users_tweeter WHERE id = $1"
	var user Users
	err := pg.DB.QueryRow(query, userID).Scan(&user.ID, &user.Name, &user.Email, &user.BirthDate, &user.Bio, &user.Location, &user.Nickname)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
