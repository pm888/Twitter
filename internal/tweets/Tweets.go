package Tweets

import (
	Postgresql "Twitter_like_application/internal/database/postgresql"
	"Twitter_like_application/internal/services"
	Serviceuser "Twitter_like_application/internal/users"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

var db *Postgresql.ServicePostgresql

func CreateTweet(w http.ResponseWriter, r *http.Request) {
	var newTweet Serviceuser.Tweet
	err := json.NewDecoder(r.Body).Decode(&newTweet)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	query := "INSERT INTO tweets (user_id, content, created_at) VALUES ($1, $2, $3) RETURNING id"
	err = db.DB.QueryRow(query, newTweet.UserID, newTweet.Text, time.Now()).Scan(&newTweet.TweetID)
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
	err := db.DB.QueryRow(query, tweetID).Scan(&tweet.TweetID, &tweet.UserID, &tweet.Text)
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
	result, err := db.DB.Exec(query, newContent, tweetID)
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
	result, err := db.DB.Exec(query, tweetID)
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
	err = db.DB.QueryRow(query, userID, tweetID).Scan(&count)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if count > 0 {
		http.Error(w, "Tweet already liked", http.StatusBadRequest)
		return
	}

	query = "INSERT INTO likes (user_id, tweet_id) VALUES ($1, $2)"
	_, err = db.DB.Exec(query, userID, tweetID)
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
	err = db.DB.QueryRow(query, userID, tweetID).Scan(&count)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if count == 0 {
		http.Error(w, "Tweet not liked", http.StatusBadRequest)
		return
	}

	query = "DELETE FROM likes WHERE user_id = $1 AND tweet_id = $2"
	_, err = db.DB.Exec(query, userID, tweetID)
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
	err = db.DB.QueryRow(query, userID, tweetID).Scan(&count)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if count > 0 {
		http.Error(w, "Tweet already retweeted", http.StatusBadRequest)
		return
	}

	query = "INSERT INTO retweets (user_id, tweet_id, created_at) VALUES ($1, $2, $3)"
	_, err = db.DB.Exec(query, userID, tweetID, time.Now())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
func GetPopularTweets(w http.ResponseWriter, r *http.Request) {
	query := "SELECT id, user_id, content FROM tweets ORDER BY likes DESC LIMIT 10"

	rows, err := db.DB.Query(query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var tweets []Serviceuser.Tweet

	for rows.Next() {
		var tweet Serviceuser.Tweet
		err := rows.Scan(&tweet.TweetID, &tweet.UserID, &tweet.Text)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tweets = append(tweets, tweet)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tweets)
}
func SearchTweets(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query")
	if query == "" {
		http.Error(w, "Missing search query", http.StatusBadRequest)
		return
	}

	searchQuery := "%" + query + "%"
	query = "SELECT id, user_id, content FROM tweets WHERE content ILIKE $1"

	rows, err := db.DB.Query(query, searchQuery)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var tweets []Serviceuser.Tweet

	for rows.Next() {
		var tweet Serviceuser.Tweet
		err := rows.Scan(&tweet.TweetID, &tweet.UserID, &tweet.Text)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tweets = append(tweets, tweet)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tweets)
}
func GetFollowingTweets(w http.ResponseWriter, r *http.Request) {
	currentUserID, err := services.GetCurrentUserID(r)

	subscribedUserIDs, err := services.GetSubscribedUserIDs(currentUserID, db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(subscribedUserIDs) == 0 {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]Serviceuser.Tweet{})
		return
	}

	userIDStr := ""
	for i, userID := range subscribedUserIDs {
		if i > 0 {
			userIDStr += ","
		}
		userIDStr += strconv.Itoa(userID)
	}

	query := fmt.Sprintf("SELECT id, user_id, content FROM tweets WHERE user_id IN (%s)", userIDStr)

	rows, err := db.DB.Query(query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var tweets []Serviceuser.Tweet

	for rows.Next() {
		var tweet Serviceuser.Tweet
		err := rows.Scan(&tweet.TweetID, &tweet.UserID, &tweet.Text)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tweets = append(tweets, tweet)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tweets)
}
