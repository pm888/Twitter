package services

import (
	"Twitter_like_application/internal/database"
	"encoding/json"
	"net/http"
	"time"
)

var Tweets map[int]Serviceuser.Tweet

func GetTweetsLast100(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Tweets)
}

//func GeTweets100(tweets map[int]Serviceuser.Tweet) []Serviceuser.Tweet {
//	var tweetSlice []Serviceuser.Tweet
//	for _, tweet := range tweets {
//		tweetSlice = append(tweetSlice, tweet)
//	}
//
//	sort.Slice(tweetSlice, func(i, j int) bool {
//		return tweetSlice[i].CreatedAt.After(tweetSlice[j].CreatedAt)
//	})
//
//	if len(tweetSlice) > 100 {
//		tweetSlice = tweetSlice[:100]
//	}
//
//	return tweetSlice
//}

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

func RetweetTheTweets(w http.ResponseWriter, r *http.Request) {

}
