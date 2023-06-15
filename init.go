package main

import (
	"bufio"
	"io"
	"log"
	"os"

	"github.com/sashabaranov/go-openai"
	"github.com/tiktoken-go/tokenizer"
)

const (
	Reset         = "\033[0m"
	Red           = "\033[31m"
	Green         = "\033[32m"
	Yellow        = "\033[33m"
	Blue          = "\033[34m"
	Purple        = "\033[35m"
	Cyan          = "\033[36m"
	Gray          = "\033[37m"
	White         = "\033[97m"
	lenColorCode  = len(Red)
	lenColorReset = len(Reset)
	lenColor      = lenColorCode + lenColorReset
)

var (
	OPENAI_API_KEY string
	client         *openai.Client
	messages       []openai.ChatCompletionMessage
	systemMessage  = openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleSystem,
		Content: "You are a helpful assitent and your name is Jin",
	}
	reader   *bufio.Reader
	tictoken tokenizer.Codec
)

func init() {
	runMode := os.Getenv("RUN_MODE")
	if runMode != "dev" {
		log.SetOutput(io.Discard)
	}
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	reader = bufio.NewReader(os.Stdin)

	OPENAI_API_KEY = APIKey()
	client = openai.NewClient(OPENAI_API_KEY)

	messages = make([]openai.ChatCompletionMessage, 0, 10)
	messages = append(messages, systemMessage)

	var err error
	tictoken, err = tokenizer.Get(tokenizer.Cl100kBase)
	if err != nil {
		log.Fatal(err)
	}
}
