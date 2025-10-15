package dto

// ErrorResponse represents a standard error response.
type ErrorResponse struct {
	ErrorCode string `json:"error_code"`
	Message   string `json:"message"`
}

// SuccessResponse represents a successful operation (use for POST/PUT/DELETE, except special cases).
type SuccessResponse struct {
	ID      string `json:"id"`
	Message string `json:"message"`
}
