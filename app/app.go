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
