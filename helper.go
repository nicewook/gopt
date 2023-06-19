package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/sashabaranov/go-openai"
)

// Model	     Input	              Output
// 4K context	 $0.0015 / 1K tokens	$0.002 / 1K tokens
// 16K context $0.003  / 1K tokens	$0.004 / 1K tokens
const perInputToken = float32(0.0000015)
const perOutputToken = float32(0.000002)

func colorStr(color, msg string) string {
	return color + msg + Reset
}

// chatComplete send request and get response from the OpenAI
// it uses 'gpt-3.5-turbo'
func chatComplete(messages []openai.ChatCompletionMessage) (openai.ChatCompletionResponse, error) {
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:     "gpt-3.5-turbo-0613",
			Messages:  messages,
			MaxTokens: ModelMaxCompletionToken,
		},
	)
	return resp, err
}

func calcPrice(iTokens, oTokens int) float32 {
	return float32(iTokens)*perInputToken + float32(oTokens)*perOutputToken
}

func prepareElapsedTime(eTime time.Duration) string {

	return fmt.Sprintf("elapsed time: %s",
		colorStr(Green, fmt.Sprintf("%.2fms", eTime.Seconds())),
	)
}

func prepareTokenInfo(u openai.Usage) string {

	return fmt.Sprintf("Token. prompt: %s, completion: %s, total: %s, money spent: %s",
		colorStr(Green, fmt.Sprintf("%d", u.PromptTokens)),
		colorStr(Green, fmt.Sprintf("%d", u.CompletionTokens)),
		colorStr(Green, fmt.Sprintf("%d", u.TotalTokens)),
		colorStr(Green, fmt.Sprintf("$%.6f", calcPrice(u.PromptTokens, u.CompletionTokens))),
	)
}

func prepareCumulativeTokenInfo(totalPromptTokens, totalCompletionTokens int) string {

	return fmt.Sprintf("cumulative total: %s, money spent: %s",
		colorStr(Green, fmt.Sprintf("%d", totalPromptTokens+totalCompletionTokens)),
		colorStr(Green, fmt.Sprintf("$%.6f", calcPrice(totalPromptTokens, totalCompletionTokens))),
	)
}

// command
func helpMessage() string {
	help := colorStr(Green, "help")
	config := colorStr(Green, "config")
	context := colorStr(Green, "context")
	reset := colorStr(Green, "reset")
	exit := colorStr(Green, "exit")
	q := colorStr(Green, "q")

	return fmt.Sprintf(`Usage:
  - %s - Displays this help message.
  - %s - Displays configuration information. 
  - %s - Displays the conversation context which reserved at the moment.
  - %s - Reset all the conversation context.
  - %s or %s - Exits the app.
	`, help, config, context, reset, exit, q)
}

func commandExecute(input string) bool {
	switch input {
	case "":
		fallthrough
	case "help":
		fmt.Println(helpMessage())

	case "config":
		fmt.Println()
		fmt.Println(colorStr(Green, "not yet implemented."))
		fmt.Println()

	case "context":
		fmt.Println()
		if len(messages) == 0 {
			fmt.Println(colorStr(Green, "no contexts yet."))
			fmt.Println()
			break
		}
		fmt.Println(colorStr(Green, "all chatting context:"))
		fmt.Println()
		for _, m := range messages {
			fmt.Println(m)
		}
		fmt.Println()
	case "reset":
		messages = []openai.ChatCompletionMessage{systemMessage}
		fmt.Println()
		fmt.Println(colorStr(Green, "reset all the conversion context."))
		fmt.Println()

	case "exit":
		fallthrough
	case "q":
		fmt.Println()
		fmt.Println(colorStr(Green, "Have a great day!"))
		os.Exit(0)

	default:
		return false
	}
	return true
}
