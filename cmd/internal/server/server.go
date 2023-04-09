package server

import (
	Serviceuser "Twitter_like_application/cmd/internal/database"
	"encoding/json"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

var tweets []Serviceuser.Tweet

func Server() {
	r := mux.NewRouter()
	//r.HandleFunc("/tweets", getTweets).Methods("GET")
	r.HandleFunc("/tweets", createUser).Methods("POST")
	http.ListenAndServe("localhost:8080", r)
}

//func getTweets(w http.ResponseWriter, r *http.Request) {
//	w.Header().Set("Content-Type", "application/json")
//	json.NewEncoder(w).Encode(tweets)
//}

//func createTweet(w http.ResponseWriter, r *http.Request) {
//	var tweet Serviceuser.Tweet
//	_ = json.NewDecoder(r.Body).Decode(&tweet)
//	tweet.ID = time.Now().Format("20060102150405")
//	tweet.CreatedAt = time.Now()
//	tweets = append(tweets, tweet)
//	w.Header().Set("Content-Type", "application/json")
//	w.WriteHeader(http.StatusCreated)
//	json.NewEncoder(w).Encode(tweet)
//}

func createUser(w http.ResponseWriter, r *http.Request) {
	//for i, _ := range Serviceuser.UserDate {
	//	fmt.Println(i)
	//
	//}
	var newUser Serviceuser.Users
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
	Serviceuser.Put(&newUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newUser)
}
