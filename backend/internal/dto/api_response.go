package dto

import "github.com/gin-gonic/gin"

type ApiResponse struct {
	Success   bool        `json:"success"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data,omitempty"` // omitempty: nếu data là nil thì không hiển thị
	ErrorCode string      `json:"error_code,omitempty"`
}

func SendSuccess(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, ApiResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func SendError(c *gin.Context, statusCode int, message string, errorCode string) {
	c.JSON(statusCode, ApiResponse{
		Success:   false,
		Message:   message,
		ErrorCode: errorCode,
	})
}
