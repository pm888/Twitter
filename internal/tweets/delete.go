package tweets

import (
	"Twitter_like_application/internal/database/pg"
	"Twitter_like_application/internal/services"
	"database/sql"
	"github.com/gorilla/mux"
	"net/http"
)

func DeleteTweet(w http.ResponseWriter, r *http.Request) {
	tweetID := mux.Vars(r)["tweet_id"]
	userID := r.Context().Value("userID").(int)

	deleteQuery := "DELETE FROM tweets WHERE tweet_id = $1 AND user_id = $2 RETURNING true"
	var exists bool
	err := pg.DB.QueryRow(deleteQuery, tweetID, userID).Scan(&exists)
	if err != nil {
		if err == sql.ErrNoRows {
			services.ReturnErr(w, "Tweet not found", http.StatusNotFound)
		} else {
			services.ReturnErr(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	if !exists {
		services.ReturnErr(w, "You are not authorized to delete this tweet", http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
}
