package main

import (
	"fmt"

	"github.com/sashabaranov/go-openai"
)

func main() {

	// usage
	// help, model, token and bill info
	// only one line question - 마지막 글자가 Ctrl+D이면 끝내게 하기?
	fmt.Println("todo: usage")
	fmt.Println()

	for {
		// get user input
		userInput := getUserInput(reader)

		// TODO: command execution
		// help, msg(show), clear, sysmsg(get, set), config
		excuted := commandExecute(userInput)
		if excuted {
			continue
		}

		// make messages
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleUser,
			Content: userInput,
		})
		// TODO: if too long, remove from the oldest couple, except the system message

		resp, err := getResponse(messages)
		if err != nil {
			fmt.Println(colorStr(Red, fmt.Sprintf("ChatCompletion error: %v", err)))
			fmt.Println()
			continue
		}

		content := resp.Choices[0].Message.Content
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleAssistant,
			Content: content,
		})
		fmt.Println(content)
		fmt.Println()
	}
}
