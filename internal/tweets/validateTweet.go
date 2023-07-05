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
		v.ValidErr["name"] += "long name,"
	}
	return true
}
func RegisterUsersValidations(tweetValid *TweetValid) error {
	err := tweetValid.Validate.RegisterValidation("checkTweetText", func(fl validator.FieldLevel) bool {
		return CheckTweetText(fl, tweetValid)
	})
	if err != nil {
		return err
	}
	return nil
}
func CheckVisibility(newTweet *CreatNewTweet, v *TweetValid) bool {
	count := 0
	if newTweet.Public {
		v.ValidErr["visibility"] += "public true,"
		count++
	}
	if newTweet.OnlyMe {
		v.ValidErr["visibility"] += "onlyme true,"
		count++
	}
	if newTweet.OnlyFollowers {
		v.ValidErr["visibility"] += "onlyfollowers true,"
		count++
	}
	if newTweet.OnlyMutualFollowers {
		v.ValidErr["visibility"] += "onlymutualFollowers true,"
		count++
	}
	if count != 1 {
		v.ValidErr["visibility"] += "You need to choose only one variant "
		return false
	}
	return count == 1
}
