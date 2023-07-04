package users

import (
	"Twitter_like_application/internal/database/pg"
	"Twitter_like_application/internal/services"
	"encoding/json"
	"fmt"
	"gopkg.in/go-playground/validator.v9"
	"net/http"
	"strconv"
	"strings"
)

func EditProfile(w http.ResponseWriter, r *http.Request) {
	userValid := &UserValid{
		validate: validator.New(),
		validErr: make(map[string]string),
	}

	err := userValid.validate.RegisterValidation("checkPassword", func(fl validator.FieldLevel) bool {
		return CheckPassword(fl, userValid)
	})
	if err != nil {
		services.ReturnErr(w, "err.Error()", http.StatusInternalServerError)
	}
	err = userValid.validate.RegisterValidation("checkName", func(fl validator.FieldLevel) bool {
		return CheckName(fl, userValid)
	})
	if err != nil {
		services.ReturnErr(w, "err.Error()", http.StatusInternalServerError)
	}

	err = userValid.validate.RegisterValidation("checkDataTime", func(fl validator.FieldLevel) bool {
		return CheckDataTime(fl, userValid)
	})
	if err != nil {
		services.ReturnErr(w, "err.Error()", http.StatusInternalServerError)
	}
	err = userValid.validate.RegisterValidation("checkNickname", func(fl validator.FieldLevel) bool {
		return CheckNickName(fl, userValid)
	})
	if err != nil {
		services.ReturnErr(w, "err.Error()", http.StatusInternalServerError)
	}
	err = userValid.validate.RegisterValidation("checkBio", func(fl validator.FieldLevel) bool {
		return CheckBio(fl, userValid)
	})
	if err != nil {
		services.ReturnErr(w, "err.Error()", http.StatusInternalServerError)
	}
	err = userValid.validate.RegisterValidation("checkLocation", func(fl validator.FieldLevel) bool {
		return CheckLocation(fl, userValid)
	})
	if err != nil {
		services.ReturnErr(w, "err.Error()", http.StatusInternalServerError)
	}
	err = userValid.validate.RegisterValidation("email", func(fl validator.FieldLevel) bool {
		return CheckEmailVal(fl, userValid)
	})
	if err != nil {
		services.ReturnErr(w, "err.Error()", http.StatusInternalServerError)
	}
	userID := r.Context().Value("userID").(int)
	err = updateProfile(r, userID, userValid)
	if err != nil {
		services.ReturnErr(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"status":  "success",
		"message": "Profile updated successfully",
	}
	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(response); err != nil {
		services.ReturnErr(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func updateProfile(r *http.Request, userID int, v *UserValid) error {
	var (
		hashedPassword []byte
		key            = []string{}
		updatedProfile = EditUser{}
		values         = []any{}
	)

	err := json.NewDecoder(r.Body).Decode(&updatedProfile)
	if err != nil {
		return fmt.Errorf("failed to decode request body: %v", err)
	}
	err = v.validate.Struct(updatedProfile)
	if err != nil {
		if validationErrs, ok := err.(validator.ValidationErrors); ok {
			for _, e := range validationErrs {
				v.failedValidations = append(v.failedValidations, e.Field())
			}
		}
		return v
	}
	if updatedProfile.Name != "" {
		values = append(values, updatedProfile.Name)
		key = append(key, " name = $"+strconv.Itoa(len(key)+1))
	}
	if updatedProfile.Password != "" {
		err, hashedPassword = services.HashedPassword(updatedProfile.Password)
		if err != nil {
			return err
		}
		values = append(values, string(hashedPassword))
		key = append(key, " password = $"+strconv.Itoa(len(key)+1))
	}
	if updatedProfile.Email != "" {
		values = append(values, updatedProfile.Email)
		key = append(key, " email = $"+strconv.Itoa(len(key)+1))
	}
	if updatedProfile.BirthDate != "" {
		values = append(values, updatedProfile.BirthDate)
		key = append(key, " birthdate = $"+strconv.Itoa(len(key)+1))
	}

	if updatedProfile.Nickname != "" {
		values = append(values, updatedProfile.Nickname)
		key = append(key, " nickname = $"+strconv.Itoa(len(key)+1))
	}
	if updatedProfile.Bio != "" {
		values = append(values, updatedProfile.Bio)
		key = append(key, " bio = $"+strconv.Itoa(len(key)+1))
	}

	if updatedProfile.Location != "" {
		values = append(values, updatedProfile.Location)
		key = append(key, " location = $"+strconv.Itoa(len(key)+1))

	}
	values = append(values, userID)
	keyString := strings.Join(key, ", ")
	query := fmt.Sprintf("UPDATE users_tweeter SET %s WHERE id = $%d", keyString, len(values))
	_, err = pg.DB.Exec(query, values...)
	if err != nil {
		return err
	}
	return err
}
