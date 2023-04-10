package Serviceuser

import (
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

type UserDataSTR struct {
	UserData map[int]*Users
}

var (
	UserDate = make(map[int]*Users)
)

func Put(u *Users) bool {
	for _, user := range UserDate {
		if user.Email == u.Email {
			return false

		}

	}
	u.ID = len(UserDate) + 1
	UserDate[u.ID] = u

	return true

}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	var newUser Users
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if newUser.Name == "" || newUser.Email == "" || newUser.Password == "" {
		http.Error(w, "Invalid user data", http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	newUser.Password = string(hashedPassword)
	ret := Put(&newUser)
	if ret == false {
		fmt.Fprint(w, "This user is alredy added")
		return
	} else {
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(newUser)
	}
	return
}
