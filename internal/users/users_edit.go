package users

import (
	"Twitter_like_application/internal/database/pg"
	"Twitter_like_application/internal/services"
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func EditProfile(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int)
	err := updateProfile(w, r, userID)
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
func updateProfile(w http.ResponseWriter, r *http.Request, userID int) error {
	nameRegex := regexp.MustCompile(`^[a-zA-Z]+$`)
	var updatedProfile Users
	err := json.NewDecoder(r.Body).Decode(&updatedProfile)
	if err != nil {
		return fmt.Errorf("failed to decode request body: %v", err)
	}
	values := []any{}
	key := []string{}
	if updatedProfile.Name != "" {
		if !nameRegex.MatchString(updatedProfile.Name) {
			services.ReturnErr(w, "Invalid name format", http.StatusBadRequest)
			return err
		}
		if len(updatedProfile.Name) > 70 {
			services.ReturnErr(w, "Name exceeds maximum length", http.StatusBadRequest)
			return err
		}
		values = append(values, updatedProfile.Name)
		key = append(key, " name = $"+strconv.Itoa(len(key)+1))
	}
	if updatedProfile.Password != "" {
		if !nameRegex.MatchString(updatedProfile.Password) {
			services.ReturnErr(w, "Invalid password format", http.StatusBadRequest)
			return err
		}
		if len(updatedProfile.Password) > 100 {
			services.ReturnErr(w, "Password exceeds maximum length", http.StatusBadRequest)
			return err
		}
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(updatedProfile.Password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return err
		}
		values = append(values, string(hashedPassword))
		key = append(key, " password = $"+strconv.Itoa(len(key)+1))
	}

	if updatedProfile.Email != "" {
		emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
		if !emailRegex.MatchString(updatedProfile.Email) {
			services.ReturnErr(w, "Invalid email format", http.StatusBadRequest)
			return err
		}
		if len(updatedProfile.Email) > 100 {
			services.ReturnErr(w, "Email exceeds maximum length", http.StatusBadRequest)
			return err
		}
		values = append(values, updatedProfile.Email)
		key = append(key, " email = $"+strconv.Itoa(len(key)+1))
	}
	if updatedProfile.BirthDate != "" {
		_, err := time.Parse("2006-01-02", updatedProfile.BirthDate)
		if err != nil {
			services.ReturnErr(w, "Invalid birth date format", http.StatusBadRequest)
			return err
		}
		values = append(values, updatedProfile.BirthDate)
		key = append(key, " birthdate = $"+strconv.Itoa(len(key)+1))
	}
	if updatedProfile.Nickname != "" {
		if !nameRegex.MatchString(updatedProfile.Nickname) {
			services.ReturnErr(w, "Invalid nickname format", http.StatusBadRequest)
			return err
		}
		if len(updatedProfile.Nickname) > 100 {
			services.ReturnErr(w, "Nickname exceeds maximum length", http.StatusBadRequest)
			return err
		}
		values = append(values, updatedProfile.Nickname)
		key = append(key, " nickname = $"+strconv.Itoa(len(key)+1))
	}
	if updatedProfile.Bio != "" {
		if !nameRegex.MatchString(updatedProfile.Bio) {
			services.ReturnErr(w, "Invalid bio format", http.StatusBadRequest)
			return err
		}
		if len(updatedProfile.Bio) > 400 {
			services.ReturnErr(w, "Bio exceeds maximum length", http.StatusBadRequest)
			return err
		}
		values = append(values, updatedProfile.Bio)
		key = append(key, " bio = $"+strconv.Itoa(len(key)+1))
	}
	if updatedProfile.Location != "" {
		if !nameRegex.MatchString(updatedProfile.Location) {
			services.ReturnErr(w, "Invalid location format", http.StatusBadRequest)
			return err
		}
		if len(updatedProfile.Location) > 200 {
			services.ReturnErr(w, "Location exceeds maximum length", http.StatusBadRequest)
			return err
		}
		values = append(values, updatedProfile.Location)
		key = append(key, " location = $"+strconv.Itoa(len(key)+1))

	}
	values = append(values, userID)
	keystring := strings.Join(key, ", ")
	query := fmt.Sprintf("UPDATE users_tweeter SET %s WHERE id = $%d", keystring, len(values))
	fmt.Println(query)
	fmt.Println(values)
	_, err = pg.DB.Exec(query, values...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}
	return nil
}
