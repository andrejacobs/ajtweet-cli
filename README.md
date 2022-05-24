# ajtweet-cli

# TODO:
- Need to set a default for the send.max value, because nothing will be sent if the user doesn't specify this value.
- Need to implement a lock file while send is run, so that add, delete can't be used.

## Overview

TODO: Post my Twitter blog post and link here
TODO: Reference the spec

## Install

## Configuration

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

## Schedule ajtweet

## License

ajtweet is released under the MIT license. See [LICENSE](LICENSE) for details.
