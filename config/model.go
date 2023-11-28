package config

import (
	"fmt"

	"github.com/sashabaranov/go-openai"
)

// go-openai models 2023-11-28
// GPT4TurboPreview      = "gpt-4-1106-preview"
// GPT4VisionPreview     = "gpt-4-vision-preview"
// GPT4                  = "gpt-4"
// GPT432K               = "gpt-4-32k"
// GPT40613              = "gpt-4-0613"
// GPT432K0613           = "gpt-4-32k-0613"

const (
	DefaultModel = openai.GPT4
)

var gpt4Models = ModelsInfo{
	"gpt-4-1106-preview": {
		GeneralDescription: "GPT-4 TurboNew - The latest GPT-4 model with improved instruction following, JSON mode, reproducible outputs, parallel function calling, and more. Returns a maximum of 4,096 output tokens. This preview model is not yet suited for production traffic.",
		ContextWindow:      "128,000 tokens",
		TrainingData:       "Up to Apr 2023",
	},
	"gpt-4-vision-preview": {
		GeneralDescription: "GPT-4 Turbo with visionNew - Ability to understand images, in addition to all other GPT-4 Turbo capabilities. Returns a maximum of 4,096 output tokens. This is a preview model version and not suited yet for production traffic.",
		ContextWindow:      "128,000 tokens",
		TrainingData:       "Up to Apr 2023",
	},
	"gpt-4": {
		GeneralDescription: "Currently points to gpt-4-0613. See continuous model upgrades.",
		ContextWindow:      "8,192 tokens",
		TrainingData:       "Up to Sep 2021",
	},
	"gpt-4-32k": {
		GeneralDescription: "Currently points to gpt-4-32k-0613. See continuous model upgrades.",
		ContextWindow:      "32,768 tokens",
		TrainingData:       "Up to Sep 2021",
	},
	"gpt-4-0613": {
		GeneralDescription: "Snapshot of gpt-4 from June 13th 2023 with improved function calling support.",
		ContextWindow:      "8,192 tokens",
		TrainingData:       "Up to Sep 2021",
	},
	"gpt-4-32k-0613": {
		GeneralDescription: "Snapshot of gpt-4-32k from June 13th 2023 with improved function calling support.",
		ContextWindow:      "32,768 tokens",
		TrainingData:       "Up to Sep 2021",
	},
}

// ModelDescription stores detailed information about each model
type ModelDescription struct {
	GeneralDescription string
	ContextWindow      string
	TrainingData       string
}

// ModelsInfo holds all models' descriptions
type ModelsInfo map[string]ModelDescription

// DescribeModel prints the description of a given model
func DescribeModel(modelName string) {
	desc, exists := gpt4Models[modelName]
	if !exists {
		fmt.Println("Model not found.")
		return
	}
	fmt.Printf("Model: %s\n", modelName)
	fmt.Printf(". General Description: %s\n", desc.GeneralDescription)
	fmt.Printf(". Context Window: %s\n", desc.ContextWindow)
	fmt.Printf(". Training Data: %s\n\n", desc.TrainingData)
}

// ListModels returns a list of all model names
func ListModels() []string {
	var models []string
	for model := range gpt4Models {
		models = append(models, model)
	}
	return models
}
