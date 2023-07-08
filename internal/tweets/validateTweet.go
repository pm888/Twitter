package tweets

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"strings"
)

var (
	maxLenghtTweet = 400
)

func (v *TweetValid) Error() string {
	var pairs []string
	for k, v := range v.ValidErr {
		pairs = append(pairs, fmt.Sprintf("%s: %s", k, v))
	}

	result := strings.Join(pairs, "; ")
	return result
}

func CheckTweetText(fl validator.FieldLevel, v *TweetValid) bool {
	text := fl.Field().String()
	if len(text) > maxLenghtTweet {
		v.ValidErr["name"] += "long text,"
	}
	return true
}
func RegisterTweetValidations(tweetValid *TweetValid) error {
	err := tweetValid.Validate.RegisterValidation("checkTweetText", func(fl validator.FieldLevel) bool {
		return CheckTweetText(fl, tweetValid)
	})
	if err != nil {
		return err
	}
	return nil
}
func CheckVisibility(tweet interface{}, v *TweetValid) bool {
	count := 0
	var public, onlyMe, onlyFollowers, onlyMutualFollowers bool

	switch t := tweet.(type) {
	case *CreatNewTweet:
		public = t.Public
		onlyMe = t.OnlyMe
		onlyFollowers = t.OnlyFollowers
		onlyMutualFollowers = t.OnlyMutualFollowers
	case *EditTweetRequest:
		public = t.Public
		onlyMe = t.OnlyMe
		onlyFollowers = t.OnlyFollowers
		onlyMutualFollowers = t.OnlyMutualFollowers
	default:
		return false
	}

	if public {
		v.ValidErr["visibility"] += "public true,"
		count++
	}
	if onlyMe {
		v.ValidErr["visibility"] += "onlyme true,"
		count++
	}
	if onlyFollowers {
		v.ValidErr["visibility"] += "onlyfollowers true,"
		count++
	}
	if onlyMutualFollowers {
		v.ValidErr["visibility"] += "onlymutualFollowers true,"
		count++
	}

	if count != 1 {
		v.ValidErr["visibility"] += "You need to choose only one variant "
		return false
	}

	return count == 1
}
