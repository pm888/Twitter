package tweets

import (
	"Twitter_like_application/internal/database/pg"
	"Twitter_like_application/internal/services"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"net/http"
)

func EditTweet(w http.ResponseWriter, r *http.Request) {
	key := []string{}
	tweetID := mux.Vars(r)["id_tweet"]
	userID := r.Context().Value("userID").(int)
	tweetValid := &TweetValid{
		Validate: validator.New(),
		ValidErr: make(map[string]string),
	}
	if err := RegisterTweetValidations(tweetValid); err != nil {
		fmt.Println(err)
	}

	var updatedTweet EditTweetRequest
	var tweet Tweet
	err := json.NewDecoder(r.Body).Decode(&updatedTweet)
	if err != nil {
		services.ReturnErr(w, err.Error(), http.StatusBadRequest)
		return
	}
	if !CheckVisibility(&updatedTweet, tweetValid) {
		services.ReturnErr(w, tweetValid.Error(), http.StatusInternalServerError)
		return
	}
	query := "SELECT user_id, public, only_followers, only_mutual_followers, only_me FROM tweets WHERE tweet_id = $1"
	err = pg.DB.QueryRow(query, tweetID).Scan(&tweet.UserID, &tweet.Public, &tweet.OnlyFollowers, &tweet.OnlyMutualFollowers, &tweet.OnlyMe)
	if err != nil {
		if err == sql.ErrNoRows {
			services.ReturnErr(w, "Tweet not found", http.StatusNotFound)
		} else {
			services.ReturnErr(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	if tweet.UserID != userID {
		http.Error(w, "it isn't your tweet", http.StatusUnauthorized)
		return
	}
	if tweet.Public != updatedTweet.Public {

	}
	updateValues := make([]string, len(key))
	for i := range key {
		updateValues[i] = fmt.Sprintf("%s = $%d", key[i], i+1)
	}

	valuesVision := make([]any, 4)
	valuesVision[0] = func() any {
		if tweet.Public != updatedTweet.Public {
			return updatedTweet.Public
		}
		return tweet.Public
	}()
	valuesVision[1] = func() any {
		if tweet.OnlyFollowers != updatedTweet.OnlyFollowers {
			return updatedTweet.OnlyFollowers
		}
		return tweet.OnlyFollowers
	}()
	valuesVision[2] = func() any {
		if tweet.OnlyMutualFollowers != updatedTweet.OnlyMutualFollowers {
			return updatedTweet.OnlyMutualFollowers
		}
		return tweet.OnlyMutualFollowers
	}()
	valuesVision[3] = func() any {
		if tweet.OnlyMe != updatedTweet.OnlyMe {
			return updatedTweet.OnlyMe
		}
		return tweet.OnlyMe
	}()

	query = "UPDATE tweets SET text = $1, public = $2, only_followers = $3, only_mutual_followers = $4, only_me = $5 WHERE tweet_id = $6"
	_, err = pg.DB.ExecContext(r.Context(), query, updatedTweet.Text, valuesVision[0], valuesVision[1], valuesVision[2], valuesVision[3], tweetID)
	if err != nil {
		services.ReturnErr(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response := map[string]interface{}{
		"status":  "success",
		"message": "Tweet updated successfully",
	}
	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(response); err != nil {
		services.ReturnErr(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
