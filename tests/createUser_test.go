package createUser_test

import (
	Serviceuser "Twitter_like_application/internal/users"
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
		"nickname": "ivan1"
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
