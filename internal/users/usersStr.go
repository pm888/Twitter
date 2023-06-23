package users

import (
	"time"
)

type Users struct {
	ID                 int
	Name               string `json:"name"`
	Password           string `json:"password"`
	Email              string `json:"email"`
	EmailToken         string
	ConfirmEmailToken  bool
	ResetPasswordToken string
	BirthDate          string `json:"birthdate"`
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
	LikeCount           int       `json:"like_count"`
	Repost              int       `json:"repost"`
	Public              bool      `json:"public"`
	OnlyFollowers       bool      `json:"only_followers"`
	OnlyMutualFollowers bool      `json:"only_mutual_followers"`
	OnlyMe              bool      `json:"only_me"`
	LoginToken          string
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
