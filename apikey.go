package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type Config struct {
	APIKey string `json:"apikey"`
}

func configPath() (dirPath, filename string, err error) {
	var home string
	home, err = os.UserHomeDir()
	if err != nil {
		return
	}
	dirPath = filepath.Join(home, ".gopt")
	filename = "config.json"
	return
}

func getAPIKey() (string, error) {
	// apiKey from environment variable
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey != "" {
		return apiKey, nil
	}

	dirPath, filename, err := configPath()
	if err != nil {
		return apiKey, err
	}
	filePath := filepath.Join(dirPath, filename)

	b, err := os.ReadFile(filePath)
	if err != nil {
		return apiKey, err
	}

	var config Config
	if err := json.Unmarshal(b, &config); err != nil {
		return apiKey, err
	}

	apiKey = config.APIKey
	return apiKey, nil
}

func setAPIKey() string {
	fmt.Println("Fail to get OpenAI API Key. If you want to set the Key, type your key. Key will be saved on", colorStr(Red, "~/.gopt/config.json"))
	fmt.Print(colorStr(Cyan, "OpenAI API Key: "))
	userInput, _ := reader.ReadString('\n')
	userInput = strings.Replace(userInput, "\n", "", -1)

	// delete file, save file
	// marshal the Person struct to JSON
	config := Config{APIKey: userInput}
	jsonData, err := json.Marshal(config)
	if err != nil {
		log.Println("error:", err)
		os.Exit(0)
	}

	// create a new file to save the JSON data
	dirPath, filename, err := configPath()
	if err != nil {
		log.Println("error:", err)
		os.Exit(0)
	}
	if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
		log.Println("error:", err)
		os.Exit(0)
	}

	filePath := filepath.Join(dirPath, filename)
	file, err := os.Create(filePath)
	if err != nil {
		log.Println("error:", err)
		os.Exit(0)
	}
	defer file.Close()

	// write the JSON data to the file
	_, err = file.Write(jsonData)
	if err != nil {
		log.Println("error:", err)
		os.Exit(0)
	}

	// print a success message
	fmt.Println("Saved successfully, Be aware the API key is stored as plaintext.")

	return userInput
}

func APIKey() string {

	result, err := getAPIKey()
	if err == nil {
		return result
	}
	result = setAPIKey()
	return result

}
