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
	"os"
	"testing"
)

func TestPopulateFromEnv(t *testing.T) {
	auth := Authentication{
		APIKey:    "a",
		APISecret: "b",
		OAuth1: OAuth1{
			Token:  "c",
			Secret: "d",
		},
	}
	var config Config
	config.Send.Authentication = auth

	if config.Send.Authentication.APIKey != "a" ||
		config.Send.Authentication.APISecret != "b" ||
		config.Send.Authentication.OAuth1.Token != "c" ||
		config.Send.Authentication.OAuth1.Secret != "d" {
		t.Fatalf("Configuration not initialzed as expected")
	}

	os.Setenv(envAPIKey, "apiKey")
	os.Setenv(envAPISecret, "apiSecret")
	os.Setenv(envOAuth1Token, "userToken")
	os.Setenv(envOAuth1Secret, "userSecret")

	config.PopulateFromEnv()

	if config.Send.Authentication.APIKey != "apiKey" {
		t.Fatalf("Expected %q = %q. Result %q", envAPIKey, "apiKey", config.Send.Authentication.APIKey)
	}

	if config.Send.Authentication.APISecret != "apiSecret" {
		t.Fatalf("Expected %q = %q. Result %q", envAPISecret, "apiSecret", config.Send.Authentication.APISecret)
	}

	if config.Send.Authentication.OAuth1.Token != "userToken" {
		t.Fatalf("Expected %q = %q. Result %q", envOAuth1Token, "userToken", config.Send.Authentication.OAuth1.Token)
	}

	if config.Send.Authentication.OAuth1.Secret != "userSecret" {
		t.Fatalf("Expected %q = %q. Result %q", envOAuth1Secret, "userSecret", config.Send.Authentication.OAuth1.Secret)
	}

}
