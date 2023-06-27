package main

import (
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

func main() {

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
		// length check and adjust. except the system message
		// ids, tokens, err
		messages = contexLengthAdjust(messages)

		// request completion
		sTime := time.Now()
		loading.Start() // Start the spinner
		resp, err := chatComplete(messages)
		if err != nil {
			fmt.Println(colorStr(Red, fmt.Sprintf("ChatCompletion error: %v", err)))
			fmt.Println()
			continue
		}
		loading.Stop()
		eTime := time.Since(sTime)

		content := resp.Choices[0].Message.Content
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleAssistant,
			Content: content,
		})
		fmt.Printf("\n%s %s\n", colorStr(Blue, "Assistant:"), content)

		// print elapsed time, tokenInfo
		totalPromptTokens += resp.Usage.PromptTokens
		totalCompletionTokens += resp.Usage.CompletionTokens
		elapsedTime := prepareElapsedTime(eTime)
		tokenInfo := prepareTokenInfo(resp.Usage)
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
