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
	"errors"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/spf13/cobra"
)

var (
	dryRunFlag    bool
	deleteAllFlag bool
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete scheduled tweets",
	Long: `Delete scheduled tweets

Each argument must match a tweet identifier in the scheduled list.

You may also simulate the deletion process by running the command
in the dry run mode (-n, --dry-run).

Examples:

 ajtweet delete "28cf75a1-e7b3-4401-a878-4362bdc4befe" "a2fdb340-0b61-4a89-b52e-82deae2e3aa8"
    Delete two tweets with the specified identifiers.

 ajtweet delete --dry-run "28cf75a1-e7b3-4401-a878-4362bdc4befe"
    Simulate a delete by running in dry run mode.

 ajtweet delete --all
    Delete all the scheduled tweets.
`,
	Args: func(cmd *cobra.Command, args []string) error {
		argCount := len(args)
		if deleteAllFlag {
			if argCount != 0 {
				return errors.New("--all Does not expect arguments to be passed")
			}
		} else if argCount < 1 {
			return errors.New("expected identifiers to be passed as arguments")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		if deleteAllFlag {
			// Confirm with the user
			expectedConfirm := randomString(5 + rand.Intn(5))
			fmt.Fprintf(os.Stdout, "Please confirm by entering: %s\n> ", expectedConfirm)

			var confirm string
			if _, err := fmt.Scanln(&confirm); err != nil {
				fmt.Fprintf(os.Stderr, "Failed to receive user confirmation. Error: %s\n", err)
				os.Exit(1)
			}

			if confirm != expectedConfirm {
				fmt.Fprintf(os.Stderr, "Confirmation failed\n")
				os.Exit(1)
			}

			if err := application.DeleteAll(); err != nil {
				fmt.Fprintf(os.Stderr, "Failed to delete all the tweets. Error: %s\n", err)
				os.Exit(1)
			}
		} else {
			for _, idString := range args {
				if err := application.Delete(idString); err != nil {
					fmt.Fprintf(os.Stderr, "Failed to delete the tweet with identifier: %q. Error: %s\n", idString, err)
					os.Exit(1)
				}

				fmt.Fprintf(os.Stdout, "Deleting tweet with identifier: %q\n", idString)
			}
		}

		if !dryRunFlag {
			if err := application.Save(); err != nil {
				fmt.Fprintf(os.Stderr, "Failed to save the changes. Error: %s\n", err)
				os.Exit(2)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)

	deleteCmd.Flags().BoolVarP(&dryRunFlag, "dry-run", "n", false, "Tweets will not be deleted")
	deleteCmd.Flags().BoolVarP(&deleteAllFlag, "all", "a", false, "Delete all the scheduled tweets")

	rand.Seed(time.Now().UnixNano())
}

var letterRunes = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randomString(count int) string {
	result := make([]rune, count)
	letterRunesCount := len(letterRunes)
	for i := range result {
		result[i] = letterRunes[rand.Intn(letterRunesCount)]
	}
	return string(result)
}
