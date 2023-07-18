/*
Copyright © 2023 Hyunseok Jeong nicewook@hotmail.com

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
	"log"
	"os"
	"path/filepath"

	"github.com/nicewook/gopt/internal/gopt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gopt",
	Short: "gopt is OpenAI GPT powered chat completion CLI",
	Long: `gopt is OpenAI GPT powered chat completion CLI.
You can ask questions to the GPT model and get answers.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("start")
	},
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

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.gopt.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {

	home, err := os.UserHomeDir()
	cobra.CheckErr(err)
	configPath := filepath.Join(home, ".local", "gopt")
	fmt.Println(configPath)
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// Directory does not exist, create it
		err := os.Mkdir(configPath, 0755)
		if err != nil {
						fmt.Println("Error creating directory:", err)
		} else {
						fmt.Println("Directory created successfully")
		}
} 

	viper.AddConfigPath(configPath)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	// viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in, or set default value
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; create a default one.
			// viper.SetDefault("apikey", "your-api-key")
			viper.SetDefault("model", gopt.DefaultModel)
			viper.SetDefault("token", false)
			viper.SetDefault("system_message", gopt.DefaultSystemMsg)
			viper.SetDefault("max_token", gopt.DefaultMaxToken)
			err := viper.SafeWriteConfigAs(filepath.Join(configPath, gopt.ConfigFile))
			cobra.CheckErr(err)

		} else {
			// Config file was found but another error was produced
			log.Fatalf("Fatal error config file: %v", err)
		}
	}

	if !viper.IsSet("openai_api_key") {
		fmt.Print("Enter your OpenAI API key: ")
		var apiKey string
		fmt.Scanln(&apiKey)
		viper.Set("apiKey", apiKey)
		err := viper.WriteConfig()
		cobra.CheckErr(err)
	}
}
