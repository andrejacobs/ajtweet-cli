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
	"io"
	"os"

	"github.com/spf13/cobra"
)

var (
	jsonFlag bool
)

type listFunc func(io.Writer) error

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Display a list of the scheduled tweets.",
	Long: `Display a list of the scheduled tweets that still need to be sent.

The list of tweets will be displayed using ANSI colour escape codes if the
stdout supports colour output. For example output to the terminal will use
colour output whereas piping to another command will not use colours.
Colour output can also be disabled by setting the environment variable
NO_COLOR.

-j, --json Can be used to output the list into a JSON format.

Examples:

ajtweet list
	List all the tweets still to be sent.

ajtweet list --json
	List all the tweets in JSON output.

NO_COLOR=1 ajtweet list
	Disable colour output while displaying the list.

ajtweet list | less
	Disable colour output and pipe the output to less.
`,
	Run: func(cmd *cobra.Command, args []string) {

		var handler listFunc = application.List

		if jsonFlag {
			handler = application.ListJSON
		}

		if err := handler(os.Stdout); err != nil {
			fmt.Fprint(os.Stderr, err)
			os.Exit(1)
		}

		fmt.Fprintln(os.Stdout)
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	listCmd.Flags().BoolVarP(&jsonFlag, "json", "j", false, "Output the list into JSON format")
}
