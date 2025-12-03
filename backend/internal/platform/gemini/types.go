package gemini

// ContentCheckRequest represents a request to check content for violations
type ContentCheckRequest struct {
	Title     string   `json:"title"`
	Text      string   `json:"text"`
	ImageURLs []string `json:"image_urls,omitempty"`
	VideoURLs []string `json:"video_urls,omitempty"`
}

// ContentCheckResponse represents the result of content moderation
type ContentCheckResponse struct {
	IsViolation bool     `json:"is_violation"`
	Confidence  float64  `json:"confidence"`
	Categories  []string `json:"categories"`
	Reason      string   `json:"reason"`
}
