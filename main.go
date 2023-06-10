package main

import (
	"fmt"
	"log"
	"strconv"

	tsize "github.com/kopoli/go-terminal-size"
	"github.com/sashabaranov/go-openai"
)

func main() {

	// usage
	// help, model, token and bill info
	// only one line question - 마지막 글자가 Ctrl+D이면 끝내게 하기?
	fmt.Println("todo: usage")
	fmt.Println()

	var totalToken int
	for {
		// get user input
		userInput := getUserInput(reader)

		// TODO: reserved command execution
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
		totalToken += resp.Usage.TotalTokens
		tokenInfo := prepareTokenInfo(resp.Usage)
		cumulativeTokenInfo := prepareCumulativeTokenInfo(totalToken)

		messages = append(messages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleAssistant,
			Content: content,
		})
		fmt.Println(content)

		// print tokenInfo
		s, err := tsize.GetSize()
		if err == nil {
			log.Println("Current size is", s.Width, "by", s.Height)
		}
		fmtString := "%" + strconv.Itoa(s.Width) + "v\n"
		log.Println("format string:", fmtString)
		fmt.Printf(fmtString, tokenInfo)
		fmt.Printf(fmtString, cumulativeTokenInfo)
	}
}
