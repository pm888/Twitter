package services

import (
	Postgresql "Twitter_like_application/internal/database/postgresql"
	Serviceuser "Twitter_like_application/internal/users"
	"bufio"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"net/smtp"
	"net/url"
	"strconv"
	"strings"
)

func CheckEmail(newUser *Serviceuser.Users) string {
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
func ResetPasswordPlusEmail(user *Serviceuser.Users) {
	resetToken := GenerateResetToken()
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
	body := fmt.Sprintf("Reset your password: click this link:\n%s", confirmURL.String())

	var auth = smtp.PlainAuth("", "your email", "password", "your site/token")
	err := smtp.SendMail("your email:587", auth, "your site/token", []string{to}, []byte(fmt.Sprintf("Subject: %s\n\n%s", subject, body)))
	if err != nil {
		return
	}
	return

}
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

func UserExists(userID string, s *Postgresql.ServicePostgresql) bool {
	query := "SELECT EXISTS (SELECT 1 FROM users WHERE id = $1)"
	var exists bool
	err := s.DB.QueryRow(query, userID).Scan(&exists)
	if err != nil {
		return false
	}
	return exists
}

func IsUserFollowing(currentUserID, targetUserID int, s *Postgresql.ServicePostgresql) bool {
	query := "SELECT EXISTS (SELECT 1 FROM subscriptions WHERE user_id = $1 AND target_user_id = $2)"
	var exists bool
	err := s.DB.QueryRow(query, currentUserID, targetUserID).Scan(&exists)
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
func GetSubscribedUserIDs(userID int, s *Postgresql.ServicePostgresql) ([]int, error) {
	query := "SELECT subscribed_user_id FROM subscriptions WHERE user_id = $1"
	rows, err := s.DB.Query(query, userID)
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
func GetUserCount(s *Postgresql.ServicePostgresql) (int, error) {
	query := "SELECT COUNT(*) FROM users"
	var count int
	err := s.DB.QueryRow(query).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func GetTweetCount(s *Postgresql.ServicePostgresql) (int, error) {
	query := "SELECT COUNT(*) FROM tweets"
	var count int
	err := s.DB.QueryRow(query).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}
