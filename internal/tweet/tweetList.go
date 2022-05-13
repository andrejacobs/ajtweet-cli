package tweet

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
)

var (
	ErrExists = errors.New("Tweet already exists in the list")
	// ErrNotExists = errors.New("Tweet does not exist in the list")
)

type TweetList struct {
	tweets []Tweet
}

func (list *TweetList) Find(id uuid.UUID) (bool, int) {
	for i, tweet := range list.tweets {
		if tweet.Id == id {
			return true, i
		}
	}

	return false, -1
}

func (list *TweetList) Add(tweet Tweet) error {
	if found, _ := list.Find(tweet.Id); found {
		return fmt.Errorf("%w: %q", ErrExists, tweet.Id)
	}

	list.tweets = append(list.tweets, tweet)
	return nil
}
