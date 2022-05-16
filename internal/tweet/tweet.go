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

// tweet is an internal package to provide the core functionality required by ajtweet.
package tweet

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// Tweet represents a single scheduled tweet to be sent to Twitter.
type Tweet struct {
	Id            uuid.UUID // The unique identifier for the tweet.
	Message       string    // The message to be posted to twitter.
	ScheduledTime time.Time // The preferred scheduled time at which the tweet needs to be sent.
}

// Create a new Tweet given the specified message and preferred scheduled time.
func New(message string, scheduledTime time.Time) Tweet {
	tweet := Tweet{
		Id:            uuid.New(),
		Message:       message,
		ScheduledTime: scheduledTime,
	}
	return tweet
}

// Return true if the tweet needs to be sent as soon as possible.
func (tweet Tweet) SendNow() bool {
	return tweet.ScheduledTime.Before(time.Now())
}

// Stringer implementation.
func (tweet Tweet) String() string {
	return fmt.Sprintf("id: %s, time: %s, tweet: %s", tweet.Id, tweet.ScheduledTime, tweet.Message)
}
