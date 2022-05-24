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

// cmd provides the CLI interface for the ajtweet application.
package cmd

import (
	"fmt"
	"os"

	"github.com/andrejacobs/ajtweet-cli/app"
	"github.com/andrejacobs/ajtweet-cli/internal/buildinfo"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

var application app.Application

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "ajtweet",
	Version: buildinfo.VersionString(),
	Short:   "Schedule tweets to be sent to Twitter",
	Long: `Schedule tweets to be sent to Twitter.

Tweets are scheduled using the "ajtweet add" command and will be sent
according to the preferred schedule when "ajtweet send" command is run.

Configuration:
 ajtweet will look for a configuration file named ".ajtweet" and a 
 supported extension (.yaml, .toml, .ini) in the following directories
 (in this specified order):
     ./             Current working directory
     $HOME/         User's home directory
     /etc/ajtweet

 For example: The configuration file $HOME/.ajtweet.ini will be found and use
 before the file /etc/ajtweet/.ajtweet.yaml

 --config path
    Can be used to explicitly specify the configuration file to be used.

TODO: Need to document the available config values. Maybe have a command to generate an example.

Examples:

 ajtweet add "Send this tweet asap"
 ajtweet add --scheduledAt "2022-05-23T21:22:42Z" "Send this later"

 date | xargs -0 ajtweet add
    Pass the output from date as the message argument expected by add.

 ajtweet list
 ajtweet list --json

 ajtweet delete "a2fdb340-0b61-4a89-b52e-82deae2e3aa8"
 ajtweet delete --dry-run "a2fdb340-0b61-4a89-b52e-82deae2e3aa8"
 ajtweet delete --all
`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Persistent flags that are available to every subcommand
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.ajtweet.yaml)")

	versionTemplate := `{{printf "%s: %s - %s\n" .Name .Short .Version}}`
	rootCmd.SetVersionTemplate(versionTemplate)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search for a config file with base name ".ajtweet" (extension will determine the type of format used, e.g. yaml)
		// starting at the current working directory
		viper.AddConfigPath(".")
		// then check the $HOME directory
		viper.AddConfigPath(home)
		// finally check in /etc/ajtweet
		viper.AddConfigPath(fmt.Sprintf("/etc/%s/", rootCmd.Name()))

		viper.SetConfigType("yaml")
		viper.SetConfigName(".ajtweet")

		//NOTE: Viper supports the following formats: JSON, TOML, YAML, HCL (HashiCorp), INI, envfile or Java properties
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading config file: %s. Error: %s\n", viper.ConfigFileUsed(), err)
		os.Exit(1)
	}

	initApplication()
}

// Initialize the main Application "context" used by the CLI commands.
func initApplication() {
	var appConfig app.Config
	if err := viper.Unmarshal(&appConfig); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing the configuration: %s", err)
		os.Exit(1)
	}

	appConfig.PopulateFromEnv()

	if appConfig.Datastore.Filepath == "" {
		appConfig.Datastore.Filepath = "./ajtweets-data.json"
	}

	if err := application.Configure(appConfig); err != nil {
		fmt.Fprintf(os.Stderr, "Error configuring the application: %s", err)
		os.Exit(1)
	}
}
