package tweets

import (
	"Twitter_like_application/internal/database/pg"
	_ "Twitter_like_application/internal/database/pg"
	Serviceuser "Twitter_like_application/internal/users"
	"database/sql"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
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

	query = `INSERT INTO tweets (tweet_id, user_id, text, created_at, parent_tweet_id, public, only_followers, only_mutual_followers, only_me,retweet)
	 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING tweet_id`
	err = pg.DB.QueryRowContext(r.Context(), query, newTweet.TweetID, newTweet.UserID, newTweet.Text, time.Now(), newTweet.ParentTweetId, newTweet.Public, newTweet.OnlyFollowers, newTweet.OnlyMutualFollowers, newTweet.OnlyMe, newTweet.Retweet).Scan(&newTweet.TweetID)
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

func LikeTweet(w http.ResponseWriter, r *http.Request) {
	idTweet := mux.Vars(r)["id_tweet"]
	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var exists bool

	query := "SELECT EXISTS (SELECT 1 FROM likes WHERE user_id = $1 AND tweet_id = $2)"
	err := pg.DB.QueryRow(query, userID, idTweet).Scan(&exists)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if exists {
		http.Error(w, "Tweet already liked", http.StatusBadRequest)
		return
	}

	query = "INSERT INTO likes (tweet_id, user_id, timestamp) VALUES ($1, $2, $3)"
	_, err = pg.DB.Exec(query, idTweet, userID, time.Now())
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
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
