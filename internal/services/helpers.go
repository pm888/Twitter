package services

import "C"
import (
	Postgresql "Twitter_like_application/internal/database/pg"
	"golang.org/x/crypto/bcrypt"
	"regexp"
	"time"

	//Serviceuser "Twitter_like_application/internal/users"

	"bufio"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
)

type ErrResponse struct {
	Errtext string `json:"errtext"`
}
type CheckVal struct {
	CheckNameRegex     *regexp.Regexp
	CheckEmailRegex    *regexp.Regexp
	CheckBioRegex      *regexp.Regexp
	CheckPasswordRegex *regexp.Regexp
	CheckNicknameRegex *regexp.Regexp
}

var (
	CheckNameRegex     *regexp.Regexp
	CheckEmailRegex    *regexp.Regexp
	CheckBioRegex      *regexp.Regexp
	CheckPasswordRegex *regexp.Regexp
	CheckNicknameRegex *regexp.Regexp
)

func GenerateResetToken() string {
	const resetTokenLength = 32
	tokenBytes := make([]byte, resetTokenLength)
	_, err := rand.Read(tokenBytes)
	if err != nil {
		panic(err)
	}
	return hex.EncodeToString(tokenBytes)
}

func ConvertStringToNumber(str string) (int, error) {
	num, err := strconv.Atoi(str)
	if err != nil {
		return 0, err
	}
	return num, nil
}

func UserExists(userID string) bool {
	query := "SELECT EXISTS (SELECT 1 FROM users WHERE id = $1)"
	var exists bool
	err := Postgresql.DB.QueryRow(query, userID).Scan(&exists)
	if err != nil {
		return false
	}
	return exists
}

func IsUserFollowing(currentUserID, targetUserID int) bool {
	query := "SELECT EXISTS (SELECT 1 FROM subscriptions WHERE user_id = $1 AND target_user_id = $2)"
	var exists bool
	err := Postgresql.DB.QueryRow(query, currentUserID, targetUserID).Scan(&exists)
	if err != nil {
		return false
	}
	return exists
}

func GetCurrentUserID(r *http.Request) (string, error) {
	cookie, err := r.Cookie("session")
	if err == http.ErrNoCookie {
		return "", errors.New("No session cookie found")
	} else if err != nil {
		return "", err
	}

	userID, err := ExtractUserIDFromSessionCookie(cookie.Value)

	return userID, nil
}
func ExtractUserIDFromSessionCookie(cookieValue string) (string, error) {
	req, err := http.ReadRequest(bufio.NewReader(strings.NewReader("GET / HTTP/1.0\r\nCookie: session=" + cookieValue + "\r\n\r\n")))
	if err != nil {
		return "", err
	}

	cookie, err := req.Cookie("session")
	if err != nil {
		return "", err
	}

	userID := cookie.Value
	return userID, nil
}
func GetSubscribedUserIDs(userID string) ([]int, error) {
	query := "SELECT subscribed_user_id FROM subscriptions WHERE user_id = $1"
	rows, err := Postgresql.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subscribedUserIDs []int

	for rows.Next() {
		var subscribedUserID int
		err := rows.Scan(&subscribedUserID)
		if err != nil {
			return nil, err
		}
		subscribedUserIDs = append(subscribedUserIDs, subscribedUserID)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return subscribedUserIDs, nil
}
func GetUserCount() (int, error) {
	query := "SELECT COUNT(*) FROM users"
	var count int
	err := Postgresql.DB.QueryRow(query).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func GetTweetCount() (int, error) {
	query := "SELECT COUNT(*) FROM tweets"
	var count int
	err := Postgresql.DB.QueryRow(query).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}
func ReturnErr(w http.ResponseWriter, err string, code int) {
	var errj ErrResponse
	errj.Errtext = err
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(errj)
}
func Reg() {
	CheckNameRegex = regexp.MustCompile(`^[a-zA-Zа-яА-ЯёЁ\s.\-]+$`)
	CheckBioRegex = regexp.MustCompile(`^[a-zA-Zа-яА-ЯёЁ\s.\-]+$`)
	CheckEmailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	CheckPasswordRegex = regexp.MustCompile(`^[a-zA-Z0-9!@#$%^&*()-_=+{}\[\]|\\;:'",.<>/?]+$`)
	CheckNicknameRegex = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)

}

func CheckName(w http.ResponseWriter, name string) {
	if !CheckNameRegex.MatchString(name) {
		ReturnErr(w, "Invalid name format", http.StatusBadRequest)
		return
	}
	if len(name) > 100 {
		ReturnErr(w, "Name exceeds maximum length", http.StatusBadRequest)
		return
	}
}
func CheckEmail(w http.ResponseWriter, email string) {
	if !CheckEmailRegex.MatchString(email) {
		ReturnErr(w, "Invalid email format", http.StatusBadRequest)
		return
	}
	if len(email) > 100 {
		ReturnErr(w, "Email exceeds maximum length", http.StatusBadRequest)
		return
	}
}
func CheckBio(w http.ResponseWriter, bio string) {
	if !CheckBioRegex.MatchString(bio) {
		ReturnErr(w, "Invalid bio format", http.StatusBadRequest)
		return
	}
	if len(bio) > 400 {
		ReturnErr(w, "Bio exceeds maximum length", http.StatusBadRequest)
		return
	}
}
func CheckLocation(w http.ResponseWriter, location string) {
	if !CheckBioRegex.MatchString(location) {
		ReturnErr(w, "Invalid location format", http.StatusBadRequest)
		return
	}
	if len(location) > 100 {
		ReturnErr(w, "Location exceeds maximum length", http.StatusBadRequest)
		return
	}
}
func CheckNickname(w http.ResponseWriter, nickname string) {
	if !CheckNicknameRegex.MatchString(nickname) {
		ReturnErr(w, "Invalid nickname format", http.StatusBadRequest)
		return
	}
	if len(nickname) > 100 {
		ReturnErr(w, "Nikname maximum length", http.StatusBadRequest)
		return
	}
}
func CheckBirthDate(w http.ResponseWriter, birthdate string) {
	_, err := time.Parse("2006-01-02", birthdate)
	if err != nil {
		ReturnErr(w, "Invalid birth date format", http.StatusBadRequest)
		return
	}
}
func CheckPassword(w http.ResponseWriter, password string) {
	if !CheckPasswordRegex.MatchString(password) {
		ReturnErr(w, "Invalid password format", http.StatusBadRequest)
		return
	}
	if len(password) > 100 {
		ReturnErr(w, "Password maximum length 100", http.StatusBadRequest)
		return
	}
	if len(password) < 8 {
		ReturnErr(w, "Password minimum length 8", http.StatusBadRequest)
		return
	}
	if !containsUppercase(password) {
		ReturnErr(w, "Password must contain at least one uppercase letter", http.StatusBadRequest)
		return
	}
	if !containsSpecialCharacter(password) {
		ReturnErr(w, "Password must contain at least one special character", http.StatusBadRequest)
		return
	}
}

func containsUppercase(s string) bool {
	for _, char := range s {
		if char >= 'A' && char <= 'Z' {
			return true
		}
	}
	return false
}

func containsSpecialCharacter(s string) bool {
	specialCharacters := "!@#$%^&*()-_=+{}[]|\\;:'\",.<>/?"
	for _, char := range s {
		if strings.ContainsRune(specialCharacters, char) {
			return true
		}
	}
	return false
}
func HashedPassword(w http.ResponseWriter, password string) []byte {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	return hashedPassword
}
