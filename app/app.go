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
package app

import (
	"time"

	"github.com/andrejacobs/ajtweet-cli/internal/tweet"
)

type Application struct {
	config Config
	tweets tweet.TweetList
}

func (app *Application) Configure(config Config) error {
	app.config = config

	if err := app.tweets.Load(app.config.Datastore.Filepath); err != nil {
		return err
	}
	return nil
}

func (app *Application) Save() error {
	//AJ### TODO: Need to make the save atomic so that we corrupt good data
	if err := app.tweets.Save(app.config.Datastore.Filepath); err != nil {
		return err
	}
	return nil
}

func (app *Application) Add(message string, scheduledTimeString string) error {
	scheduledTime, err := parseTime(scheduledTimeString)
	if err != nil {
		return err
	}

	tweet := tweet.New(message, scheduledTime)
	if err := app.tweets.Add(tweet); err != nil {
		return err
	}
	return nil
}

func parseTime(timeString string) (time.Time, error) {
	return time.Parse(time.RFC3339, timeString)
}
