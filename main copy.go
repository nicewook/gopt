package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/chzyer/readline"
	"github.com/sashabaranov/go-openai"
)

func main2() {

	// usage
	// help, model, token and bill info
	// only one line question - 마지막 글자가 Ctrl+D이면 끝내게 하기?
	fmt.Println()
	fmt.Println(helpMessage())

	var totalPromptTokens, totalCompletionTokens int

	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	appConfigPath := filepath.Join(homeDir, ".gopt")

	var completer = readline.NewPrefixCompleter(
		readline.PcItem("help"),
		readline.PcItem("config"),
		readline.PcItem("context"),
		readline.PcItem("reset"),
		readline.PcItem("clear"),
		readline.PcItem("exit"),
	)
	readlineConfig := &readline.Config{
		Prompt:          colorStr(Cyan, "gpt> "),
		HistoryFile:     filepath.Join(appConfigPath, "readline-history"),
		AutoComplete:    completer,
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",
	}
	rl, err := readline.NewEx(readlineConfig)
	if err != nil {
		log.Fatal(err)
	}
	defer rl.Close()
	rl.CaptureExitSignal()

	loading := spinner.New([]string{".", "..", "...", "....", "....."}, 150*time.Millisecond)
	loading.Prefix = colorStr(Yellow, "loading")
	loading.Color("yellow")

	// start loop
	for {
		// get user input
		userInput, err := rl.Readline()
		if err == readline.ErrInterrupt {
			if len(userInput) == 0 {
				break
			} else {
				continue
			}
		} else if err == io.EOF {
			break
		}
		userInput = strings.TrimSpace(userInput)

		// TODO: reserved command execution
		// help, msg(show), clear, sysmsg(get, set), config
		if commandExecute(rl.Terminal, userInput) {
			continue
		}

		// make messages
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleUser,
			Content: userInput,
		})

		// request completion

		var (
			resp            openai.ChatCompletionResponse
			isStreamMode    bool = true
			responseContent string
			curUsage        openai.Usage
			promptTokenLen  int
		)
		// length check and adjust. except the system message
		sTime := time.Now()
		messages, promptTokenLen = contexLengthAdjust(messages)

		log.Println("stream mode:", isStreamMode)
		fmt.Println()
		if isStreamMode {
			stream, err := chatCompleteStream(messages)
			if err != nil {
				fmt.Println(colorStr(Red, fmt.Sprintf("ChatCompletion error: %v", err)))
				fmt.Println()
				continue
			}
			defer stream.Close()

			fmt.Print(colorStr(Blue, "Assistant: "))
			for {
				response, err := stream.Recv()
				if errors.Is(err, io.EOF) {
					fmt.Println()
					log.Println("Stream finished")
					break
				}

				if err != nil {
					log.Printf("Stream error: %v\n", err)
					break
				}

				responseContent += response.Choices[0].Delta.Content
				fmt.Print(response.Choices[0].Delta.Content)
			}
			completionTokensLen := NumTokensFromText(responseContent, GPT3Dot5Turbo0613)
			curUsage = openai.Usage{
				PromptTokens:     promptTokenLen,
				CompletionTokens: completionTokensLen,
				TotalTokens:      promptTokenLen + completionTokensLen,
			}
		} else {
			loading.Start() // Start the spinner
			resp, err = chatComplete(messages)
			if err != nil {
				fmt.Println(colorStr(Red, fmt.Sprintf("ChatCompletion error: %v", err)))
				fmt.Println()
				continue
			}
			loading.Stop()
			responseContent = resp.Choices[0].Message.Content
			curUsage = resp.Usage
			fmt.Println(colorStr(Blue, "Assistant:"), responseContent)

		}
		eTime := time.Since(sTime)

		respMessage := openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleAssistant,
			Content: responseContent,
		}
		messages = append(messages, respMessage)

		// print elapsed time, tokenInfo
		totalPromptTokens += curUsage.PromptTokens
		totalCompletionTokens += curUsage.CompletionTokens

		if !isStreamMode {
			log.Println("tiktoken calc and real response Ussage compare in stream mode")
			log.Printf("prompt: %d, %d", curUsage.PromptTokens, resp.Usage.PromptTokens)
			log.Printf("completion: %d, %d", curUsage.CompletionTokens, resp.Usage.CompletionTokens)
			log.Printf("total: %d, %d", curUsage.TotalTokens, resp.Usage.TotalTokens)
		}

		elapsedTime := prepareElapsedTime(eTime)
		tokenInfo := prepareTokenInfo(curUsage)
		cumulativeTokenInfo := prepareCumulativeTokenInfo(totalPromptTokens, totalCompletionTokens)

		fmt.Printf(rightAlignWithColorWords(1), elapsedTime)
		fmt.Printf(rightAlignWithColorWords(4), tokenInfo)
		fmt.Printf(rightAlignWithColorWords(2), cumulativeTokenInfo)

	}
}

func rightAlignWithColorWords(wordCount int) string {
	// color prompt should be counted for right alignment
	virtualWidth := wordCount*lenColor + readline.GetScreenWidth()
	return "%" + strconv.Itoa(virtualWidth) + "v\n"
}
