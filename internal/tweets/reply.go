package tweets

import (
	"Twitter_like_application/internal/database/pg"
	"Twitter_like_application/internal/services"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

func Reply(w http.ResponseWriter, r *http.Request) {
	idTweet := mux.Vars(r)["id_tweet"]
	userID := r.Context().Value("userID").(int)
	var newReply TweetReply
	err := json.NewDecoder(r.Body).Decode(&newReply)
	if err != nil {
		services.ReturnErr(w, err.Error(), http.StatusBadRequest)
		return
	}
	query := `INSERT INTO reply (text_reply,parent_tweet_id,user_id,timestamp)
		VALUES ($1, $2, $3, $4) RETURNING reply_id`
	err = pg.DB.QueryRowContext(r.Context(), query, newReply.Text, idTweet, userID, time.Now()).Scan(&newReply.ReplyId)
	if err != nil {
		services.ReturnErr(w, err.Error(), http.StatusInternalServerError)
		return

	}
	response := map[string]interface{}{
		"status":  "success",
		"message": "Reply added",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

}
