/*
Copyright © 2022 André Jacobs

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
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
