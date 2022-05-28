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
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/andrejacobs/ajtweet-cli/internal/tweet"
	"github.com/fatih/color"
	"github.com/google/uuid"
)

var (
	ErrLockfileExists = errors.New("another instance is running and have acquired the lock file")
)

// The main "context" used in the application.
type Application struct {
	config Config
	tweets tweet.TweetList
}

// Configure and load any existing tweets to be used by the Application.
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

	whiteBold := color.New(color.FgWhite, color.Bold).SprintFunc()
	greenBold := color.New(color.FgGreen, color.Bold).SprintFunc()
	cyan := color.New(color.FgCyan).SprintFunc()

	for _, tw := range app.tweets.List() {
		if _, err := fmt.Fprintf(out, "id: %s\n", cyan(tw.Id)); err != nil {
			return err
		}

		if _, err := fmt.Fprintf(out, "time: %s", tw.ScheduledTime.Format(time.RFC3339)); err != nil {
			return err
		}

		if tw.SendNow() {
			if _, err := fmt.Fprintf(out, " %s", greenBold("[send now!]")); err != nil {
				return err
			}
		}

		if _, err := fmt.Fprintf(out, "\ntweet: %s\n\n", whiteBold(tw.Message)); err != nil {
			return err
		}
	}

	return nil
}

// Write the list of scheduled tweets that still need to be sent in a JSON encoding to the specified io.Writer.
func (app *Application) ListJSON(out io.Writer) error {
	tweets := app.tweets.List()
	jsonData, err := json.Marshal(tweets)
	if err != nil {
		return err
	}

	out.Write(jsonData)
	return nil
}

// Delete the tweet matching the specified identifier.
func (app *Application) Delete(idString string) error {
	id, err := uuid.Parse(idString)
	if err != nil {
		return err
	}

	return app.tweets.Delete(id)
}

// Delete all the tweets.
func (app *Application) DeleteAll() error {
	return app.tweets.DeleteAll()
}

var (
	// The tweet has already been added to the list.
	ErrMissingAuth = errors.New("authentication parameters are missing")
)

// Send any scheduled tweets.
func (app *Application) Send(out io.Writer, dryRun bool) error {

	var client *twitterClient

	configure := func(out io.Writer, dryRun bool) error {
		if app.config.Send.Authentication.APIKey == "" {
			return fmt.Errorf("%w: API Key", ErrMissingAuth)
		}

		if app.config.Send.Authentication.APISecret == "" {
			return fmt.Errorf("%w: API Secret", ErrMissingAuth)
		}

		if app.config.Send.Authentication.OAuth1.Token == "" {
			return fmt.Errorf("%w: OAuth 1 User token", ErrMissingAuth)
		}

		if app.config.Send.Authentication.OAuth1.Secret == "" {
			return fmt.Errorf("%w: OAuth 1 User secret", ErrMissingAuth)
		}

		var err error
		client, err = newOAuth1Client(app.config.Send.Authentication)
		if err != nil {
			return err
		}

		return nil
	}

	greenBold := color.New(color.FgHiGreen, color.Bold).SprintFunc()

	actual := func(out io.Writer, dryRun bool, tweet tweet.Tweet) error {
		if dryRun {
			return nil
		}

		id, err := sendTweet(client, tweet.Message)
		if err != nil {
			return err
		}

		fmt.Fprintf(out, "Twitter identifier: %s\n", greenBold(id))

		return nil
	}

	return app.send(out, dryRun, configure, actual)
}

type sendConfigure func(out io.Writer, dryRun bool) error
type sendActual func(out io.Writer, dryRun bool, tweet tweet.Tweet) error

func (app *Application) send(out io.Writer, dryRun bool,
	configure sendConfigure, actual sendActual) error {

	if err := configure(out, dryRun); err != nil {
		return err
	}

	sendable := app.tweets.ToSend(app.config.Send.Max, time.Now())
	sendCount := len(sendable)

	whiteBold := color.New(color.FgWhite, color.Bold).SprintFunc()
	cyan := color.New(color.FgCyan).SprintFunc()

	for i, tweet := range sendable {
		fmt.Fprintf(out, "Sending %d of %d\n", i+1, sendCount)

		if _, err := fmt.Fprintf(out, "id: %s\n", cyan(tweet.Id)); err != nil {
			return err
		}

		if _, err := fmt.Fprintf(out, "tweet: %s\n\n", whiteBold(tweet.Message)); err != nil {
			return err
		}

		if err := actual(out, dryRun, tweet); err != nil {
			return err
		}

		if err := app.tweets.Delete(tweet.Id); err != nil {
			return err
		}

		if !dryRun {
			if err := app.Save(); err != nil {
				return err
			}
		}

		if (app.config.Send.Delay > 0) && (i < sendCount-1) {
			fmt.Fprintf(out, "Delaying for %d seconds ...\n", app.config.Send.Delay)
			time.Sleep(time.Duration(app.config.Send.Delay) * time.Second)
		}
	}

	return nil
}

func parseTime(timeString string) (time.Time, error) {
	return time.Parse(time.RFC3339, timeString)
}

func (app *Application) isLocked() bool {
	if _, err := os.Stat(app.config.Lockfile); err == nil {
		return true
	}
	return false
}

func (app *Application) AcquireLock() error {
	if app.isLocked() {
		return fmt.Errorf("%w: %s", ErrLockfileExists, app.config.Lockfile)
	}

	pid := fmt.Sprintf("pid: %d\n", os.Getpid())
	if err := os.WriteFile(app.config.Lockfile, []byte(pid), 0644); err != nil {
		return err
	}

	return nil
}

func (app *Application) ReleaseLock() error {
	return os.Remove(app.config.Lockfile)
}
