// internal/api/gemini.go
package api

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/reek004/MarloweQuill/internal/config"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

type GeminiClient struct {
	client *genai.Client
	config *config.Config
}

func NewGeminiClient(cfg *config.Config) *GeminiClient {
	client, err := genai.NewClient(context.Background(), option.WithAPIKey(cfg.GeminiAPIKey))
	if err != nil {
		log.Fatal("Error creating Gemini client:", err)
	}

	return &GeminiClient{
		client: client,
		config: cfg,
	}
}

func (g *GeminiClient) GenerateContent(prompt string) (string, error) {
	// Try each model from configuration
	for _, modelName := range g.config.DefaultModels {
		model := g.client.GenerativeModel(modelName)
		
		// Create parts from the prompt
		parts := []genai.Part{
			genai.Text(prompt),
		}

		resp, err := model.GenerateContent(context.Background(), parts...)
		if err != nil {
			log.Printf("Error with model %s: %v", modelName, err)
			continue
		}

		// Extract text from the response
		if len(resp.Candidates) > 0 && len(resp.Candidates[0].Content.Parts) > 0 {
			var contentParts []string
			for _, part := range resp.Candidates[0].Content.Parts {
				if textPart, ok := part.(genai.Text); ok {
					contentParts = append(contentParts, string(textPart))
				}
			}

			if len(contentParts) > 0 {
				return strings.Join(contentParts, "\n"), nil
			}
		}
	}

	return "", fmt.Errorf("failed to generate content with any available model")
}

func (g *GeminiClient) Close() {
	g.client.Close()
}