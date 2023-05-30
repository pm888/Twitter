package Tweets

import (
	Postgresql "Twitter_like_application/internal/database/postgresql"
	"Twitter_like_application/internal/services"
	Serviceuser "Twitter_like_application/internal/users"
	"encoding/json"
	"net/http"
	"time"
)

var s *Postgresql.ServicePostgresql

func CreateTweet(w http.ResponseWriter, r *http.Request) {
	var newTweet Serviceuser.Tweet
	err := json.NewDecoder(r.Body).Decode(&newTweet)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	query := "INSERT INTO tweets (user_id, content, created_at) VALUES ($1, $2, $3) RETURNING id"
	err = s.DB.QueryRow(query, newTweet.UserID, newTweet.Text, time.Now()).Scan(&newTweet.TweetID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newTweet)
}

func GetTweet(w http.ResponseWriter, r *http.Request) {
	tweetID := r.URL.Query().Get("tweet_id")
	if tweetID == "" {
		http.Error(w, "Missing tweet ID", http.StatusBadRequest)
		return
	}

	query := "SELECT id, user_id, content FROM tweets WHERE id = $1"
	var tweet Serviceuser.Tweet
	err := s.DB.QueryRow(query, tweetID).Scan(&tweet.TweetID, &tweet.UserID, &tweet.Text)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tweet)
}

func UpdateTweet(w http.ResponseWriter, r *http.Request) {
	tweetID := r.URL.Query().Get("tweet_id")
	newContent := r.FormValue("text")
	intId, err := services.ConvertStringToNumber(tweetID)
	if tweetID == "" {
		http.Error(w, "Missing tweet ID", http.StatusBadRequest)
		return
	}
	if newContent == "" {
		http.Error(w, "Missing new tweet content", http.StatusBadRequest)
		return
	}
	query := "UPDATE tweets SET content = $1 WHERE id = $2"
	result, err := s.DB.Exec(query, newContent, tweetID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "Tweet not found", http.StatusNotFound)
		return
	}

	var updatedTweet = Serviceuser.Tweet{TweetID: intId, Text: newContent}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedTweet)
}
func DeleteTweet(w http.ResponseWriter, r *http.Request) {
	tweetID := r.URL.Query().Get("tweet_id")
	if tweetID == "" {
		http.Error(w, "Missing tweet ID", http.StatusBadRequest)
		return
	}

	query := "DELETE FROM tweets WHERE id = $1"
	result, err := s.DB.Exec(query, tweetID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "Tweet not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func LikeTweet(w http.ResponseWriter, r *http.Request) {
	tweetID := r.FormValue("tweet_id")
	if tweetID == "" {
		http.Error(w, "Missing tweet ID", http.StatusBadRequest)
		return
	}

	userID, err := services.GetCurrentUserID(r)

	query := "SELECT COUNT(*) FROM likes WHERE user_id = $1 AND tweet_id = $2"
	var count int
	err = s.DB.QueryRow(query, userID, tweetID).Scan(&count)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if count > 0 {
		http.Error(w, "Tweet already liked", http.StatusBadRequest)
		return
	}

	query = "INSERT INTO likes (user_id, tweet_id) VALUES ($1, $2)"
	_, err = s.DB.Exec(query, userID, tweetID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
func UnlikeTweet(w http.ResponseWriter, r *http.Request) {
	tweetID := r.FormValue("tweet_id")
	if tweetID == "" {
		http.Error(w, "Missing tweet ID", http.StatusBadRequest)
		return
	}

	userID, err := services.GetCurrentUserID(r)

	query := "SELECT COUNT(*) FROM likes WHERE user_id = $1 AND tweet_id = $2"
	var count int
	err = s.DB.QueryRow(query, userID, tweetID).Scan(&count)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if count == 0 {
		http.Error(w, "Tweet not liked", http.StatusBadRequest)
		return
	}

	query = "DELETE FROM likes WHERE user_id = $1 AND tweet_id = $2"
	_, err = s.DB.Exec(query, userID, tweetID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
func Retweet(w http.ResponseWriter, r *http.Request) {
	tweetID := r.FormValue("tweet_id")
	if tweetID == "" {
		http.Error(w, "Missing tweet ID", http.StatusBadRequest)
		return
	}

	userID, err := services.GetCurrentUserID(r)

	query := "SELECT COUNT(*) FROM retweets WHERE user_id = $1 AND tweet_id = $2"
	var count int
	err = s.DB.QueryRow(query, userID, tweetID).Scan(&count)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if count > 0 {
		http.Error(w, "Tweet already retweeted", http.StatusBadRequest)
		return
	}

	query = "INSERT INTO retweets (user_id, tweet_id, created_at) VALUES ($1, $2, $3)"
	_, err = s.DB.Exec(query, userID, tweetID, time.Now())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
