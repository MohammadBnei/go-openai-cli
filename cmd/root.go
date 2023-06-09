/*
Copyright © 2023 Mohammad-Amine BANAEI mohammadamine.banaei@pm.me

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
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "go-openai-cli",
	Short: "Go-OpenAI-CLI is a command-line interface that allows users to generate text using OpenAI's GPT-3 language generation service.",
	Long:  `Go-OpenAI-CLI is a command-line interface tool that provides users with convenient access to OpenAI's GPT-3 language generation service. With this app, users can easily send prompts to the OpenAI API and receive generated responses, which can then be printed on the command-line or saved to a markdown file. Go-OpenAI-CLI is an excellent tool for creatives, content creators, chatbot developers and virtual assistants, as they can use it to quickly generate text for various purposes. By configuring their OpenAI API key and model, users can customize the behavior of the app to suit their specific needs. Moreover, Go-OpenAI-CLI is an open-source project that welcomes contributions from the community, and it is licensed under the MIT License.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
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

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/config/go-openai-cli/config.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if os.Getenv("CONFIG") != "" {
		cfgFile = os.Getenv("CONFIG")
	}
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		viper.AddConfigPath(home + "/config/go-openai-cli")
		viper.SetConfigType("yaml")
		viper.SetConfigName("config")
	}

	err := viper.BindPFlag("OPENAI_KEY", rootCmd.Flags().Lookup("OPENAI_KEY"))
	if err != nil {
		fmt.Println(err)
		return
	}
	err = viper.BindPFlag("messages-length", rootCmd.Flags().Lookup("messages-length"))
	if err != nil {
		fmt.Println(err)
		return
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
