package services

import (
	Serviceuser "Twitter_like_application/cmd/internal/database"
	"encoding/json"
	"net/http"
	"time"
)

var Tweets map[int]Serviceuser.Tweet

func GetTweets(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Tweets)
}

func CreateTweet(w http.ResponseWriter, r *http.Request) {
	var tweet Serviceuser.Tweet
	var user Serviceuser.Users
	_ = json.NewDecoder(r.Body).Decode(&tweet)
	//tweet.ID = time.Now().Format("20060102150405")
	tweet.CreatedAt = time.Now()
	Tweets[user.ID] = tweet
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(tweet)
}
func LikeTweet(w http.ResponseWriter, r *http.Request) {
	var tweet Serviceuser.Tweet
	_ = json.NewDecoder(r.Body).Decode(&tweet)
	tweet.Like++

}
