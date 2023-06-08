package main

import (
	"bufio"
	"os"

	"github.com/sashabaranov/go-openai"
)

const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	Purple = "\033[35m"
	Cyan   = "\033[36m"
	Gray   = "\033[37m"
	White  = "\033[97m"
)

var (
	OPENAI_API_KEY string
	client         *openai.Client
	messages       []openai.ChatCompletionMessage
	reader         *bufio.Reader
)

func init() {
	OPENAI_API_KEY = os.Getenv("OPENAI_API_KEY")
	client = openai.NewClient(OPENAI_API_KEY)
	messages = make([]openai.ChatCompletionMessage, 0, 10)
	reader = bufio.NewReader(os.Stdin)
}
