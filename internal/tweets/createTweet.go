package tweets

import (
	"Twitter_like_application/internal/database/pg"
	"Twitter_like_application/internal/services"
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"net/http"
	"time"
)

func CreateNewTweet(w http.ResponseWriter, r *http.Request) {
	tweetValid := &TweetValid{
		Validate: validator.New(),
		ValidErr: make(map[string]string),
	}
	if err := RegisterTweetValidations(tweetValid); err != nil {
		fmt.Println(err)
	}
	userID := r.Context().Value("userID").(int)
	var newTweet CreatNewTweet
	err := json.NewDecoder(r.Body).Decode(&newTweet)
	if err != nil {
		services.ReturnErr(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = tweetValid.Validate.Struct(newTweet)
	if !CheckVisibility(&newTweet, tweetValid) {
		services.ReturnErr(w, tweetValid.Error(), http.StatusInternalServerError)
		return
	}

	query := `INSERT INTO tweets (user_id, text, created_at, public, only_followers, only_mutual_followers, only_me)
		VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING tweet_id`
	err = pg.DB.QueryRowContext(r.Context(), query, userID, newTweet.Text, time.Now(), newTweet.Public, newTweet.OnlyFollowers, newTweet.OnlyMutualFollowers, newTweet.OnlyMe).Scan(&newTweet.TweetID)
	if err != nil {
		services.ReturnErr(w, err.Error(), http.StatusInternalServerError)
		return

	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newTweet)

	return
}
