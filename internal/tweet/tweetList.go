package tweet

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/google/uuid"
)

var (
	ErrExists = errors.New("Tweet already exists in the list")
	// ErrNotExists = errors.New("Tweet does not exist in the list")
)

type TweetList struct {
	Tweets []Tweet
}

func (list *TweetList) Find(id uuid.UUID) (bool, int) {
	for i, tweet := range list.Tweets {
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

	list.Tweets = append(list.Tweets, tweet)
	return nil
}

func (list *TweetList) Load(filePath string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}

	if len(data) == 0 {
		return nil
	}

	return json.Unmarshal(data, list)
}

func (list *TweetList) Save(filePath string) error {
	jsonData, err := json.Marshal(list)
	if err != nil {
		return err
	}

	return os.WriteFile(filePath, jsonData, 0644)
}
