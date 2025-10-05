package apperror

type AppError struct {
	Code    string
	Message string
}

var (
	ErrInvalidCredentials = &AppError{Code: "INVALID_CREDENTIALS", Message: "invalid credentials"}
	ErrUsernameExists     = &AppError{Code: "USERNAME_EXISTS", Message: "username already exists"}
	ErrEmailExists        = &AppError{Code: "EMAIL_EXISTS", Message: "email already exists"}
	ErrUserNotFound       = &AppError{Code: "USER_NOT_FOUND", Message: "user not found"}
	ErrInternal           = &AppError{Code: "INTERNAL_ERROR", Message: "internal server error"}
	ErrInvalidToken       = &AppError{Code: "INVALID_TOKEN", Message: "invalid token"}
	ErrInvalidClaims      = &AppError{Code: "INVALID_CLAIMS", Message: "invalid claims"}
	ErrInvalidIssuer      = &AppError{Code: "INVALID_ISSUER", Message: "invalid issuer"}
	ErrInvalidAudience    = &AppError{Code: "INVALID_AUDIENCE", Message: "invalid audience"}
	ErrTokenInvalidated   = &AppError{Code: "TOKEN_INVALIDATED", Message: "token has been invalidated"}
)

func (e *AppError) Error() string {
	return e.Message
}

func Code(err error) string {
	if appErr, ok := err.(*AppError); ok {
		return appErr.Code
	}
	return ErrInternal.Code
}
func NewError(originalErr error, code, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
	}
}
