package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/sashabaranov/go-openai"
)

// $0.002 dollar for 1K tokens
// $0.000002 dollar for 1 tokens
const perToken = float32(0.000002)

func colorStr(color, msg string) string {
	return color + msg + Reset
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

func prepareTokenInfo(u openai.Usage) string {

	return fmt.Sprintf("Token. prompt: %s, completion: %s, total: %s, money spent: %s",
		colorStr(Green, fmt.Sprintf("%d", u.PromptTokens)),
		colorStr(Green, fmt.Sprintf("%d", u.CompletionTokens)),
		colorStr(Green, fmt.Sprintf("%d", u.TotalTokens)),
		colorStr(Green, fmt.Sprintf("$%.6f", float32(u.TotalTokens)*perToken)))
}

func prepareCumulativeTokenInfo(totalToken int) string {

	return fmt.Sprintf("cumulative total: %s, money spent: %s",
		colorStr(Green, fmt.Sprintf("%d", totalToken)),
		colorStr(Green, fmt.Sprintf("$%.6f", float32(totalToken)*perToken)))
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
		fallthrough
	case "q":
		fmt.Println("exit statement")
		os.Exit(0)
	default:
		return false
	}
	return true
}
