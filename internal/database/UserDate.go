package Serviceuser

import (
	"Twitter_like_application/internal/services"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
	"html/template"
	"net/http"
)

type UserDataSTR struct {
	UserData map[int]*Users
}

var (
	UserData  = make(map[int]*Users)
	TwittData = make(map[int]*Tweet)
)

////for map
//func Put(u *Users) bool {
//	for _, user := range UserData {
//		if user.Email == u.Email {
//			return false
//
//		}
//
//	}
//	u.ID = len(UserData) + 1
//	UserData[u.ID] = u
//
//	return true
//
//}

// Creat user(PostgreSQL)
func CreateUser(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("postgres", "postgresql://username:password@localhost/dbname?sslmode=disable")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	var newUser Users
	err = json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if newUser.Name == "" || newUser.Email == "" || newUser.Password == "" || newUser.Nickname == "" {
		http.Error(w, "Invalid user data", http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	newUser.Password = string(hashedPassword)

	query := `INSERT INTO users (name, password, email, nickname) VALUES ($1, $2, $3, $4) RETURNING id`
	err = db.QueryRow(query, newUser.Name, newUser.Password, newUser.Email, newUser.Nickname).Scan(&newUser.ID)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			http.Error(w, "This user is already added", http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	userToken := services.CheckEmail(&newUser)
	newUser.EmailToken = userToken

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newUser)
}

// Creat user(map)
//func CreateUser(w http.ResponseWriter, r *http.Request) {
//	var newUser Users
//	err := json.NewDecoder(r.Body).Decode(&newUser)
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusBadRequest)
//		return
//	}
//
//	if newUser.Name == "" || newUser.Email == "" || newUser.Password == "" || newUser.Nickname == "" {
//		http.Error(w, "Invalid user data", http.StatusBadRequest)
//		return
//	}
//
//	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusInternalServerError)
//		return
//	}
//	newUser.Password = string(hashedPassword)
//	ret := Put(&newUser)
//	if ret == false {
//		fmt.Fprint(w, "This user is alredy added")
//		return
//	} else {
//		userToken := services.CheckEmail(&newUser)
//		newUser.EmailToken = userToken
//		w.WriteHeader(http.StatusCreated)
//		json.NewEncoder(w).Encode(newUser)
//	}
//	return
//}

func LoginUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		usermail := r.FormValue("usermail")
		password := r.FormValue("password")
		for _, name := range UserData {
			if name.Email == usermail || name.Password == password {
				cookie := &http.Cookie{
					Name:  "session",
					Value: "authenticated",
				}
				http.SetCookie(w, cookie)
				http.Redirect(w, r, "/", http.StatusSeeOther)
			}
			t, _ := template.ParseFiles("login.html")
			t.Execute(w, nil)
		}
	}
}

func LogoutUser(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session")
	if err != http.ErrNoCookie {
		cookie = &http.Cookie{
			Name:   "session",
			Value:  "",
			Path:   "/",
			MaxAge: -1,
		}
		http.SetCookie(w, cookie)
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	var deleteUser Users
	err := json.NewDecoder(r.Body).Decode(&deleteUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	//for map
	for id, _ := range UserData {
		if deleteUser.ID == id {
			delete(UserData, id)
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

}

func Following(w http.ResponseWriter, r *http.Request) {
	var user FollowingForUser
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	writer := user.Writer
	subscriber := user.Subscriber
	//for map
	for id, _ := range UserData {
		if id == subscriber {
			UserData[writer].Followers = append(UserData[writer].Following, subscriber)
			UserData[subscriber].Following = append(UserData[subscriber].Followers, writer)
		}
	}

}

func ExploreMyaccaunt(w http.ResponseWriter, r *http.Request) {
	var user Users
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	for id, _ := range UserData {
		if id == user.ID {
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, user)
		}
	}

}

func SeeMyTimeline() {
	//var user Users

}
func ExploreOtherUsers() {

}

func EditmyProfile(w http.ResponseWriter, r *http.Request) {
	var newuser ReplaceMyData
	err := json.NewDecoder(r.Body).Decode(&newuser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	hashedNewPassword, err := bcrypt.GenerateFromPassword([]byte(newuser.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	//for map
	for _, users := range UserData {
		if users.Name != newuser.NewName {
			users.Name = newuser.NewName
		} else if users.Email != newuser.NewEmail {
			users.Email = newuser.NewEmail
		} else if users.Nickname != newuser.NewNickname {
			users.Nickname = newuser.NewNickname
		} else if users.BirthDate != newuser.NewBirthDate {
			users.BirthDate = newuser.NewBirthDate
		} else if users.Bio != newuser.NewBio {
			users.Bio = newuser.NewBio
		} else if users.Password != string(hashedNewPassword) {
			users.Password = string(hashedNewPassword)
		} else if users.Location != newuser.NewLocation {
			users.Location = newuser.NewLocation
		}
		return

	}
}

func ResetPassword(w http.ResponseWriter, r *http.Request) {
	var userResPass Users
	err := json.NewDecoder(r.Body).Decode(&userResPass)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	services.ResetPasswordPlusEmail(&userResPass)
}

func Home(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session")
	if err != nil || cookie.Value != "authenticated" {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	t, _ := template.ParseFiles("home.html")
	t.Execute(w, nil)
}
