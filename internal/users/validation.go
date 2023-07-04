package users

import (
	"gopkg.in/go-playground/validator.v9"
	"net/mail"
	"regexp"
	"strings"
	"time"
	"unicode"
)

func HasDigit(password string) bool {
	for _, char := range password {
		if unicode.IsDigit(char) {
			return true
		}
	}
	return false
}

func HasCommonWord(password string) bool {
	commonWords := []string{"password", "12345678", "87654321", "qwerty123"}
	for _, word := range commonWords {
		if strings.Contains(password, word) {
			return true
		}
	}
	return false
}

func HasSequence(password string) bool {
	sequences := []string{"123", "abc", "xyz"}
	for _, sequence := range sequences {
		if strings.Contains(password, sequence) {
			return true
		}
	}
	return false
}
func CheckPassword(fl validator.FieldLevel, v *UserValid) bool {
	password := fl.Field().String()
	if len(password) < 8 {
		v.validErr["password"] += "short,"
		return false
	}
	if len(password) > 100 {
		v.validErr["password"] += "long,"
		return false
	}
	hasUpperCase := false
	hasSpecialChar := false
	hasDigit := false
	hasSequence := false
	hasCommonWord := false

	for _, char := range password {
		if unicode.IsUpper(char) {
			hasUpperCase = true
		} else if !unicode.IsLetter(char) && !unicode.IsNumber(char) {
			hasSpecialChar = true
		}
	}

	if !HasDigit(password) == false {
		hasDigit = true
	}

	if HasSequence(password) == false {
		hasSequence = true
	}

	if HasCommonWord(password) == false {
		hasCommonWord = true
	}
	if hasUpperCase == false {
		v.validErr["password"] += "uppercase,"
	}
	if hasSpecialChar == false {
		v.validErr["password"] += "special character,"
	}
	if hasDigit == false {
		v.validErr["password"] += "digit,"
	}
	if hasSequence == false {
		v.validErr["password"] += "sequence,"
	}
	if hasCommonWord == false {
		v.validErr["password"] += "common word,"
	}
	if (hasUpperCase && hasSpecialChar && hasDigit && hasSequence && hasCommonWord) == false {
		return false
	}

	return true
}

func CheckDataTime(fl validator.FieldLevel, v *UserValid) bool {
	dateStr := fl.Field().String()
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		v.validErr["datatime"] += "uncorrect date,"
		return false
	}
	currentDate := time.Now()
	if date.After(currentDate) {
		v.validErr["datatime"] += "invalid period more,"
		return true
	}

	return true
}

func CheckName(fl validator.FieldLevel, v *UserValid) bool {
	name := fl.Field().String()
	u := NameVal{}
	if len(name) > 100 {
		v.validErr["name"] += "long name,"
		u.long = true
	}
	pattern := "^[\\p{L}\\s]+$"
	match, _ := regexp.MatchString(pattern, name)
	if match == false {
		v.validErr["name"] += "digit or special character,"
		u.realName = true
	}
	if (u.long && u.realName) == false {
		return false
	}
	return true
}

func CheckNickName(fl validator.FieldLevel, v *UserValid) bool {
	nickname := fl.Field().String()
	if len(nickname) > 100 {
		v.validErr["nickname"] = "long"
	}
	return true

}
func CheckBio(fl validator.FieldLevel, v *UserValid) bool {
	nickname := fl.Field().String()
	if len(nickname) > 400 {
		v.validErr["nickname"] = "long"
	}
	return true

}
func CheckLocation(fl validator.FieldLevel, v *UserValid) bool {
	nickname := fl.Field().String()
	if len(nickname) > 100 {
		v.validErr["nickname"] = "long"
	}
	return true

}
func CheckEmailVal(fl validator.FieldLevel, v *UserValid) bool {
	email := fl.Field().String()
	_, err := mail.ParseAddress(email)
	if err != nil {
		v.validErr["email"] = "not correct email"
		return false
	}
	return true
}
