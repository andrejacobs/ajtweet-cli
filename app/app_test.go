package app

import (
	"errors"
	"os"
	"testing"
	"time"
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
