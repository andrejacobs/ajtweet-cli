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
	"sort"
	"time"

	"github.com/google/uuid"
)

var (
	// The tweet has already been added to the list.
	ErrExists    = errors.New("Tweet already exists in the list")
	ErrNotExists = errors.New("Tweet does not exist in the list")
)

// TweetList manages a collection of Tweets.
type TweetList struct {
	Tweets []Tweet
}

// Check if a tweet with the specified identifier can be found in the list.
// If the tweet can be found then true and the index of the tweet will be returned.
// If the tweet can not be found then false and -1 will be returned.
func (list *TweetList) Find(id uuid.UUID) (bool, int) {
	for i, tweet := range list.Tweets {
		if tweet.Id == id {
			return true, i
		}
	}

	return false, -1
}

// Add a new tweet to the list.
func (list *TweetList) Add(tweet Tweet) error {
	if found, _ := list.Find(tweet.Id); found {
		return fmt.Errorf("%w: %q", ErrExists, tweet.Id)
	}

	list.Tweets = append(list.Tweets, tweet)
	return nil
}

// Delete the tweet matching the specified identifier.
// If the tweet could not be found then an error will be returned.
func (list *TweetList) Delete(id uuid.UUID) error {
	found, index := list.Find(id)
	if !found {
		return fmt.Errorf("%w: %q", ErrNotExists, id)
	}

	list.Tweets = append(list.Tweets[:index], list.Tweets[index+1:]...)

	return nil
}

// Delete all the tweets from the list
func (list *TweetList) DeleteAll() error {
	list.Tweets = nil
	return nil
}

// Return a slice of tweets ordered by which tweets need to be sent first.
func (list *TweetList) List() []Tweet {
	result := make([]Tweet, len(list.Tweets))
	copy(result, list.Tweets)

	sort.SliceStable(result, func(i, j int) bool {
		return result[i].ScheduledTime.Before(result[j].ScheduledTime)
	})

	return result
}

// Return a slice of tweets that need to be send according to the specific time.
// max Is the maximum number of tweets to return.
// now Is the time to compare the tweet's scheduledAt time against.
func (list *TweetList) ToSend(max int, now time.Time) []Tweet {
	if max <= 0 {
		return nil
	}

	sendable := list.filter(func(tweet Tweet) bool {
		return tweet.SendWhen(now)
	})

	sort.SliceStable(sendable, func(i, j int) bool {
		return sendable[i].ScheduledTime.Before(sendable[j].ScheduledTime)
	})

	possibleMax := min(len(sendable), max)
	result := make([]Tweet, possibleMax)

	copy(result, sendable)
	return result
}

// Load the list of tweets from a JSON encoded file at the specified filePath.
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

// Save the list of tweets to a JSON encoded file at the specified filePath.
func (list *TweetList) Save(filePath string) error {
	jsonData, err := json.Marshal(list)
	if err != nil {
		return err
	}

	return os.WriteFile(filePath, jsonData, 0644)
}

// This is just bonkers that you have to implement this yourself!
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Determine if a tweet needs to be kept (true) or removed (false)
type filterFunc func(tweet Tweet) bool

// Filter the list of tweets given the filterFunc as criteria
func (list *TweetList) filter(keep filterFunc) []Tweet {
	result := make([]Tweet, 0, len(list.Tweets))

	for _, tweet := range list.Tweets {
		if keep(tweet) {
			result = append(result, tweet)
		}
	}

	return result
}
