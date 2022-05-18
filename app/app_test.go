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
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"testing"
	"time"

	"github.com/andrejacobs/ajtweet-cli/internal/tweet"
	"github.com/google/uuid"
)

func TestAdd(t *testing.T) {
	app := Application{}

	if count := len(app.tweets.Tweets); count != 0 {
		t.Fatalf("Expected 0 tweets, Result: %d", count)
	}

	if err := app.Add("Tweet 1", "1942-04-24T10:42:42Z"); err != nil {
		t.Fatal(err)
	}

	if err := app.Add("Tweet 2", "2006-03-05T14:55:02Z"); err != nil {
		t.Fatal(err)
	}

	if count := len(app.tweets.Tweets); count != 2 {
		t.Fatalf("Expected 2 tweets, Result: %d", count)
	}

	expectedTime1, _ := parseTime("1942-04-24T10:42:42Z")
	if app.tweets.Tweets[0].Message != "Tweet 1" || app.tweets.Tweets[0].ScheduledTime != expectedTime1 {
		t.Fatal("Tweet 1 does not meet expectations")
	}

	expectedTime2, _ := parseTime("2006-03-05T14:55:02Z")
	if app.tweets.Tweets[1].Message != "Tweet 2" || app.tweets.Tweets[1].ScheduledTime != expectedTime2 {
		t.Fatal("Tweet 2 does not meet expectations")
	}

}

func TestAddInvalidTime(t *testing.T) {
	app := Application{}

	if err := app.Add("Tweet 1", "NOT VALID"); errors.Is(err, &time.ParseError{}) {
		t.Fatal(err)
	}
}

func TestConfigureAndSave(t *testing.T) {
	app := Application{}

	// Create a temporary file and delete it so that we can use the filename
	tempFile, err := os.CreateTemp("", "")
	if err != nil {
		t.Fatalf("Failed to create a temporary file. Error: %s", err)
	}
	os.Remove(tempFile.Name())
	defer os.Remove(tempFile.Name())

	// Configure the app
	config := Config{
		Datastore: Datastore{
			Filepath: tempFile.Name(),
		},
	}

	if err := app.Configure(config); err != nil {
		t.Fatalf("Failed to configure the app. Error: %s", err)
	}

	// Add a tweet
	if err := app.Add("Tweet 1", "1942-04-24T10:42:42Z"); err != nil {
		t.Fatal(err)
	}

	// Save
	if err := app.Save(); err != nil {
		t.Fatalf("Failed to save. Error: %s", err)
	}

	// Create a new app and check that configure loaded correctly
	app2 := Application{}
	if err := app2.Configure(config); err != nil {
		t.Fatalf("Failed to configure the app. Error: %s", err)
	}

	expectedTime1, _ := parseTime("1942-04-24T10:42:42Z")
	if app.tweets.Tweets[0].Message != "Tweet 1" || app.tweets.Tweets[0].ScheduledTime != expectedTime1 {
		t.Fatal("Failed to load the tweets as expected")
	}

}

func TestList(t *testing.T) {
	app := Application{}

	if err := app.Add("Tweet 1", time.Now().Format(time.RFC3339)); err != nil {
		t.Fatal(err)
	}

	if err := app.Add("Tweet 2", time.Now().Add(-5*time.Minute).Format(time.RFC3339)); err != nil {
		t.Fatal(err)
	}

	if err := app.Add("Tweet 3", time.Now().Add(5*time.Minute).Format(time.RFC3339)); err != nil {
		t.Fatal(err)
	}

	var buffer bytes.Buffer
	if err := app.List(&buffer); err != nil {
		t.Fatal(err)
	}

	expectedTweets := app.tweets.List()
	expected := fmt.Sprintf(
		`id: %s
time: %s [send now!]
tweet: %s

id: %s
time: %s [send now!]
tweet: %s

id: %s
time: %s
tweet: %s

`, expectedTweets[0].Id, expectedTweets[0].ScheduledTime.Format(time.RFC3339), expectedTweets[0].Message,
		expectedTweets[1].Id, expectedTweets[1].ScheduledTime.Format(time.RFC3339), expectedTweets[1].Message,
		expectedTweets[2].Id, expectedTweets[2].ScheduledTime.Format(time.RFC3339), expectedTweets[2].Message)

	result := buffer.String()
	if result != expected {
		t.Fatalf("Expected:\n%q\n\nResult:\n%q\n", expected, result)
	}
}

func TestListJSON(t *testing.T) {
	app := Application{}

	if err := app.Add("Tweet 1", time.Now().Format(time.RFC3339)); err != nil {
		t.Fatal(err)
	}

	if err := app.Add("Tweet 2", time.Now().Add(-5*time.Minute).Format(time.RFC3339)); err != nil {
		t.Fatal(err)
	}

	var buffer bytes.Buffer
	if err := app.ListJSON(&buffer); err != nil {
		t.Fatal(err)
	}

	expectedTweets := app.tweets.List()
	expected := fmt.Sprintf(`[{"id":"%s","message":"%s","scheduledTime":"%s"},{"id":"%s","message":"%s","scheduledTime":"%s"}]`,
		expectedTweets[0].Id, expectedTweets[0].Message, expectedTweets[0].ScheduledTime.Format(time.RFC3339),
		expectedTweets[1].Id, expectedTweets[1].Message, expectedTweets[1].ScheduledTime.Format(time.RFC3339))

	result := buffer.String()
	if result != expected {
		t.Fatalf("Expected:\n%q\n\nResult:\n%q\n", expected, result)
	}
}

func TestDelete(t *testing.T) {
	app := Application{}

	if err := app.Add("Tweet 1", time.Now().Format(time.RFC3339)); err != nil {
		t.Fatal(err)
	}

	if err := app.Add("Tweet 2", time.Now().Add(-5*time.Minute).Format(time.RFC3339)); err != nil {
		t.Fatal(err)
	}

	if err := app.Delete(app.tweets.Tweets[0].Id.String()); err != nil {
		t.Fatal(err)
	}

	if err := app.Delete(app.tweets.Tweets[0].Id.String()); err != nil {
		t.Fatal(err)
	}

	if len(app.tweets.Tweets) != 0 {
		t.Fatal("Expected all tweets to have been deleted")
	}

	if err := app.Delete("NOT VALID UUID"); err == nil {
		t.Fatal("Expected an UUID error to be raised")
	}

	if err := app.Delete(uuid.NewString()); err == nil {
		t.Fatal("Expected an error since the item does not exist in the list")
	}

}

func TestDeleteAll(t *testing.T) {
	app := Application{}

	app.Add("Tweet 1", time.Now().Format(time.RFC3339))
	app.Add("Tweet 2", time.Now().Format(time.RFC3339))

	if err := app.DeleteAll(); err != nil {
		t.Fatal(err)
	}

	if len(app.tweets.Tweets) != 0 {
		t.Fatal("Expected all tweets to have been deleted")
	}
}

func TestSend(t *testing.T) {
	app := Application{}

	tempFile, err := os.CreateTemp("", "")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tempFile.Name())

	app.config.Datastore.Filepath = tempFile.Name()
	app.config.Send.Max = 100

	app.Add("Tweet 1", time.Now().Format(time.RFC3339))
	app.Add("Tweet 2", time.Now().Format(time.RFC3339))
	app.Add("Tweet 3", time.Now().Add(5*time.Minute).Format(time.RFC3339))

	sendable := make([]tweet.Tweet, 2)
	copy(sendable, app.tweets.Tweets)

	var buffer bytes.Buffer
	if err := app.Send(&buffer); err != nil {
		t.Fatal(err)
	}

	expectedOut := fmt.Sprintf(`Sending 1 of 2
id: %s
tweet: %s

Sending 2 of 2
id: %s
tweet: %s

`, sendable[0].Id, sendable[0].Message,
		sendable[1].Id, sendable[1].Message)

	if buffer.String() != expectedOut {
		t.Fatalf("Expected output:\n%q\nResult:\n%q", expectedOut, buffer.String())
	}

	if count := len(app.tweets.Tweets); count != 1 {
		t.Fatalf("Expected 1 tweet to remain. Result: %d", count)
	}
}

func TestSendMaxAndEmpty(t *testing.T) {
	app := Application{}

	tempFile, err := os.CreateTemp("", "")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tempFile.Name())

	const total = 10
	app.config.Datastore.Filepath = tempFile.Name()
	app.config.Send.Max = total / 2

	for i := 0; i < total; i++ {
		app.Add("Tweet", time.Now().Format(time.RFC3339))
	}

	if len(app.tweets.Tweets) != total {
		t.Fatalf("Expected %d tweets ready to be sent", total)
	}

	// Send first batch
	if err := app.Send(io.Discard); err != nil {
		t.Fatal(err)
	}

	if count := len(app.tweets.Tweets); count != app.config.Send.Max {
		t.Fatalf("Expected %d tweets remaining to be sent. Result: %d", app.config.Send.Max, count)
	}

	// Send second batch
	if err := app.Send(io.Discard); err != nil {
		t.Fatal(err)
	}

	if count := len(app.tweets.Tweets); count != 0 {
		t.Fatalf("Expected 0 tweets remaining to be sent. Result: %d", count)
	}

	// Send when there is nothing to send
	var buffer bytes.Buffer
	if err := app.Send(&buffer); err != nil {
		t.Fatal(err)
	}

	if buffer.String() != "" {
		t.Fatalf(`Expected "". Result: %q`, buffer.String())
	}
}
