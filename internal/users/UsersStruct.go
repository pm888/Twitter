package Serviceuser

import (
	"time"
)

type Users struct {
	ID                 int    `json:"id"`
	Name               string `json:"name"`
	Password           string `json:"password"`
	Email              string `json:"email"`
	EmailToken         string
	ConfirmEmailToken  bool
	ResetPasswordToken string
	BirthDate          string `json:"birth_date"`
	Nickname           string `json:"nickname"`
	Bio                string `json:"bio"`
	Location           string `json:"location"`
	Tweet
	Following []int
	Followers []int
}

type ReplaceMyData struct {
	NewName      string `json:"new_name"`
	NewPassword  string `json:"new_password"`
	NewEmail     string `json:"new_email"`
	NewBirthDate string `json:"new_birth_date"`
	NewNickname  string `json:"new_nickname"`
	NewBio       string `json:"new_bio"`
	NewLocation  string `json:"new_location"`
}

type Tweet struct {
	TweetID             int       `json:"tweet_id"`
	UserID              int       `json:"user_id"`
	Author              string    `json:"author"`
	Text                string    `json:"text"`
	CreatedAt           time.Time `json:"created_at"`
	Like                int
	Repost              int
	Public              bool `json:"public"`
	OnlyFollowers       bool `json:"only_followers"`
	OnlyMutualFollowers bool `json:"only_mutual_followers"`
	OnlyMe              bool `json:"only_me"`
}
type ReplayTweet struct {
	Tweet
}

type DeleteUserST struct {
	UserIdDeleting int `json:"delete_id"`
}

type ResetPasswordUser struct {
	UserResetPassword int `json:"user_reset_password"`
}

type FollowingForUser struct {
	Writer     int `json:"writer"`
	Subscriber int `json:"subscriber"`
	Users
}
