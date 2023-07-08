package users

import (
	"Twitter_like_application/internal/database/pg"
	"Twitter_like_application/internal/services"
	"database/sql"
	"encoding/json"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"regexp"
	"time"
)

func CreateUser(w http.ResponseWriter, r *http.Request) {
	newUser := &Users{}
	err := json.NewDecoder(r.Body).Decode(newUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if newUser.Name != "" {
		nameRegex := regexp.MustCompile(`^[a-zA-Z]+$`)
		if !nameRegex.MatchString(newUser.Name) {
			services.ReturnErr(w, "Invalid name format", http.StatusBadRequest)
			return
		}

		if len(newUser.Name) > 70 {
			services.ReturnErr(w, "Name exceeds maximum length", http.StatusBadRequest)
			return
		}
	}
	if newUser.Email != "" {
		emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
		if !emailRegex.MatchString(newUser.Email) {
			http.Error(w, "Invalid email format", http.StatusBadRequest)
			return
		}
		if len(newUser.Email) > 200 {
			services.ReturnErr(w, "Email exceeds maximum length", http.StatusBadRequest)
			return
		}
	}
	if newUser.Password != "" {
		nameRegex := regexp.MustCompile(`^[a-zA-Z]+$`)
		if !nameRegex.MatchString(newUser.Password) {
			services.ReturnErr(w, "Invalid password format", http.StatusBadRequest)
			return
		}

		if len(newUser.Password) > 100 {
			services.ReturnErr(w, "Password exceeds maximum length", http.StatusBadRequest)
			return
		}
	}
	if newUser.BirthDate != "" {
		_, err := time.Parse("2006-01-02", newUser.BirthDate)
		if err != nil {
			http.Error(w, "Invalid birth date format", http.StatusBadRequest)
			return
		}

		query := `SELECT id FROM users_tweeter WHERE email = $1`
		var existingUserID int
		err = pg.DB.QueryRow(query, newUser.Email).Scan(&existingUserID)
		if err == nil {
			services.ReturnErr(w, "User with this email already exists", http.StatusBadRequest)
			return
		} else if err != sql.ErrNoRows {
			services.ReturnErr(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if newUser.Name == "" || newUser.Email == "" || newUser.Password == "" || newUser.BirthDate == "" {
			services.ReturnErr(w, "Invalid user data", http.StatusBadRequest)
			return
		}
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
		if err != nil {
			services.ReturnErr(w, err.Error(), http.StatusInternalServerError)
			return
		}
		newUser.Password = string(hashedPassword)
		query = `INSERT INTO users_tweeter (name, password, email, nickname, location, bio, birthdate) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`
		err = pg.DB.QueryRow(query, newUser.Name, newUser.Password, newUser.Email, newUser.Nickname, newUser.Location, newUser.Bio, newUser.BirthDate).Scan(&newUser.ID)
		if err != nil {
			if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
				services.ReturnErr(w, "This user is already added", http.StatusBadRequest)
				return
			}
			services.ReturnErr(w, err.Error(), http.StatusInternalServerError)
			return
		}

		userToken := CheckEmail(newUser)
		newUser.EmailToken = userToken

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(newUser)
	}
}
