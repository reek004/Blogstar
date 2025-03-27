package main

import (
	"fmt"
	"log"

	"github.com/reek004/MarloweQuill/internal/api"
	"github.com/reek004/MarloweQuill/internal/config"
	"github.com/reek004/MarloweQuill/internal/content"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Create Gemini client
	geminiClient := api.NewGeminiClient(cfg)
	defer geminiClient.Close()

	// Create content generator
	generator := content.NewContentGenerator(geminiClient)

	// Generate different types of content
	contentTypes := []content.GenerationOptions{
		{
			ContentType: content.BlogPost,
			Topic:       "Artificial Intelligence in 2024",
			Tone:        "professional",
			Length:      500,
		},
		{
			ContentType: content.SocialMedia,
			Topic:       "Climate Change Awareness",
			Tone:        "motivational",
			Length:      100,
		},
	}

	// Generate and save multiple content types
	for _, opts := range contentTypes {
		fmt.Printf("Generating %s content...\n", opts.ContentType)
		
		generatedContent, err := generator.GenerateContent(opts)
		if err != nil {
			log.Printf("Error generating %s content: %v", opts.ContentType, err)
			continue
		}

		// Save the generated content
		filename, err := generator.SaveContent(generatedContent, opts.ContentType)
		if err != nil {
			log.Printf("Error saving %s content: %v", opts.ContentType, err)
			continue
		}

		fmt.Printf("Generated %s content saved to %s\n", opts.ContentType, filename)
		fmt.Println("Content Preview:")
		fmt.Println(generatedContent[:min(500, len(generatedContent))] + "...")
	}
}

// Helper function to get minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
