package server

import (
	Serviceuser "Twitter_like_application/cmd/internal/database"
	"github.com/gorilla/mux"
	"net/http"
)

var tweets []Serviceuser.Tweet

func Server() {
	r := mux.NewRouter()
	//r.HandleFunc("/tweets", getTweets).Methods("GET")
	r.HandleFunc("/tweets", Serviceuser.CreateUser).Methods("POST")
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
