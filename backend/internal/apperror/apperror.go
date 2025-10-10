package apperror

import (
	"errors"
	"net/http"
)

type AppError struct {
	Code    string
	Message string
}

// Error implements the error interface for AppError
func (e AppError) Error() string {
	return e.Message
}

// Code extracts the error Code from an error, returning the AppError Code if it's an AppError, otherwise returns INTERNAL_ERROR
func Code(err error) string {
	if isAppError(err) {
		return err.(AppError).Code
	}
	return ErrInternal.Code
}

// Message extracts the error Message from an error, returning the AppError Message if it's an AppError, otherwise returns a generic internal error Message
func Message(err error) string {
	if isAppError(err) {
		return err.(AppError).Message
	}
	return ErrInternal.Message
}

// isAppError checks if an error is an AppError (safe to expose to frontend)
func isAppError(err error) bool {
	var appError AppError
	ok := errors.As(err, &appError)
	return ok
}

// isErrorType checks if err matches any of the provided target errors
func isErrorType(err error, targets ...error) bool {
	for _, target := range targets {
		if errors.Is(err, target) {
			return true
		}
	}
	return false
}

// StatusFromError maps custom errors to HTTP status codes
func StatusFromError(err error) int {
	switch {
	// 400 Bad Request
	case isErrorType(err, ErrBadRequest, ErrInvalidID, ErrInvalidMembershipData):
		return http.StatusBadRequest
	// 401 Unauthorized
	case isErrorType(err, ErrInvalidCredentials, ErrInvalidToken, ErrInvalidClaims, ErrInvalidIssuer, ErrInvalidAudience, ErrTokenInvalidated):
		return http.StatusUnauthorized
	// 403 Forbidden
	case isErrorType(err, ErrForbidden, ErrUserInactive, ErrUserNotMember):
		return http.StatusForbidden
	// 404 Not Found
	case isErrorType(err, ErrUserNotFound, ErrCommunityNotFound, ErrMembershipNotFound):
		return http.StatusNotFound
	// 409 Conflict
	case isErrorType(err, ErrUsernameExists, ErrEmailExists, ErrCommunityNameExists, ErrAlreadyMember):
		return http.StatusConflict
	// 500 Internal Server Error
	case isErrorType(err, ErrInternal, ErrNoFieldsToUpdate, ErrMembershipCreateFailed, ErrMembershipDeleteFailed):
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}

var (
	// Auth-related
	ErrInvalidCredentials = AppError{Code: "INVALID_CREDENTIALS", Message: "Invalid username or password"}
	ErrInvalidToken       = AppError{Code: "INVALID_TOKEN", Message: "Invalid or expired token"}
	ErrInvalidClaims      = AppError{Code: "INVALID_CLAIMS", Message: "Invalid token claims"}
	ErrInvalidIssuer      = AppError{Code: "INVALID_ISSUER", Message: "Invalid token issuer"}
	ErrInvalidAudience    = AppError{Code: "INVALID_AUDIENCE", Message: "Invalid token audience"}
	ErrTokenInvalidated   = AppError{Code: "TOKEN_INVALIDATED", Message: "Token has been invalidated"}
	ErrForbidden          = AppError{Code: "FORBIDDEN", Message: "You do not have permission to perform this action"}
	ErrBadRequest         = AppError{Code: "BAD_REQUEST", Message: "Bad request"}

	// Generic
	ErrInternal         = AppError{Code: "INTERNAL_ERROR", Message: "Internal server error"}
	ErrNoFieldsToUpdate = AppError{Code: "NO_FIELDS_TO_UPDATE", Message: "No fields provided to update"}
	ErrInvalidID        = AppError{Code: "INVALID_ID", Message: "Invalid ID format"}

	// User-related
	ErrUserNotFound   = AppError{Code: "USER_NOT_FOUND", Message: "User not found"}
	ErrUsernameExists = AppError{Code: "USERNAME_EXISTS", Message: "Username already exists"}
	ErrEmailExists    = AppError{Code: "EMAIL_EXISTS", Message: "Email already exists"}
	ErrUserInactive   = AppError{Code: "USER_INACTIVE", Message: "User account is inactive"}

	// Community-related
	ErrCommunityNotFound   = AppError{Code: "COMMUNITY_NOT_FOUND", Message: "Community not found"}
	ErrCommunityNameExists = AppError{Code: "COMMUNITY_NAME_EXISTS", Message: "Community name already exists"}
	ErrUserNotMember       = AppError{Code: "USER_NOT_MEMBER", Message: "User is not a member of this community"}

	// Membership-related
	ErrMembershipNotFound     = AppError{Code: "MEMBERSHIP_NOT_FOUND", Message: "Membership not found"}
	ErrAlreadyMember          = AppError{Code: "ALREADY_MEMBER", Message: "User is already a member of this community"}
	ErrMembershipCreateFailed = AppError{Code: "MEMBERSHIP_CREATE_FAILED", Message: "Failed to create membership"}
	ErrMembershipDeleteFailed = AppError{Code: "MEMBERSHIP_DELETE_FAILED", Message: "Failed to delete membership"}
	ErrInvalidMembershipData  = AppError{Code: "INVALID_MEMBERSHIP_DATA", Message: "Invalid membership data"}
)
