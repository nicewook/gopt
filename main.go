package main

import (
	"fmt"
	"log"
	"strconv"
	"time"

	tsize "github.com/kopoli/go-terminal-size"
	"github.com/sashabaranov/go-openai"
)

func main() {

	// usage
	// help, model, token and bill info
	// only one line question - 마지막 글자가 Ctrl+D이면 끝내게 하기?
	fmt.Println()
	fmt.Println(helpMessage())

	var totalPromptTokens, totalCompletionTokens int
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
		// length check and adjust. except the system message


		// request completion 
		sTime := time.Now()
		resp, err := getResponse(messages)
		if err != nil {
			fmt.Println(colorStr(Red, fmt.Sprintf("ChatCompletion error: %v", err)))
			fmt.Println()
			continue
		}
		eTime := time.Since(sTime)

		content := resp.Choices[0].Message.Content
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleAssistant,
			Content: content,
		})
		fmt.Println(content)

		// print elapsed time, tokenInfo
		totalPromptTokens += resp.Usage.PromptTokens
		totalCompletionTokens += resp.Usage.CompletionTokens
		elapsedTime := prepareElapsedTime(eTime)
		tokenInfo := prepareTokenInfo(resp.Usage)
		cumulativeTokenInfo := prepareCumulativeTokenInfo(totalPromptTokens, totalCompletionTokens)

		s, err := tsize.GetSize()
		if err == nil {
			log.Println("Current size is", s.Width, "by", s.Height)
		}
		fmtStr1Color := "%" + strconv.Itoa(s.Width+1*lenColor) + "v\n"
		fmtStr4Color := "%" + strconv.Itoa(s.Width+4*lenColor) + "v\n"
		fmtStr2Color := "%" + strconv.Itoa(s.Width+2*lenColor) + "v\n"
		fmt.Printf(fmtStr1Color, elapsedTime)
		fmt.Printf(fmtStr4Color, tokenInfo)
		fmt.Printf(fmtStr2Color, cumulativeTokenInfo)

	}
}
