/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Config your gopt. model, token usage display, etc.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("config called")
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
	Run: getConfig,
}

func getConfig(cmd *cobra.Command, args []string) {
	key := args[0]
	value := viper.Get(key)
	fmt.Printf("%s: %v\n", key, value)
}

var setCmd = &cobra.Command{
	Use:   "set",
	Short: "Set a value in the configuration",
	Args:  cobra.ExactArgs(2),
	Run: setConfig,
}

func setConfig(cmd *cobra.Command, args []string) {
	key := args[0]
	value := args[1]
	viper.Set(key, value)
	err := viper.WriteConfig()
	if err != nil {
		fmt.Println("Error writing config:", err)
	}
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List the configuration",
	Run: listConfig,
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
	Run: resetConfig,
}

func resetConfig(cmd *cobra.Command, args []string) {
	viper.Set("model", "")
	viper.Set("token", "")
	viper.Set("max_token", 0)
	viper.Set("system_message", "")
	err := viper.WriteConfig()
	if err != nil {
		fmt.Println("Error writing config:", err)
	}
}
