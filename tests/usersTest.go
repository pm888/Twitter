package usersTest

import (
	Serviceuser "Twitter_like_application/internal/users"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {
	w := httptest.NewRecorder()
	reqBody := strings.NewReader(`{
		"name": "Ivan Ivanov",
		"email": "ivanov@example.com",
		"password": "password123",
		"nickname": "ivan1",
        "birthdate":"01.01.2021",
        "nickname":"nickname",
	    "bio":"bio",
        "location":"testlocation"
	}`)
	r, err := http.NewRequest("POST", "/create-user", reqBody)
	assert.NoError(t, err)
	Serviceuser.CreateUser(w, r)

	assert.Equal(t, http.StatusCreated, w.Code)

	var newUser Serviceuser.Users
	err = json.Unmarshal(w.Body.Bytes(), &newUser)
	assert.NoError(t, err)

	assert.NotEmpty(t, newUser.ID)
	assert.NotEmpty(t, newUser.EmailToken)
}

func TestLoginUsers(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(Serviceuser.LoginUsers))
	defer server.Close()

	formData := strings.NewReader("usermail=test@example.com&password=123456")

	resp, err := http.Post(server.URL, "application/x-www-form-urlencoded", formData)
	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusSeeOther {
		t.Errorf("Expected status code %d, but got %d", http.StatusSeeOther, resp.StatusCode)
	}

	cookies := resp.Cookies()
	if len(cookies) == 0 || cookies[0].Name != "session" || cookies[0].Value != "authenticated" {
		t.Errorf("Expected cookie 'session' with value 'authenticated', but got %+v", cookies)
	}
}

func TestLogoutUser(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(Serviceuser.LogoutUser))
	defer server.Close()

	req, err := http.NewRequest("GET", server.URL, nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusSeeOther {
		t.Errorf("Expected status code %d, but got %d", http.StatusSeeOther, resp.StatusCode)
	}

}

func TestResetPassword(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(Serviceuser.ResetPassword))
	defer server.Close()

	user := Serviceuser.Users{
		Email: "test@example.com",
	}

	payload, err := json.Marshal(user)
	if err != nil {
		t.Fatalf("Failed to marshal JSON payload: %v", err)
	}

	req, err := http.NewRequest("POST", server.URL, bytes.NewBuffer(payload))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, resp.StatusCode)
	}
}

func TestGetUserProfile(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(Serviceuser.GetUserProfile))
	defer server.Close()

	var (
		user = Serviceuser.Users{
			ID:        888,
			Name:      "Test User",
			Email:     "test@example.com",
			BirthDate: "01.01.2021",
			Nickname:  "test_nickname",
			Bio:       "testBio",
			Location:  "Test_location",
		}
	)

	req, err := http.NewRequest("GET", server.URL, nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	var userProfile Serviceuser.Users

	req.Header.Set("X-UserID", string(userProfile.ID))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, resp.StatusCode)
	}

	err = json.NewDecoder(resp.Body).Decode(&userProfile)
	if err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}

	if userProfile.ID != user.ID {
		t.Errorf("Expected user ID %s, but got %s", user.ID, userProfile.ID)
	}
	if userProfile.Name != user.Name {
		t.Errorf("Expected user name %s, but got %s", user.Name, userProfile.Name)
	}
	if userProfile.Email != user.Email {
		t.Errorf("Expected user email %s, but got %s", user.Email, userProfile.Email)
	}
	if userProfile.Nickname != user.Nickname {
		t.Errorf("Expected user nickname %s, but got %s", user.Nickname, userProfile.Nickname)
	}
}

func TestFollowUser(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(Serviceuser.FollowUser))
	defer server.Close()

	req, err := http.NewRequest("POST", server.URL, nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	q := req.URL.Query()
	q.Add("user_id", "targetUserID")
	req.URL.RawQuery = q.Encode()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, resp.StatusCode)
	}

}
