package content

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/reek004/MarloweQuill/internal/api"
)

type ContentGenerator struct {
	geminiClient *api.GeminiClient
}

type ContentType string

const (
	BlogPost     ContentType = "blog_post"
	Article      ContentType = "article"
	SocialMedia  ContentType = "social_media"
	Script       ContentType = "script"
)

type GenerationOptions struct {
	ContentType    ContentType
	Topic          string
	Tone           string
	Length         int
	AdditionalContext string
}

func NewContentGenerator(client *api.GeminiClient) *ContentGenerator {
	return &ContentGenerator{
		geminiClient: client,
	}
}

func (g *ContentGenerator) GenerateContent(opts GenerationOptions) (string, error) {
	// Construct prompt based on content type and options
	prompt := g.constructPrompt(opts)

	// Generate content using Gemini API
	content, err := g.geminiClient.GenerateContent(prompt)
	if err != nil {
		return "", err
	}

	return content, nil
}

func (g *ContentGenerator) constructPrompt(opts GenerationOptions) string {
	basePrompt := fmt.Sprintf("Write a %s about '%s'", opts.ContentType, opts.Topic)
	
	if opts.Tone != "" {
		basePrompt += fmt.Sprintf(" in a %s tone", opts.Tone)
	}
	
	if opts.Length > 0 {
		basePrompt += fmt.Sprintf(". Aim for approximately %d words", opts.Length)
	}
	
	if opts.AdditionalContext != "" {
		basePrompt += fmt.Sprintf(". Additional context: %s", opts.AdditionalContext)
	}
	
	return basePrompt
}

func (g *ContentGenerator) SaveContent(content string, contentType ContentType) (string, error) {
	// Create output directory if it doesn't exist
	outputDir := "generated_content"
	err := os.MkdirAll(outputDir, os.ModePerm)
	if err != nil {
		return "", err
	}

	// Generate filename with timestamp
	timestamp := time.Now().Format("20060102_150405")
	filename := filepath.Join(outputDir, fmt.Sprintf("%s_%s.txt", contentType, timestamp))

	// Write content to file
	err = os.WriteFile(filename, []byte(content), 0644)
	if err != nil {
		return "", err
	}

	return filename, nil
}