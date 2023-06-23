package tweets

import (
	"Twitter_like_application/internal/database/pg"
	_ "Twitter_like_application/internal/database/pg"
	"Twitter_like_application/internal/services"
	Serviceuser "Twitter_like_application/internal/users"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"time"
)

func CreateTweet(w http.ResponseWriter, r *http.Request) {
	var newTweet Serviceuser.Tweet
	err := json.NewDecoder(r.Body).Decode(&newTweet)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	cookie, err := r.Cookie("session")
	if err != nil || cookie.Value == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	tokenString := cookie.Value

	query := "SELECT logintoken FROM users_tweeter WHERE id = $1"
	stmt, err := pg.DB.PrepareContext(r.Context(), query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	var dbToken string
	err = stmt.QueryRowContext(r.Context(), newTweet.UserID).Scan(&dbToken)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if tokenString != dbToken {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	query = `INSERT INTO tweets (user_id, author, text, created_at, like_count, repost, public, only_followers, only_mutual_followers, only_me)
	 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING tweet_id`
	err = pg.DB.QueryRowContext(r.Context(), query, newTweet.UserID, newTweet.Author, newTweet.Text, time.Now(), newTweet.LikeCount, newTweet.Repost, newTweet.Public, newTweet.OnlyFollowers, newTweet.OnlyMutualFollowers, newTweet.OnlyMe).Scan(&newTweet.TweetID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func GetTweet(w http.ResponseWriter, r *http.Request) {
	tweetID := r.URL.Query().Get("tweet_id")
	if tweetID == "" {
		http.Error(w, "Missing tweet ID", http.StatusBadRequest)
		return
	}

	query := "SELECT id, user_id, content FROM tweets WHERE id = $1"
	var tweet Serviceuser.Tweet
	err := pg.DB.QueryRow(query, tweetID).Scan(&tweet.TweetID, &tweet.UserID, &tweet.Text)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tweet)
}
func EditTweet(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int)

	var updatedTweet Serviceuser.Tweet
	err := json.NewDecoder(r.Body).Decode(&updatedTweet)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	query := "SELECT user_id FROM tweets WHERE tweet_id = $1"
	var tweetUserID int
	err = pg.DB.QueryRow(query, updatedTweet.TweetID).Scan(&tweetUserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if userID != tweetUserID {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	query = "UPDATE tweets SET text = $1 WHERE tweet_id = $2"
	_, err = pg.DB.Exec(query, updatedTweet.Text, updatedTweet.TweetID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"status":  "success",
		"message": "Tweet updated successfully",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func LikeTweet(w http.ResponseWriter, r *http.Request) {
	var like Serviceuser.Tweeter_like
	err := json.NewDecoder(r.Body).Decode(&like)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	query := "SELECT COUNT(*) FROM tweets WHERE user_id = $1 AND tweet_id = $2"
	var count int
	err = pg.DB.QueryRow(query, like.Autor, like.Id_post).Scan(&count)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if count > 0 {
		http.Error(w, "Tweet already liked", http.StatusBadRequest)
		return
	}

	query = "INSERT INTO tweets (autor, id_post, whose_like) VALUES ($1, $2, $3)"
	_, err = pg.DB.Exec(query, like.Autor, like.Id_post, like.Whose_like)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func UnlikeTweet(w http.ResponseWriter, r *http.Request) {
	idTweet := mux.Vars(r)["id_tweet"]
	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	query := "DELETE FROM likes WHERE user_id = $1 AND tweet_id = $2 RETURNING true"
	var exists bool
	err := pg.DB.QueryRow(query, userID, idTweet).Scan(&exists)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Tweet not liked", http.StatusBadRequest)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

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
	err = pg.DB.QueryRow(query, userID, tweetID).Scan(&count)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Tweet not liked", http.StatusBadRequest)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

}

func Retweet(w http.ResponseWriter, r *http.Request) {
	tweetID := mux.Vars(r)["id_tweet"]
	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	query := `
		SELECT EXISTS (
			SELECT 1
			FROM retweets
			WHERE tweet_id = $1 AND user_id = $2
			LIMIT 1
		), t.text
		FROM tweets t
		WHERE t.tweet_id = $1
		LIMIT 1
	`
	var exists bool
	var tweetText string
	err := pg.DB.QueryRow(query, tweetID, userID).Scan(&exists, &tweetText)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if tweetText == "" {
		http.Error(w, "Tweet not found", http.StatusNotFound)
		return
	}

	query = "INSERT INTO retweets (tweet_id, user_id, timestamp) VALUES ($1, $2, $3)"
	_, err = pg.DB.Exec(query, tweetID, userID, time.Now())
	if err != nil {
		tx.Rollback()
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	query = `
		INSERT INTO tweets (user_id, text, created_at, visibility, retweet)
		SELECT $1, $2, $3, visibility, $4
		FROM tweets
		WHERE tweet_id = $4
		LIMIT 1
	`
	_, err = pg.DB.Exec(query, userID, tweetText, time.Now(), tweetID)

	if tweetText == "" {
		http.Error(w, "Tweet not found", http.StatusNotFound)
		return
	}

	query = "INSERT INTO retweets (tweet_id, user_id, timestamp) VALUES ($1, $2, $3)"
	_, err = pg.DB.Exec(query, tweetID, userID, time.Now())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	public := true
	onlyFollowers := false
	onlyMutualFollowers := false
	onlyMe := false

	query = "INSERT INTO tweets (user_id, text, created_at, public, only_followers, only_mutual_followers, only_me, retweet) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)"
	_, err = pg.DB.Exec(query, userID, fmt.Sprintf(tweetText), time.Now(), public, onlyFollowers, onlyMutualFollowers, onlyMe, tweetID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func GetPopularTweets(w http.ResponseWriter, r *http.Request) {
	query := "SELECT id, user_id, content FROM tweets ORDER BY likes DESC LIMIT 10"

	rows, err := pg.DB.Query(query)
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

	rows, err := pg.DB.Query(query, searchQuery)
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

	subscribedUserIDs, err := services.GetSubscribedUserIDs(currentUserID)
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

	rows, err := pg.DB.Query(query)
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
