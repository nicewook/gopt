package cmd

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/chzyer/readline"
	"github.com/fatih/color"
	"github.com/sashabaranov/go-openai"
)

// Model	     Input	              Output
// 4K context	 $0.0015 / 1K tokens	$0.002 / 1K tokens
// 16K context $0.003  / 1K tokens	$0.004 / 1K tokens
const perInputToken = float32(0.0000015)
const perOutputToken = float32(0.000002)

// chatComplete send request and get response from the OpenAI
func chatComplete(model string, messages []openai.ChatCompletionMessage) (openai.ChatCompletionResponse, error) {
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:     model,
			Messages:  messages,
			MaxTokens: ModelMaxCompletionToken,
		},
	)

	return resp, err
}

// chatCompleteStream send request and get response as stream from the OpenAI
func chatCompleteStream(model string, messages []openai.ChatCompletionMessage) (*openai.ChatCompletionStream, error) {
	stream, err := client.CreateChatCompletionStream(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:     model,
			Messages:  messages,
			MaxTokens: ModelMaxCompletionToken,
			Stream:    true,
		},
	)
	return stream, err
}

func calcPrice(iTokens, oTokens int) float32 {
	return float32(iTokens)*perInputToken + float32(oTokens)*perOutputToken
}

func prepareElapsedTime(eTime time.Duration) string {
	green := color.New(color.FgGreen).SprintfFunc()
	return fmt.Sprintf("elapsed time: %s", green("%.2fms", eTime.Seconds()))
}

func prepareTokenInfo(u openai.Usage) string {
	greenF := color.New(color.FgGreen).SprintfFunc()
	return fmt.Sprintf("Token. prompt: %s, completion: %s, total: %s, money spent: %s",
		greenF("%d", u.PromptTokens),
		greenF("%d", u.CompletionTokens),
		greenF("%d", u.TotalTokens),
		greenF("$%.6f", calcPrice(u.PromptTokens, u.CompletionTokens)),
	)
}

func prepareCumulativeTokenInfo(totalPromptTokens, totalCompletionTokens int) string {
	greenF := color.New(color.FgGreen).SprintfFunc()

	return fmt.Sprintf("cumulative total: %s, money spent: %s",
		greenF("%d", totalPromptTokens+totalCompletionTokens),
		greenF("$%.6f", calcPrice(totalPromptTokens, totalCompletionTokens)),
	)
}

// command
func helpMessage() string {
	green := color.New(color.FgGreen).SprintFunc()

	help := green("help")
	config := green("config")
	context := green("context")
	reset := green("reset")
	clear := green("clear")
	exit := green("exit")
	q := green("q")

	return fmt.Sprintf(`Usage:
  - %s - Displays this help message.
  - %s - Displays configuration information. 
  - %s - Displays the conversation context which reserved at the moment.
  - %s - Reset all the conversation context.
  - %s - Clear terminal.
  - %s or %s - Exits the app.
	`, help, config, context, reset, clear, exit, q)
}

func commandExecute(w io.Writer, input string) bool {

	switch input {
	case "":
		fmt.Println()
		log.Println("do nothing")
	case "h":
		fallthrough
	case "help":
		fmt.Println(helpMessage())

	case "config":
		fmt.Println()
		color.Green("not yet implemented.")
		// fmt.Println(green("not yet implemented."))
		fmt.Println()

	case "context":
		fmt.Println()
		if len(messages) == 0 {
			color.Green("no contexts yet.")
			fmt.Println()
			break
		}
		color.Green("all chatting context:")
		fmt.Println()
		for _, m := range messages {
			fmt.Println(m)
		}
		fmt.Println()
	case "reset":
		messages = []openai.ChatCompletionMessage{systemMessage}
		fmt.Println()
		color.Green("reset all the conversion context.")
		fmt.Println()

	case "clear":
		readline.ClearScreen(w)

	case "exit":
		fallthrough
	case "q":
		fmt.Println()
		color.Green("Have a great day!")
		os.Exit(0)

	default:
		return false
	}
	return true
}
