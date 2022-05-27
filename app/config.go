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

import "os"

// Configuration data used by the Application.
type Config struct {
	Datastore Datastore
	Send      Send

	Lockfile string // File path of where the lock file will be created.
}

// Datastore configures how the tweets are stored by the Application.
type Datastore struct {
	Filepath string // File path of where the tweets should be stored.
}

// Send parameters
type Send struct {
	Max   int // The maximum number of tweets to send in this call of the app.
	Delay int // The number of seconds to delay between each sending of a tweet.

	Authentication Authentication
}

// Authentication details for the Twitter API
type Authentication struct {
	APIKey    string `mapstructure:"api_key"`    // Consumer / API Key
	APISecret string `mapstructure:"api_secret"` // Consumer / API Secret

	OAuth1 OAuth1
}

// OAuth1.0a authentication for the Twitter API
type OAuth1 struct {
	Token  string // User token
	Secret string // User token secret
}

const (
	envAPIKey       = "AJTWEET_API_KEY"
	envAPISecret    = "AJTWEET_API_SECRET"
	envOAuth1Token  = "AJTWEET_ACCESS_TOKEN"
	envOAuth1Secret = "AJTWEET_ACCESS_SECRET"

	defaultSendMax      = 10
	defaultSendDelay    = 1
	defaultSendLockfile = "./ajtweet.lock"
)

// Create a new Config and set the default values required
func NewConfig() Config {
	var config Config
	config.Send.Max = defaultSendMax
	config.Send.Delay = defaultSendDelay
	config.Lockfile = defaultSendLockfile
	return config
}

// Configure values from matching environment variables
func (config *Config) PopulateFromEnv() {
	if value, present := os.LookupEnv(envAPIKey); present {
		config.Send.Authentication.APIKey = value
	}

	if value, present := os.LookupEnv(envAPISecret); present {
		config.Send.Authentication.APISecret = value
	}

	if value, present := os.LookupEnv(envOAuth1Token); present {
		config.Send.Authentication.OAuth1.Token = value
	}

	if value, present := os.LookupEnv(envOAuth1Secret); present {
		config.Send.Authentication.OAuth1.Secret = value
	}
}
