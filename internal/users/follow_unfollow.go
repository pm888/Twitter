package users

import (
	"Twitter_like_application/internal/database/pg"
	"Twitter_like_application/internal/services"
	"encoding/json"
	"fmt"
	"net/http"
)

func FollowUser(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int)

	var targetUserID TargetUser
	err := json.NewDecoder(r.Body).Decode(&targetUserID)
	if err != nil {
		services.ReturnErr(w, err.Error(), http.StatusBadRequest)
		return
	}

	var count int
	err = pg.DB.QueryRow("SELECT COUNT(*) FROM followers_subscriptions WHERE follower_id = $1 AND subscription_id = $2", userID, targetUserID.ID).Scan(&count)
	if err != nil {
		services.ReturnErr(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if count > 0 {
		services.ReturnErr(w, "User is already subscribed to the target user", http.StatusBadRequest)
		return
	}

	_, err = pg.DB.Exec("INSERT INTO followers_subscriptions (follower_id, subscription_id) VALUES ($1, $2)", userID, targetUserID.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"status":  "success",
		"message": fmt.Sprintf("id %d folloer to id %d", userID, targetUserID.ID),
	}
	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(response); err != nil {
		services.ReturnErr(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func UnfollowUser(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int)

	var targetUserID TargetUser
	err := json.NewDecoder(r.Body).Decode(&targetUserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = pg.DB.Exec("DELETE FROM followers_subscriptions WHERE follower_id = $1 AND subscription_id = $2", userID, targetUserID.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"status":  "success",
		"message": fmt.Sprintf("id %d unfolloer from id %d", userID, targetUserID.ID),
	}
	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(response); err != nil {
		services.ReturnErr(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
