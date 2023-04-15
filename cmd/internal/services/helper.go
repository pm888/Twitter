package services

import (
	Serviceuser "Twitter_like_application/cmd/internal/database"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/smtp"
	"net/url"
)

func CheckEmai(newUser *Serviceuser.Users) string {
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

func ConfirmEmail(token string, user *Serviceuser.Users) error {
	for id, _ := range Serviceuser.UserDate {
		if user.ID == id || token == user.EmailTocken {
			user.ConfirmEmailToken = true
		}
	}

	return nil
}
