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

// app package provides the core interface for ajtweet.
// It is used in decoupling from the CLI and thus makes it easier to unit-test.
package app

import (
	"fmt"
	"io"
	"time"

	"github.com/andrejacobs/ajtweet-cli/internal/tweet"
)

// The main "context" used in the application.
type Application struct {
	config Config
	tweets tweet.TweetList
}

// Configure and load any exsiting tweets to be used by the Application.
func (app *Application) Configure(config Config) error {
	app.config = config

	if err := app.tweets.Load(app.config.Datastore.Filepath); err != nil {
		return err
	}
	return nil
}

// Save any changes made by the Application.
func (app *Application) Save() error {
	//AJ### TODO: Need to make the save atomic so that we corrupt good data
	if err := app.tweets.Save(app.config.Datastore.Filepath); err != nil {
		return err
	}
	return nil
}

// Add a new scheduled tweet to the Application.
// The scheduledTimeString must be in the RFC 3339 standard, e.g. 2006-03-05T10:42:01Z
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

// Write the list of scheduled tweets that still need to be sent to the specified io.Writer.
func (app *Application) List(out io.Writer) error {

	for _, tw := range app.tweets.List() {
		if _, err := fmt.Fprintf(out, "id: %s\n", tw.Id); err != nil {
			return err
		}

		if _, err := fmt.Fprintf(out, "time: %s", tw.ScheduledTime.Format(time.RFC3339)); err != nil {
			return err
		}

		if tw.SendNow() {
			if _, err := fmt.Fprint(out, " [send now!]"); err != nil {
				return err
			}
		}
		if _, err := fmt.Fprintf(out, "\ntweet: %s\n\n", tw.Message); err != nil {
			return err
		}
	}

	return nil
}

func parseTime(timeString string) (time.Time, error) {
	return time.Parse(time.RFC3339, timeString)
}
