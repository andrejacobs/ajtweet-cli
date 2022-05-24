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
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	sendDryRunFlag bool
)

// sendCmd represents the send command
var sendCmd = &cobra.Command{
	Use:   "send",
	Short: "Send the scheduled tweets to Twitter",
	Long: `Send the scheduled tweets to Twitter

The list of scheduled tweets will be checked against the current time to
determine which tweets need to be sent as soon as possible. Only the tweets
that have a scheduled time that is before the current system time will be
sent.

You may also simulate the send process by running the command
in the dry run mode (-n, --dry-run).

Rate limiting:
 You can limit the number of tweets that will be sent per run of the app
 by specifying the send.max value in the configuration file.

 You can also specify a delay in seconds that the app needs to wait after
 each tweet was sent using the send.delay value in the configuration file.

 Example YAML to send a maximum number of 100 tweets and delay 5 seconds
 between each sending of a tweet.

    send:
        max: 100
        delay: 5

Authentication:
 Please see the Authentication section from the root command's
 help on how to configure the required authentication needed to use the
 Twitter APIs (i.e. ajtweet --help).

Examples:
 ajtweet send
 ajtweet send --dry-run	
`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := application.Send(os.Stdout, sendDryRunFlag); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to send. Error: %s\n", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(sendCmd)

	sendCmd.Flags().BoolVarP(&sendDryRunFlag, "dry-run", "n", false, "Tweets will not be sent to Twitter and also not be deleted")
}
