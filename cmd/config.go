/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"strconv"

	"github.com/nicewook/gopt/internal/gopt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Config your gopt. model, token usage display, etc.",
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("config called")
	},
}

func init() {
	configCmd.AddCommand(getCmd, setCmd, listCmd, resetCmd)
	rootCmd.AddCommand(configCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// configCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// configCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

var getCmd = &cobra.Command{
	Use:   "get config",
	Short: "Get a value from the configuration",
	Args:  cobra.ExactArgs(1),
	Run:   getConfig,
}

func getConfig(cmd *cobra.Command, args []string) {
	key := args[0]
	if !viper.IsSet(key) {
		fmt.Printf("Key %s is not valid in the configuration.\n", key)
		return
	}

	value := viper.Get(key)
	fmt.Printf("%s: %v\n", key, value)
}

var setCmd = &cobra.Command{
	Use:   "set",
	Short: "Set a value in the configuration",
	Args:  cobra.ExactArgs(2),
	Run:   setConfig,
}

func setConfig(cmd *cobra.Command, args []string) {
	key, value := args[0], args[1]
	if !viper.IsSet(key) {
		fmt.Printf("Key %s is not valid in the configuration.\n", key)
		return
	}

	// model values validation
	if key == "model" {

		var isValid bool
		for _, model := range gopt.GetModels() {
			if value == model {
				isValid = true
				break
			}
		}

		if !isValid {
			fmt.Println("Value for key 'model' is not valid. Refer to available model list")
			for _, model := range gopt.GetModels() {
				fmt.Println("-", model)
			}
			return
		}

		// change model

	}

	// Let's try to convert the value string to int
	if intValue, err := strconv.Atoi(value); err == nil {
		viper.Set(key, intValue)
	} else if boolValue, err := strconv.ParseBool(value); err == nil {
		// The value is not an int, let's try with bool
		viper.Set(key, boolValue)
	} else {
		// The value is neither int nor bool, so we keep it as string
		viper.Set(key, value)
	}

	err := viper.WriteConfig()
	if err != nil {
		fmt.Println("Error writing config:", err)
	}
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List the configuration",
	Run:   listConfig,
}

func listConfig(cmd *cobra.Command, args []string) {
	settings := viper.AllSettings()
	for key, value := range settings {
		if key != "openai_api_key" {
			fmt.Printf("%s: %v\n", key, value)
		}
	}
}

var resetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Reset the configuration",
	Run:   resetConfig,
}

func resetConfig(cmd *cobra.Command, args []string) {
	viper.Set("model", gopt.DefaultModel)
	viper.Set("token", false)
	viper.Set("system_message", gopt.DefaultSystemMsg)
	viper.Set("max_token", gopt.DefaultMaxToken)
	err := viper.WriteConfig()
	if err != nil {
		fmt.Println("Error writing config:", err)
	}

	fmt.Println("config reset to default.")
	fmt.Println("--")
	listConfig(cmd, args)
}
