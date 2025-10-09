package apperror

type AppError struct {
	Code    string
	Message string
}

var (
	// Auth-related
	ErrInvalidCredentials = &AppError{Code: "INVALID_CREDENTIALS", Message: "invalid credentials"}
	ErrInvalidToken       = &AppError{Code: "INVALID_TOKEN", Message: "invalid token"}
	ErrInvalidClaims      = &AppError{Code: "INVALID_CLAIMS", Message: "invalid claims"}
	ErrInvalidIssuer      = &AppError{Code: "INVALID_ISSUER", Message: "invalid issuer"}
	ErrInvalidAudience    = &AppError{Code: "INVALID_AUDIENCE", Message: "invalid audience"}
	ErrTokenInvalidated   = &AppError{Code: "TOKEN_INVALIDATED", Message: "token has been invalidated"}
	ErrForbidden          = &AppError{Code: "FORBIDDEN", Message: "you do not have permission to perform this action"}

	// Generic
	ErrInternal         = &AppError{Code: "INTERNAL_ERROR", Message: "internal server error"}
	ErrNoFieldsToUpdate = &AppError{Code: "NO_FIELDS_TO_UPDATE", Message: "no fields provided to update"}
	ErrInvalidID        = &AppError{Code: "INVALID_ID", Message: "invalid id"}

	// User-related
	ErrUserNotFound   = &AppError{Code: "USER_NOT_FOUND", Message: "user not found"}
	ErrUsernameExists = &AppError{Code: "USERNAME_EXISTS", Message: "username already exists"}
	ErrEmailExists    = &AppError{Code: "EMAIL_EXISTS", Message: "email already exists"}

	// Community-related
	ErrCommunityNotFound   = &AppError{Code: "COMMUNITY_NOT_FOUND", Message: "community not found"}
	ErrCommunityNameExists = &AppError{Code: "COMMUNITY_NAME_EXISTS", Message: "community name already exists"}
	ErrUserNotMember       = &AppError{Code: "USER_NOT_MEMBER", Message: "user is not a member of this community"}

	// Membership-related
	ErrMembershipNotFound     = &AppError{Code: "MEMBERSHIP_NOT_FOUND", Message: "membership not found"}
	ErrAlreadyMember          = &AppError{Code: "ALREADY_MEMBER", Message: "user is already a member of this community"}
	ErrMembershipCreateFailed = &AppError{Code: "MEMBERSHIP_CREATE_FAILED", Message: "failed to create membership"}
	ErrMembershipDeleteFailed = &AppError{Code: "MEMBERSHIP_DELETE_FAILED", Message: "failed to delete membership"}
	ErrInvalidMembershipData  = &AppError{Code: "INVALID_MEMBERSHIP_DATA", Message: "invalid membership data provided"}
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
