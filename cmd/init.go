package cmd

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/chzyer/readline"
	"github.com/sashabaranov/go-openai"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/nicewook/gopt/config"
)

const ( // color for width calcluation. many of them not used but reserve
	Reset         = "\033[0m"
	Red           = "\033[31m"
	Green         = "\033[32m"
	Yellow        = "\033[33m"
	Blue          = "\033[34m"
	Purple        = "\033[35m"
	Cyan          = "\033[36m"
	Gray          = "\033[37m"
	White         = "\033[97m"
	lenColorCode  = len(Red)
	lenColorReset = len(Reset)
	lenColor      = lenColorCode + lenColorReset
)

func rightAlignWithColorWords(wordCount int) string {
	// color prompt should be counted for right alignment
	virtualWidth := wordCount*lenColor + readline.GetScreenWidth()
	return "%" + strconv.Itoa(virtualWidth) + "v\n"
}

var (
	client        *openai.Client
	messages      []openai.ChatCompletionMessage
	systemMessage = openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleSystem,
		Content: config.DefaultSystemMsg,
	}
)

func init() {

	setLog()
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.gopt.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func setLog() {
	runMode := os.Getenv("RUN_MODE")
	if runMode != "dev" {
		log.SetOutput(io.Discard)
	}
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {

	home, err := os.UserHomeDir()
	cobra.CheckErr(err)
	configPath := filepath.Join(home, ".local", "gopt")

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// Directory does not exist, create it
		err := os.Mkdir(configPath, 0755)
		if err != nil {
			log.Fatal("Error creating directory:", err)
		}
	}

	viper.AddConfigPath(configPath)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	// If a config file is found, read it in, or set default value
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; create a default one.
			viper.SetDefault("model", config.DefaultModel)
			viper.SetDefault("token", false)
			viper.SetDefault("system_message", config.DefaultSystemMsg)
			viper.SetDefault("max_token", config.DefaultMaxToken)
			err := viper.SafeWriteConfigAs(filepath.Join(configPath, "config.yaml"))
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
		viper.Set("openai_api_key", apiKey)
		err := viper.WriteConfig()
		cobra.CheckErr(err)
	}

	// init variables
	client = openai.NewClient(viper.GetString("openai_api_key"))
	messages = make([]openai.ChatCompletionMessage, 0, 10)
	messages = append(messages, systemMessage)
}
