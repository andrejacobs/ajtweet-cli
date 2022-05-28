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
	"time"

	"github.com/spf13/cobra"
)

var (
	scheduledAtFlag string
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new tweet to be sent to Twitter",
	Long: `Add a new tweet to be sent to Twitter at a preferred scheduled time.

-t, --scheduledAt specifies the preferred time in the RFC3339 format at which
you would like the tweet to be sent at. If no time is specified then the
current time will be used.
NOTE: The actual time the tweet will be sent depends entirely on when the send
command is run. Hence why this is the preferred time and not "guaranteed time".

Example RFC3339 format: YYYY-MM-DDTHH:mm:ssZ, e.g. 2022-05-16T19:39Z
	
Tweets are stored as per the application's configuration. Please see the 
main help section for more details (ajtweet help)

Examples:

 ajtweet add "Please send this tweet as soon as you can"
    Add a tweet to be sent the next time ajtweet send is run.

 ajtweet add --scheduledAt "2032-05-16T19:42:00Z" "Send this tweet a year from now"
    Add a tweet to be sent at the preferred scheduled time.
	`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		if scheduledAtFlag == "" {
			scheduledAtFlag = time.Now().Format(time.RFC3339)
		}

		if err := application.Add(args[0], scheduledAtFlag); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to add tweet. Error: %s\n", err)
			cleanupAndExit(1)
		}

		if err := application.Save(); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to save the changes. Error: %s\n", err)
			cleanupAndExit(2)
		}
	},
}

func init() {
	rootCmd.AddCommand(addCmd)

	addCmd.Flags().StringVarP(&scheduledAtFlag, "scheduledAt", "t", "", "Scheduled date time according to RFC3339 standard")
}
