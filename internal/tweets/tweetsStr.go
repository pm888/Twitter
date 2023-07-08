package tweets

import (
	"github.com/go-playground/validator/v10"
	"time"
)

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
type CreatNewTweet struct {
	TweetID             int
	Text                string `json:"text" validate:"required,checkTweetText"`
	CreatedAt           time.Time
	Public              bool `json:"public"`
	OnlyFollowers       bool `json:"only_followers"`
	OnlyMutualFollowers bool `json:"only_mutual_followers"`
	OnlyMe              bool `json:"only_me"`
}

type TweetValid struct {
	Validate *validator.Validate
	ValidErr map[string]string
}

type EditTweetRequest struct {
	Text                string `json:"text" validate:"required,checkTweetText"`
	Public              bool   `json:"public"`
	OnlyFollowers       bool   `json:"only_followers"`
	OnlyMutualFollowers bool   `json:"only_mutual_followers"`
	OnlyMe              bool   `json:"only_me"`
}
