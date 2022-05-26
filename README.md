# ajtweet-cli

A command line interface to schedule tweets to be sent to Twitter.

## Overview

ajtweet was started out of my need for scheduling tweets without going on to Twitter while also serving as a project to learn how to build Go CLI apps.

Please see my [spec](https://github.com/andrejacobs/specs/blob/main/tools/ajtweet/ajtweet%20-%20spec.md) for more insights.

## Prerequisites

You will need to have a registered developer account with Twitter along with the following keys, secrets and tokens:

* Consumer API key and secret.
* OAuth 1.0a User access token and secret.

Please see my [blog post](TODO Post my Twitter blog post and link here) on how you can get these values.

## Install

## Configuration

ajtweet requires a configuration file in order to function. Since I am using [Viper](https://github.com/spf13/viper) the configuration file can be in YAML, TOML, JSON etc. format.

The app will look for a configuration file named ".ajtweet" and a supported extension (.yaml, .toml, .ini etc.) in the following directories in this specified order:

      ./             Current working directory
      $HOME/         User's home directory
      /etc/ajtweet

For example: The configuration file `$HOME/.ajtweet.ini` will be found and used before the file `/etc/ajtweet/.ajtweet.yaml`.

You can also explicitly specify the configuration file to be used using the `--config` flag.

Environment variables can also be used to override some of the configuration values. For example the authentication values.

### Authentication

You will need to have a registered developer account with Twitter to be able to access the Twitter v2 APIs.

You will need the following:

* Consumer API Key
    - config path: send.authentication.api_key
    - environment: AJTWEET_API_KEY

* Consumer API Secret
    - config path: send.authentication.api_secret
    - environment: AJTWEET_API_SECRET

* OAuth 1.0 User access token
    - config path: send.authentication.oauth1.token
    - environment: AJTWEET_ACCESS_TOKEN

* OAuth 1.0 User access secret
    - config path: send.authentication.oauth1.secret
    - environment: AJTWEET_ACCESS_SECRET

See the file `env_example` in this repository for an example environment file that can be sourced in your shell session to set the required authentication values before running ajtweet.

        $ cp env_example .env
        $ source .env
        $ ajtweet send --dry-run

### Example YAML configuration

The following is an example YAML configuration file you can use to configure ajtweet. Name the file `.ajtweet.yaml` and store it in one of the search directories as mentioned earlier.

    datastore:
        filepath: /home/andre/ajtweet/tweets.json

    lockfile: /var/ajtweet/ajtweet.lock

    send:
        max: 100
        delay: 5

        authentication:
            api_key: your_consumer_key_for_twitter
            api_secret: your_consumer_secret

            oauth1:
                token: user_access_token
                secret: user_access_secret

* datastore: Specifies the file to be used for storing the schedued tweets. At the moment only a JSON file is supported. Please ensure the directories exist before running the app.
* lockfile: The path of where the lock file will be created. The default is ./ajtweet.lock.
* send.max: The maximum number of tweets to be sent during a call to the `send` command. Default value is 10.
* send.delay: The time in seconds to wait after each tweet before sending the next one. Default value is 1 second.
* authentication: Specify the Twitter API key and secret along with the OAuth 1.0a user access token and secret. Please note that you can use environment variables instead as mentioned in the Authentication section.

## Add tweets

You schedule tweets using the `add` command. Tweets will only be sent to Twitter when you run the `send` command.

The RFC3339 time format is used for specifying a preferred scheduled time. That is a time at which you would like the tweet to be sent, however the actual time at which a tweet is sent is determined by when the `send` command is run.

Example RFC3339 format: YYYY-MM-DDTHH:mm:ssZ, e.g. 2022-05-16T19:39Z

The preferred time is specified using the `-t` or `--scheduledAt` flag followed by the time in the RFC3339 format. The default value is the time at which you ran the `add` command. Thus it will be sent as soon as `send` is run.

* Schedule a tweet to be sent as soon as possible (i.e. next time `ajtweet send` is run).

        ajtweet add "Please send this tweet as soon as you can"

* Schedule a tweet to be sent no earlier than a specified time.

        ajtweet add --scheduledAt "2032-05-16T19:42:00Z" "Send this tweet a year from now"

## List tweets

Run the `list` command to see the list of scheduled tweets that still need to be sent.

The list of tweets will be displayed using ANSI colour escape codes if the STDOUT supports colour output. For example output to the terminal will use colour output whereas piping to another command will not use colours. Colour output can also be disabled by setting the environment variable NO_COLOR.

The list of tweets can also be displayed as an encoded JSON string using the `-j` or `--json` flags.

* Display all scheduled tweets.

        ajtweet list

        id: 8b957daf-9967-4bc2-b123-f184e0079afe
        time: 2022-05-24T20:55:07+01:00 [send now!]
        tweet: Please send this tweet as soon as you can

        id: 9896a759-77c8-434a-9759-81dccfacbb1b
        time: 2022-05-24T21:55:00Z
        tweet: Hello world

* Display all scheduled tweets as a JSON encoding.

         ajtweet list --json

         [{"id":"8b957daf-9967-4bc2-b123-f184e0079afe","message":"Please send this tweet as soon as you can","scheduledTime":"2022-05-24T20:55:07+01:00"},{"id":"9896a759-77c8-434a-9759-81dccfacbb1b","message":"Hello world","scheduledTime":"2022-05-24T21:55:00Z"}]

* Display all scheduled tweets while disabling colour output on STDOUT.

        NO_COLOR=1 ajtweet list

## Delete tweets

Tweets are uniquely identified by an identifier and you will need to pass this to the `delete` command. The identifiers can be found by using the `list` command.

You may also simulate the deletion process by running the command in the dry run mode using the `-n` or `--dry-run` flag.

To delete all the scheduled tweet you can use the `-a` or `--all` flag. You will be prompted to confirm by repeating a random string.

* Delete two tweets matching the specified identifiers.

        ajtweet delete "28cf75a1-e7b3-4401-a878-4362bdc4befe" "a2fdb340-0b61-4a89-b52e-82deae2e3aa8"

        Deleting tweet with identifier: "28cf75a1-e7b3-4401-a878-4362bdc4befe"
        Deleting tweet with identifier: "a2fdb340-0b61-4a89-b52e-82deae2e3aa8"

* Simulate deleting a tweet using the dry-run mode.

        ajtweet delete --dry-run "28cf75a1-e7b3-4401-a878-4362bdc4befe"

        Deleting tweet with identifier: "9896a759-77c8-434a-9759-81dccfacbb1b"

* Delete all the tweets.

        ajtweet delete --all

        Please confirm by entering: aR1ssKS3
        >

## Send tweets to Twitter

Tweets will only be sent when you run the `send` command.

The list of scheduled tweets will be checked against the current time to determine which tweets need to be sent as soon as possible. Only the tweets that have a preferred scheduled time that is before the current system time will be considered for sending to Twitter.

You may also simulate the send process by running the command in the dry run mode by using the `-n` or `--dry-run` flag.

You can limit the number of tweets that will be sent per run of the app by specifying the send.max value in the configuration file. The default value is 10.

You can also specify a delay in seconds that the app needs to wait after each tweet was sent using the send.delay value in the configuration file.

* Send all scheduled tweets that have been scheduled to be sent before the current system time.

        ajtweet send

        Sending 1 of 2
        id: 4a5884b0-a0ca-4b4e-ab6a-e6ae43b7b8bc
        tweet: Hello

        Delaying for 5 seconds ...

        Sending 2 of 2
        id: f2a18d0b-1ec2-4d2a-a4ec-eeba52bb78ef
        tweet: World

* Simulate sending all the scheduled tweets.

        ajtweet send --dry-run

## Single allowed instance

Only one instance of ajtweet is allowed to run at any one point in time. This is to ensure that only one program is making changes to the data store or sending tweets.

The app uses a lock file to ensure only one instance of the program is running. The lock file's path can be specified in the configuration file. The default is `./ajtweet.lock`

## Schedule ajtweet

## Dependencies

ajtweet uses the following excellent packages:

* spf13's [Cobra](https://github.com/spf13/cobra)
* spf13's [Viper](https://github.com/spf13/viper)
* michimani's [gotwi](https://github.com/michimani/gotwi)
* Fatih's [color](https://github.com/fatih/color)

## License

ajtweet is released under the MIT license. See [LICENSE](LICENSE) for details.
