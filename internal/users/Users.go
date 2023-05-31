package Serviceuser

import (
	Postgresql "Twitter_like_application/internal/database/postgresql"
	_ "Twitter_like_application/internal/database/postgresql"
	"Twitter_like_application/internal/services"
	"encoding/json"
	"fmt"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
	"html/template"
	"net/http"
	"strconv"
)

func CreateUser(w http.ResponseWriter, r *http.Request, s *Postgresql.ServicePostgresql) {
	var newUser Users
	if newUser.Name == "" || newUser.Email == "" || newUser.Password == "" || newUser.Nickname == "" {
		http.Error(w, "Invalid user data", http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	newUser.Password = string(hashedPassword)

	query := `INSERT INTO users (name, password, email, nickname) VALUES ($1, $2, $3, $4) RETURNING id`
	err = s.DB.QueryRow(query, newUser.Name, newUser.Password, newUser.Email, newUser.Nickname).Scan(&newUser.ID)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			http.Error(w, "This user is already added", http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	userToken := services.CheckEmail(&newUser)
	newUser.EmailToken = userToken

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newUser)
}

func LoginUsers(w http.ResponseWriter, r *http.Request, s *Postgresql.ServicePostgresql) {
	if r.Method == "POST" {
		usermail := r.FormValue("usermail")
		password := r.FormValue("password")

		query := "SELECT COUNT(*) FROM users WHERE email = $1 AND password = $2"
		var count int
		err := s.DB.QueryRow(query, usermail, password).Scan(&count)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if count > 0 {
			cookie := &http.Cookie{
				Name:  "session",
				Value: "authenticated",
			}
			http.SetCookie(w, cookie)
			http.Redirect(w, r, "/", http.StatusSeeOther)
		} else {
			t, _ := template.ParseFiles("login.html")
			t.Execute(w, nil)
		}
	}
}
func LogoutUser(w http.ResponseWriter, r *http.Request, s *Postgresql.ServicePostgresql) {
	sessionID := getSessionIDFromRequest(r)
	deleteSessionQuery := "DELETE FROM sessions WHERE session_id = $1"
	_, err := s.DB.Exec(deleteSessionQuery, sessionID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func getSessionIDFromRequest(r *http.Request) string {
	cookie, err := r.Cookie("session")
	if err != nil {
		return ""
	}
	return cookie.Value
}

func EditMyProfile(w http.ResponseWriter, r *http.Request, s *Postgresql.ServicePostgresql) {
	var newuser *ReplaceMyData

	hashedNewPassword, err := bcrypt.GenerateFromPassword([]byte(newuser.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = s.DB.Exec(`
		UPDATE UserData SET 
			Name = $1,
			Email = $2,
			Nickname = $3,
			BirthDate = $4,
			Bio = $5,
			Password = $6,
			Location = $7
		WHERE id = $8`,
		newuser.NewName,
		newuser.NewEmail,
		newuser.NewNickname,
		newuser.NewBirthDate,
		newuser.NewBio,
		string(hashedNewPassword),
		newuser.NewLocation,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
func ResetPassword(w http.ResponseWriter, r *http.Request) {
	var userResPass Users
	err := json.NewDecoder(r.Body).Decode(&userResPass)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	services.ResetPasswordPlusEmail(&userResPass)
}

func GetUserProfile(w http.ResponseWriter, r *http.Request, s *Postgresql.ServicePostgresql) {
	userID := GetCurrentUserID(w, r)

	query := "SELECT id, name, email, nickname FROM users WHERE id = $1"
	row := s.DB.QueryRow(query, userID)

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

func FollowUser(w http.ResponseWriter, r *http.Request, s *Postgresql.ServicePostgresql) {
	currentUserID, err := services.GetCurrentUserID(r)
	currentUserIDint, err := services.ConvertStringToNumber(currentUserID)

	targetUserID := r.FormValue("user_id")
	if targetUserID == "" {
		http.Error(w, "Missing target user ID", http.StatusBadRequest)
		return
	}
	targetUserIDint, err := services.ConvertStringToNumber(targetUserID)

	if !services.UserExists(targetUserID) {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	if services.IsUserFollowing(currentUserIDint, targetUserIDint) {
		http.Error(w, "Already following the user", http.StatusBadRequest)
		return
	}

	query := "INSERT INTO subscriptions (user_id, target_user_id) VALUES ($1, $2)"
	_, err = s.DB.Exec(query, currentUserID, targetUserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
func UnfollowUser(w http.ResponseWriter, r *http.Request, s *Postgresql.ServicePostgresql) {
	currentUserID, err := services.GetCurrentUserID(r)

	userID := r.FormValue("user_id")
	if userID == "" {
		http.Error(w, "Missing user ID", http.StatusBadRequest)
		return
	}

	query := "DELETE FROM subscriptions WHERE follower_id = $1 AND followee_id = $2"
	_, err = s.DB.Exec(query, currentUserID, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func GetFollowers(w http.ResponseWriter, r *http.Request, s *Postgresql.ServicePostgresql) {
	userID := r.FormValue("user_id")
	if userID == "" {
		http.Error(w, "Missing user ID", http.StatusBadRequest)
		return
	}

	query := "SELECT u.id, u.username FROM users u INNER JOIN subscriptions s ON u.id = s.follower_id WHERE s.followee_id = $1"
	rows, err := s.DB.Query(query, userID)
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

func GetFollowing(w http.ResponseWriter, r *http.Request, s *Postgresql.ServicePostgresql) {
	userID := r.FormValue("user_id")
	if userID == "" {
		http.Error(w, "Missing user ID", http.StatusBadRequest)
		return
	}

	query := "SELECT u.id, u.username FROM users u INNER JOIN subscriptions s ON u.id = s.followee_id WHERE s.follower_id = $1"
	rows, err := s.DB.Query(query, userID)
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

func SearchUsers(w http.ResponseWriter, r *http.Request, s *Postgresql.ServicePostgresql) {
	query := r.URL.Query().Get("query")
	if query == "" {
		http.Error(w, "Missing search query", http.StatusBadRequest)
		return
	}

	searchQuery := "%" + query + "%"
	query = "SELECT id, name, username FROM users WHERE name ILIKE $1 OR username ILIKE $1"

	rows, err := s.DB.Query(query, searchQuery)
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

func SearchTweets(w http.ResponseWriter, r *http.Request, s *Postgresql.ServicePostgresql) {
	query := r.URL.Query().Get("query")
	if query == "" {
		http.Error(w, "Missing search query", http.StatusBadRequest)
		return
	}

	searchQuery := "%" + query + "%"
	query = "SELECT id, user_id, content FROM tweets WHERE content ILIKE $1"

	rows, err := s.DB.Query(query, searchQuery)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var tweets []Tweet

	for rows.Next() {
		var tweet Tweet
		err := rows.Scan(&tweet.TweetID, &tweet.UserID, &tweet.Text)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tweets = append(tweets, tweet)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tweets)
}
func GetFollowingTweets(w http.ResponseWriter, r *http.Request, s *Postgresql.ServicePostgresql) {
	currentUserID := GetCurrentUserID(w, r)

	subscribedUserIDs, err := services.GetSubscribedUserIDs(currentUserID, s)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(subscribedUserIDs) == 0 {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]Tweet{})
		return
	}

	userIDStr := ""
	for i, userID := range subscribedUserIDs {
		if i > 0 {
			userIDStr += ","
		}
		userIDStr += strconv.Itoa(userID)
	}

	query := fmt.Sprintf("SELECT id, user_id, content FROM tweets WHERE user_id IN (%s)", userIDStr)

	rows, err := s.DB.Query(query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var tweets []Tweet

	for rows.Next() {
		var tweet Tweet
		err := rows.Scan(&tweet.TweetID, &tweet.UserID, &tweet.Text)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tweets = append(tweets, tweet)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tweets)
}
func GetStatistics(w http.ResponseWriter, r *http.Request, s *Postgresql.ServicePostgresql) {
	userCount, err := services.GetUserCount(s)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tweetCount, err := services.GetTweetCount(s)
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
