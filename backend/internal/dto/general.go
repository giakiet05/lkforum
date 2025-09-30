package dto

// ErrorResponse represents a standard error response.
type ErrorResponse struct {
	Code  string `json:"code"`
	Error string `json:"error"`
}

// SuccessResponse represents a successful operation (use for POST/PUT/DELETE, except special cases).
type SuccessResponse struct {
	ID      string `json:"id"`
	Message string `json:"message"`
}
