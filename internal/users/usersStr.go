package users

import (
	"gopkg.in/go-playground/validator.v9"
	"time"
)

type Users struct {
	ID                 int
	Name               string `json:"name" validate:"omitempty"`
	Password           string `json:"password" validate:"omitempty"`
	Email              string `json:"email" validate:"omitempty,email"`
	EmailToken         string
	ConfirmEmailToken  bool
	ResetPasswordToken string
	BirthDate          string `json:"birthdate" validate:"omitempty"`
	Nickname           string `json:"nickname" validate:"omitempty"`
	Bio                string `json:"bio" validate:"omitempty"`
	Location           string `json:"location" validate:"omitempty"`
	Tweet
}

type Tweet struct {
	TweetID             int       `json:"tweet_id"`
	UserID              int       `json:"user_id"`
	Author              string    `json:"author"`
	Text                string    `json:"text"`
	CreatedAt           time.Time `json:"created_at"`
	LikeCount           int       `json:"like_count"`
	Retweet             int       `json:"repost"`
	Public              bool      `json:"public"`
	OnlyFollowers       bool      `json:"only_followers"`
	OnlyMutualFollowers bool      `json:"only_mutual_followers"`
	OnlyMe              bool      `json:"only_me"`
	LoginToken          string
	ParentTweetId       int `json:"parent_tweet_id"`
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

type UsersLogin struct {
	Usermail string `json:"email_logIN"`
	Password string `json:"password_logIN"`
}

type Tweeter_like struct {
	Autor      int `json:"autor"`
	Id_post    int `json:"id_post"`
	Whose_like int `json:"whose_like"`
}
type UserVal struct {
	validate *validator.Validate
}
type EditUser struct {
	ID        int
	Name      string `json:"name" validate:"omitempty,checkName"`
	Password  string `json:"password" validate:"omitempty,checkPassword"`
	Email     string `json:"email" validate:"omitempty,email"`
	BirthDate string `json:"birthdate" validate:"omitempty,checkDataTime"`
	Nickname  string `json:"nickname" validate:"omitempty,checkNickname"`
	Bio       string `json:"bio" validate:"omitempty,checkBio"`
	Location  string `json:"location" validate:"omitempty,checkLocation"`
}
