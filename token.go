package main

import (
	"fmt"
	"log"

	"github.com/sashabaranov/go-openai"
)

const (
	ModelTokenLimit         = 4096
	ModelMaxCompletionToken = 1024
)

func contexLengthAdjust(messages []openai.ChatCompletionMessage) []openai.ChatCompletionMessage {

	for tokenLen := countToken(messages); tokenLen+ModelMaxCompletionToken > ModelTokenLimit; {
		log.Printf("expected token length %d is exceded the token limit %d",
			tokenLen+ModelMaxCompletionToken, ModelTokenLimit)
		log.Println("message removed:", messages[1])

		messages = append(messages[0:1], messages[2:]...) // remove oldest message, except system message
		tokenLen = countToken(messages)                   // count again
	}
	return messages
}

func countToken(messages []openai.ChatCompletionMessage) int {
	ids, _, err := tictoken.Encode(fmt.Sprintf("%v", messages))
	if err != nil {
		log.Fatal(err)
	}
	log.Println("context token length:", len(ids))
	return len(ids)
}
