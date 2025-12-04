package gemini

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/giakiet05/lkforum/internal/config"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

type GeminiClient struct {
	client     *genai.Client
	model      *genai.GenerativeModel
	config     *config.GeminiConfig
	httpClient *http.Client
}

func NewGeminiClient(cfg *config.GeminiConfig) (*GeminiClient, error) {
	if !cfg.Enabled {
		log.Println("Gemini moderation is disabled")
		return nil, nil
	}

	if cfg.APIKey == "" {
		return nil, fmt.Errorf("GEMINI_API_KEY is required")
	}

	ctx := context.Background()

	// Create Gemini client
	client, err := genai.NewClient(ctx, option.WithAPIKey(cfg.APIKey))
	if err != nil {
		return nil, fmt.Errorf("failed to create Gemini client: %w", err)
	}

	// Configure model
	model := client.GenerativeModel(cfg.Model)
	model.SetTemperature(0.2) // Low temperature for consistent moderation
	model.ResponseMIMEType = "application/json"

	// Disable safety filters - we want to analyze all content
	model.SafetySettings = []*genai.SafetySetting{
		{Category: genai.HarmCategoryHarassment, Threshold: genai.HarmBlockNone},
		{Category: genai.HarmCategoryHateSpeech, Threshold: genai.HarmBlockNone},
		{Category: genai.HarmCategorySexuallyExplicit, Threshold: genai.HarmBlockNone},
		{Category: genai.HarmCategoryDangerousContent, Threshold: genai.HarmBlockNone},
	}

	// HTTP client for downloading images
	httpClient := &http.Client{
		Timeout: time.Duration(cfg.Timeout) * time.Second,
	}

	log.Printf("Gemini client initialized with model: %s", cfg.Model)

	return &GeminiClient{
		client:     client,
		model:      model,
		config:     cfg,
		httpClient: httpClient,
	}, nil
}

func (c *GeminiClient) CheckContent(ctx context.Context, req *ContentCheckRequest) (*ContentCheckResponse, error) {
	if c == nil {
		// Moderation disabled - approve all content
		return &ContentCheckResponse{
			IsViolation: false,
			Confidence:  0,
			Categories:  []string{},
			Reason:      "Moderation disabled",
		}, nil
	}

	// Build prompt
	prompt := buildModerationPrompt(req)

	// Prepare parts
	parts := []genai.Part{genai.Text(prompt)}

	// Add images
	for _, imageURL := range req.ImageURLs {
		imgData, mimeType, err := c.downloadImage(imageURL)
		if err != nil {
			log.Printf("Failed to download image %s: %v", imageURL, err)
			continue
		}

		// Resize if too large (Gemini limit: 20MB)
		if len(imgData) > 10*1024*1024 {
			log.Printf("Image too large (%d bytes), skipping: %s", len(imgData), imageURL)
			continue
		}

		parts = append(parts, genai.ImageData(mimeType, imgData))
	}

	// Add video thumbnails (if any)
	for _, videoURL := range req.VideoURLs {
		// For videos, we just note them in the prompt
		// Actual video analysis would require more complex processing
		log.Printf("Video URL provided: %s (thumbnail check only)", videoURL)
	}

	// Call Gemini with retry
	var resp *genai.GenerateContentResponse
	var err error

	for attempt := 0; attempt < c.config.MaxRetries; attempt++ {
		resp, err = c.model.GenerateContent(ctx, parts...)
		if err == nil {
			break
		}

		log.Printf("Gemini API error (attempt %d/%d): %v", attempt+1, c.config.MaxRetries, err)
		if attempt < c.config.MaxRetries-1 {
			time.Sleep(time.Duration(attempt+1) * time.Second) // Exponential backoff
		}
	}

	if err != nil {
		return nil, fmt.Errorf("gemini API failed after %d retries: %w", c.config.MaxRetries, err)
	}

	// Parse response
	return c.parseResponse(resp)
}

func (c *GeminiClient) downloadImage(url string) ([]byte, string, error) {
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf("HTTP %d: %s", resp.StatusCode, resp.Status)
	}

	// Read image data
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", err
	}

	// Detect MIME type
	mimeType := resp.Header.Get("Content-Type")
	if mimeType == "" {
		mimeType = detectMIMEType(data)
	}

	return data, mimeType, nil
}

func detectMIMEType(data []byte) string {
	if len(data) < 12 {
		return "image/jpeg"
	}

	// Check magic numbers
	if bytes.HasPrefix(data, []byte{0xFF, 0xD8, 0xFF}) {
		return "image/jpeg"
	}
	if bytes.HasPrefix(data, []byte{0x89, 0x50, 0x4E, 0x47}) {
		return "image/png"
	}
	if bytes.HasPrefix(data, []byte("GIF")) {
		return "image/gif"
	}
	if bytes.HasPrefix(data, []byte("WEBP")) || bytes.Contains(data[:20], []byte("WEBP")) {
		return "image/webp"
	}

	return "image/jpeg" // Default
}

func (c *GeminiClient) parseResponse(resp *genai.GenerateContentResponse) (*ContentCheckResponse, error) {
	if len(resp.Candidates) == 0 {
		return nil, fmt.Errorf("no candidates in response")
	}

	candidate := resp.Candidates[0]
	if len(candidate.Content.Parts) == 0 {
		return nil, fmt.Errorf("no parts in response")
	}

	// Extract text from response
	var jsonText string
	for _, part := range candidate.Content.Parts {
		if text, ok := part.(genai.Text); ok {
			jsonText = string(text)
			break
		}
	}

	if jsonText == "" {
		return nil, fmt.Errorf("no text in response")
	}

	// Clean JSON (remove markdown code blocks if present)
	jsonText = strings.TrimSpace(jsonText)
	jsonText = strings.TrimPrefix(jsonText, "```json")
	jsonText = strings.TrimPrefix(jsonText, "```")
	jsonText = strings.TrimSuffix(jsonText, "```")
	jsonText = strings.TrimSpace(jsonText)

	// Parse JSON
	var result ContentCheckResponse
	if err := json.Unmarshal([]byte(jsonText), &result); err != nil {
		return nil, fmt.Errorf("failed to parse JSON response: %w (text: %s)", err, jsonText)
	}

	return &result, nil
}

func (c *GeminiClient) Close() error {
	if c != nil && c.client != nil {
		return c.client.Close()
	}
	return nil
}
