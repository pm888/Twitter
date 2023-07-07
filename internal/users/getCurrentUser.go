package users

import (
	"Twitter_like_application/internal/database/pg"
	"Twitter_like_application/internal/services"
	"encoding/json"
	"net/http"
)

func GetCurrentProfile(w http.ResponseWriter, r *http.Request) {
	var (
		followerCount     int
		subscriptionCount int
	)
	userID := r.Context().Value("userID").(int)

	query := `
	SELECT
		u.name, u.birthdate, u.bio, u.location, u.nickname,
		COUNT(DISTINCT fs1.subscription_id) AS subscription_count,
		COUNT(DISTINCT fs2.follower_id) AS follower_count
	FROM
		users_tweeter u
	LEFT JOIN
		followers_subscriptions fs1 ON u.id = fs1.follower_id
	LEFT JOIN
		followers_subscriptions fs2 ON u.id = fs2.subscription_id
	WHERE
		u.id = $1
	GROUP BY
		u.id, u.name, u.birthdate, u.bio, u.location, u.nickname
`
	var user GetCurrentUser
	err := pg.DB.QueryRow(query, userID).Scan(
		&user.Name, &user.BirthDate, &user.Bio, &user.Location, &user.Nickname,
		&subscriptionCount, &followerCount,
	)
	if err != nil {
		services.ReturnErr(w, err.Error(), http.StatusInternalServerError)
		return
	}
	user.Followers = subscriptionCount
	user.Following = followerCount

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
