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
	"strings"
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
func EditTweet(w http.ResponseWriter, r *http.Request) {
	idTweet := mux.Vars(r)["id_tweet"]
	var (
		key    = []string{}
		values = []any{}
	)
	apikey := r.Header.Get("X-API-KEY")
	cookie, err := r.Cookie("session")
	var sessionID string
	if apikey != "" {
		sessionID = apikey
	} else {
		sessionID = cookie.Value
	}
	if apikey == "" && (err != nil || cookie == nil) {
		fmt.Println(err)
		services.ReturnErr(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	var updatedTweet Serviceuser.Tweet
	err = json.NewDecoder(r.Body).Decode(&updatedTweet)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	query := "SELECT user_id FROM user_session WHERE login_token = $1"
	var userID int
	err = pg.DB.QueryRow(query, sessionID).Scan(&userID)
	if err != nil {
		services.ReturnErr(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if updatedTweet.Text != "" {
		values = append(values, updatedTweet.Text)
		key = append(key, "text = $"+strconv.Itoa(len(key)+1))
	}
	if updatedTweet.Public == true {
		values = append(values, updatedTweet.Public, false, false, false)
		key = append(key, "public = $"+strconv.Itoa(len(key)+1), "only_followers = $"+strconv.Itoa(len(key)+2), "only_mutual_followers = $"+strconv.Itoa(len(key)+3), "only_me = $"+strconv.Itoa(len(key)+4))
	} else if updatedTweet.OnlyFollowers == true {
		values = append(values, updatedTweet.OnlyFollowers, false, false, false)
		key = append(key, "only_followers = $"+strconv.Itoa(len(key)+1), "public = $"+strconv.Itoa(len(key)+2), "only_mutual_followers = $"+strconv.Itoa(len(key)+3), "only_me = $"+strconv.Itoa(len(key)+4))
	} else if updatedTweet.OnlyMutualFollowers == true {
		values = append(values, updatedTweet.OnlyMutualFollowers, false, false, false)
		key = append(key, "only_mutual_followers = $"+strconv.Itoa(len(key)+1), "public = $"+strconv.Itoa(len(key)+2), "only_followers = $"+strconv.Itoa(len(key)+3), "only_me = $"+strconv.Itoa(len(key)+4))
	} else if updatedTweet.OnlyMe == true {
		values = append(values, updatedTweet.OnlyMe, false, false, false)
		key = append(key, "only_me = $"+strconv.Itoa(len(key)+1), "public = $"+strconv.Itoa(len(key)+2), "only_followers = $"+strconv.Itoa(len(key)+3), "only_mutual_followers = $"+strconv.Itoa(len(key)+4))
	}
	keystring := strings.Join(key, ", ")
	values = append(values, idTweet)
	fmt.Println(idTweet, "<<<<<<")
	if cookie != nil || apikey != "" {
		query := fmt.Sprintf("UPDATE tweets SET %s WHERE tweet_id = $%d", keystring, len(values))
		fmt.Println(query)
		fmt.Println(values)
		_, err = pg.DB.Exec(query, values...)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)

		}
	} else {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
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
