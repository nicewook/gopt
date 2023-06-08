package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/sashabaranov/go-openai"
)

func colorStr(color, msg string) string {
	return fmt.Sprint(color + msg + Reset)
}

func getUserInput(reader *bufio.Reader) string {
	fmt.Print(colorStr(Cyan, "gpt> "))
	userInput, _ := reader.ReadString('\n')

	userInput = strings.Replace(userInput, "\n", "", -1)
	return userInput
}

// getResponse send request and get response from the OpenAI
// it uses 'gpt-3.5-turbo'
func getResponse(messages []openai.ChatCompletionMessage) (openai.ChatCompletionResponse, error) {
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:    openai.GPT3Dot5Turbo,
			Messages: messages,
		},
	)
	return resp, err
}

func commandExecute(input string) bool {
	switch input {
	case "help":
		fmt.Println("help statement")
	case "history":
		fmt.Println("all chatting history:")
		for _, m := range messages {
			fmt.Println(m)
		}
	case "config":
		fmt.Println("config statement")
	case "exit":
		fmt.Println("exit statement")
		os.Exit(0)
	default:
		return false
	}
	return true
}
