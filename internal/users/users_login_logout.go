package users

import (
	"Twitter_like_application/internal/database/pg"
	"Twitter_like_application/internal/services"
	"encoding/json"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

func LoginUsers(w http.ResponseWriter, r *http.Request) {
	user := &Users{}
	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		services.ReturnErr(w, err.Error(), http.StatusBadRequest)
		return
	}
	if user.Email != "" {
		services.CheckEmail(w, user.Email)
	}

	if user.Password != "" {
		services.CheckPassword(w, user.Password)
	}
	query := "SELECT id, password FROM users_tweeter WHERE email = $1"
	var userID int
	var savedPassword string
	err = pg.DB.QueryRow(query, user.Email).Scan(&userID, &savedPassword)
	if err != nil {
		services.ReturnErr(w, err.Error(), http.StatusInternalServerError)
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

		insertQuery := "INSERT INTO user_session (user_id, login_token, timestamp) VALUES ($1, $2, $3)"
		_, err = pg.DB.Exec(insertQuery, userID, cookie.Value, time.Now())
		if err != nil {
			services.ReturnErr(w, err.Error(), http.StatusInternalServerError)
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
		services.ReturnErr(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func LogoutUser(w http.ResponseWriter, r *http.Request) {
	apikey := r.Header.Get("X-API-KEY")
	if apikey == "" {
		cookie, err := r.Cookie("session")
		if err != nil {
			services.ReturnErr(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = DeleteUserSession(cookie.Value)
		if err != nil {
			services.ReturnErr(w, err.Error(), http.StatusInternalServerError)
			return
		}

		cookie = &http.Cookie{
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
	} else {
		err := DeleteUserSession(apikey)
		if err != nil {
			services.ReturnErr(w, err.Error(), http.StatusInternalServerError)
			return
		}
		response := map[string]interface{}{
			"status":  "success",
			"message": "Logged out successfully",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

func DeleteUserSession(token string) error {
	query := "DELETE FROM user_session WHERE login_token = $1"
	_, err := pg.DB.Exec(query, token)
	if err != nil {
		return err
	}

	return nil
}
