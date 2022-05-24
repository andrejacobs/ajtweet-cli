package app

import (
	"context"
	"os"

	"github.com/michimani/gotwi"
	"github.com/michimani/gotwi/tweet/managetweet"
	"github.com/michimani/gotwi/tweet/managetweet/types"
)

// Source code based on the example:
// https://github.com/michimani/gotwi/blob/main/_examples/2_post_delete_tweet/main.go
// https://github.com/andrejacobs/example_twitter_golang

type twitterClient struct {
	gotwiClient *gotwi.Client
}

func newOAuth1Client(auth Authentication) (*twitterClient, error) {

	os.Setenv(gotwi.APIKeyEnvName, auth.APIKey)
	os.Setenv(gotwi.APIKeySecretEnvName, auth.APISecret)

	in := &gotwi.NewClientInput{
		AuthenticationMethod: gotwi.AuthenMethodOAuth1UserContext,
		OAuthToken:           auth.OAuth1.Token,
		OAuthTokenSecret:     auth.OAuth1.Secret,
	}

	client, err := gotwi.NewClient(in)
	if err != nil {
		return nil, err
	}

	return &twitterClient{gotwiClient: client}, nil
}

func sendTweet(client *twitterClient, text string) (string, error) {
	p := &types.CreateInput{
		Text: gotwi.String(text),
	}

	res, err := managetweet.Create(context.Background(), client.gotwiClient, p)
	if err != nil {
		return "", err
	}

	return gotwi.StringValue(res.Data.ID), nil
}
