package main

import (
	"encoding/json"
	"log"

	"github.com/sashabaranov/go-openai"
)

const (
	ModelTokenLimit         = 4096
	ModelMaxCompletionToken = 1024
)

func contexLengthAdjust(messages []openai.ChatCompletionMessage) []openai.ChatCompletionMessage {

	tokenLen := countToken(messages)
	log.Println("tokenLen:", tokenLen)

	for tokenLen+ModelMaxCompletionToken > ModelTokenLimit {
		log.Printf("expected token length %d is exceded the token limit %d",
			tokenLen+ModelMaxCompletionToken, ModelTokenLimit)
		log.Println(colorStr(Red, "remove oldest message:"), messages[1])

		messages = append(messages[0:1], messages[2:]...) // remove oldest message, except system message
		tokenLen = countToken(messages)                   // count again
		log.Println("reduced tokenLen:", tokenLen)
	}
	return messages
}

func countToken(messages []openai.ChatCompletionMessage) int {
	b, err := json.Marshal(messages)
	if err != nil {
		log.Fatal(err)
	}
	ids, _, err := tictoken.Encode(string(b))
	if err != nil {
		log.Fatal(err)
	}
	log.Println("context token length:", len(ids))
	return len(ids)
}
