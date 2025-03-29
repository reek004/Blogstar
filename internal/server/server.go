package server

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/reek004/MarloweQuill/internal/api"
	"github.com/reek004/MarloweQuill/internal/content"
	"github.com/reek004/MarloweQuill/internal/middleware"
)

type Server struct {
	Generator   *content.ContentGenerator
	RateLimiter *middleware.RateLimiter
}

type GenerateRequest struct {
	ContentType       string `json:"content_type"`
	Topic             string `json:"topic"`
	Tone              string `json:"tone"`
	Length            int    `json:"length"`
	AdditionalContext string `json:"additional_context"`
}

type GenerateResponse struct {
	Content  string `json:"content"`
	Filename string `json:"filename,omitempty"`
	Error    string `json:"error,omitempty"`
}

func NewServer(geminiClient *api.GeminiClient) *Server {
	// Create a rate limiter: 5 requests per minute per IP
	rateLimiter := middleware.NewRateLimiter(1*time.Minute, 5)

	return &Server{
		Generator:   content.NewContentGenerator(geminiClient),
		RateLimiter: rateLimiter,
	}
}

func (s *Server) Start(port string) error {
	mux := http.NewServeMux()

	// Register handlers with rate limiting and CORS
	mux.HandleFunc("/api/generate", s.RateLimiter.MiddlewareFunc(s.enableCORS(s.handleGenerate)))

	// Add a simple health check endpoint (no rate limiting here)
	mux.HandleFunc("/health", s.enableCORS(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Server is running"))
	}))

	log.Printf("Starting server on :%s", port)
	return http.ListenAndServe(":"+port, mux)
}

// Add CORS support
func (s *Server) enableCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// Handle preflight requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next(w, r)
	}
}

func (s *Server) handleGenerate(w http.ResponseWriter, r *http.Request) {
	// Only accept POST requests
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse request body
	var req GenerateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate request
	if req.ContentType == "" || req.Topic == "" {
		http.Error(w, "Content type and topic are required", http.StatusBadRequest)
		return
	}

	// Convert request to content generation options
	opts := content.GenerationOptions{
		ContentType:       content.ContentType(req.ContentType),
		Topic:             req.Topic,
		Tone:              req.Tone,
		Length:            req.Length,
		AdditionalContext: req.AdditionalContext,
	}

	// Generate content
	generatedContent, err := s.Generator.GenerateContent(opts)
	if err != nil {
		response := GenerateResponse{
			Error: err.Error(),
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Save the generated content and get the filename
	filename, err := s.Generator.SaveContent(generatedContent, opts.ContentType)
	if err != nil {
		log.Printf("Error saving content: %v", err)
		// Continue even if saving fails
	}

	// Return the generated content
	response := GenerateResponse{
		Content:  generatedContent,
		Filename: filename,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
