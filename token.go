package main

import (
	"errors"
	"fmt"
	"log"

	"github.com/pkoukk/tiktoken-go"
	"github.com/sashabaranov/go-openai"
)

const (
	ModelTokenLimit         = 4096
	ModelMaxCompletionToken = 1024
)

func contexLengthAdjust(messages []openai.ChatCompletionMessage) []openai.ChatCompletionMessage {

	tokenLen := NumTokensFromMessages(messages, openai.GPT3Dot5Turbo+"-0613")
	log.Println("tokenLen:", tokenLen)

	for tokenLen+ModelMaxCompletionToken > ModelTokenLimit {
		log.Printf("expected token length %d is exceded the token limit %d",
			tokenLen+ModelMaxCompletionToken, ModelTokenLimit)
		log.Println(colorStr(Red, "remove oldest message:"), messages[1])

		messages = append(messages[0:1], messages[2:]...)                         // remove oldest message, except system message
		tokenLen := NumTokensFromMessages(messages, openai.GPT3Dot5Turbo+"-0613") // count again
		log.Println("reduced tokenLen:", tokenLen)
	}
	return messages
}

// below link may not work on Chrome(error: Unable to render code block)
// then, use FireFox
// OpenAI Cookbook: https://github.com/openai/openai-cookbook/blob/main/examples/How_to_count_tokens_with_tiktoken.ipynb
func NumTokensFromMessages(messages []openai.ChatCompletionMessage, model string) (numTokens int) {
	tkm, err := tiktoken.EncodingForModel(model)
	if err != nil {
		err = fmt.Errorf("encoding for model: %v", err)
		log.Println(err)
		return
	}

	var tokensPerMessage, tokensPerName int

	if model == "gpt-3.5-turbo-0613" ||
		model == "gpt-3.5-turbo-16k-0613" ||
		model == "gpt-4-0314" ||
		model == "gpt-4-32k-0314" ||
		model == "gpt-4-0613" ||
		model == "gpt-4-32k-0613" {
		tokensPerMessage = 3
		tokensPerName = -1
	} else if model == "gpt-3.5-turbo-0301" {
		tokensPerMessage = 4 // every message follows <|start|>{role/name}\n{content}<|end|>\n
		tokensPerName = -1   // if there's a name, the role is omitted
	} else if model == "gpt-3.5-turbo" {
		log.Println("warning: gpt-3.5-turbo may update over time. Returning num tokens assuming gpt-3.5-turbo-0613.")
		return NumTokensFromMessages(messages, "gpt-3.5-turbo-0613")
	} else if model == "gpt-4" {
		log.Println("warning: gpt-4 may update over time. Returning num tokens assuming gpt-4-0613.")
		return NumTokensFromMessages(messages, "gpt-4-0613")
	} else {
		err := errors.New("warning: model not found. Using cl100k_base encoding")
		log.Println(err)
		return
	}

	for _, message := range messages {
		numTokens += tokensPerMessage
		numTokens += len(tkm.Encode(message.Content, nil, nil))
		numTokens += len(tkm.Encode(message.Role, nil, nil))
		numTokens += len(tkm.Encode(message.Name, nil, nil))
		if message.Name != "" {
			numTokens += tokensPerName
		}
	}
	numTokens += 3 // every reply is primed with <|start|>assistant<|message|>
	return numTokens
}
