package users

import (
	"Twitter_like_application/internal/database/pg"
	"Twitter_like_application/internal/services"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
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
	var updatedProfile Users
	err := json.NewDecoder(r.Body).Decode(&updatedProfile)
	if err != nil {
		return fmt.Errorf("failed to decode request body: %v", err)
	}
	values := []any{}
	key := []string{}
	if updatedProfile.Name != "" {
		services.CheckName(w, updatedProfile.Name)
		values = append(values, updatedProfile.Name)
		key = append(key, " name = $"+strconv.Itoa(len(key)+1))
	}
	if updatedProfile.Password != "" {
		services.CheckPassword(w, updatedProfile.Password)
		hashedPassword := services.HashedPassword(w, updatedProfile.Password)
		values = append(values, string(hashedPassword))
		key = append(key, " password = $"+strconv.Itoa(len(key)+1))
	}

	if updatedProfile.Email != "" {
		services.CheckEmail(w, updatedProfile.Email)
		values = append(values, updatedProfile.Email)
		key = append(key, " email = $"+strconv.Itoa(len(key)+1))
	}
	if updatedProfile.BirthDate != "" {
		services.CheckBirthDate(w, updatedProfile.BirthDate)
		values = append(values, updatedProfile.BirthDate)
		key = append(key, " birthdate = $"+strconv.Itoa(len(key)+1))
	}
	if updatedProfile.Nickname != "" {
		services.CheckNickname(w, updatedProfile.Nickname)
		values = append(values, updatedProfile.Nickname)
		key = append(key, " nickname = $"+strconv.Itoa(len(key)+1))
	}
	if updatedProfile.Bio != "" {
		services.CheckBio(w, updatedProfile.Bio)
		values = append(values, updatedProfile.Bio)
		key = append(key, " bio = $"+strconv.Itoa(len(key)+1))
	}
	if updatedProfile.Location != "" {
		services.CheckLocation(w, updatedProfile.Location)
		values = append(values, updatedProfile.Location)
		key = append(key, " location = $"+strconv.Itoa(len(key)+1))

	}
	values = append(values, userID)
	keystring := strings.Join(key, ", ")
	query := fmt.Sprintf("UPDATE users_tweeter SET %s WHERE id = $%d", keystring, len(values))
	_, err = pg.DB.Exec(query, values...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}
	return nil
}
