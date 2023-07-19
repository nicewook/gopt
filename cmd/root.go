/*
Copyright © 2023 Hyunseok Jeong nicewook@hotmail.com

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

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
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gopt",
	Short: "gopt is OpenAI GPT powered chat completion CLI",
	Long: `gopt is OpenAI GPT powered chat completion CLI.
You can ask questions to the GPT model and get answers.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		startChat()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func startChat() {

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
			completionTokensLen := NumTokensFromText(responseContent, viper.GetString("model"))
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
