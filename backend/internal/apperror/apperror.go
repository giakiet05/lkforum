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
func NewError(originalErr error, code, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
	}
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
	case isErrorType(err, ErrBadRequest, ErrInvalidID, ErrInvalidMembershipData, ErrInvalidOTP, ErrOTPExpired,
		ErrInvalidGender, ErrInvalidDateFormat, ErrAgeTooYoung, ErrInvalidBirthDate, ErrInvalidProvince, ErrTooManyInterests, ErrInvalidInterest, ErrInvalidUsername, ErrPasswordTooShort, ErrPasswordTooWeak):
		return http.StatusBadRequest
	// 401 Unauthorized
	case isErrorType(err, ErrInvalidCredentials, ErrInvalidToken, ErrInvalidClaims, ErrInvalidIssuer, ErrInvalidAudience, ErrTokenInvalidated, ErrMissingAuthHeader, ErrInvalidAuthHeader, ErrMissingToken, ErrNotAuthenticated):
		return http.StatusUnauthorized
	// 403 Forbidden
	case isErrorType(err, ErrForbidden, ErrUserInactive, ErrUserNotMember, ErrEmailNotVerified, ErrAdminAccessRequired):
		return http.StatusForbidden
	// 404 Not Found
	case isErrorType(err, ErrUserNotFound, ErrCommunityNotFound, ErrCommunityDeleted, ErrMembershipNotFound, ErrPostNotFound, ErrVoteNotFound, ErrDraftNotFound, ErrEmailNotRegistered):
		return http.StatusNotFound
	// 409 Conflict
	case isErrorType(err, ErrUsernameExists, ErrEmailExists, ErrCommunityNameExists, ErrAlreadyMember, ErrEmailAlreadyVerified, ErrLoginMethodMismatch, ErrPollVoted, ErrPollCannotEdit, ErrAlreadyReported):
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
	ErrInvalidCredentials   = AppError{Code: "INVALID_CREDENTIALS", Message: "Email hoặc mật khẩu không đúng"}
	ErrInvalidToken         = AppError{Code: "INVALID_TOKEN", Message: "Token không hợp lệ hoặc đã hết hạn"}
	ErrInvalidClaims        = AppError{Code: "INVALID_CLAIMS", Message: "Thông tin token không hợp lệ"}
	ErrInvalidIssuer        = AppError{Code: "INVALID_ISSUER", Message: "Nguồn phát hành token không hợp lệ"}
	ErrInvalidAudience      = AppError{Code: "INVALID_AUDIENCE", Message: "Đối tượng token không hợp lệ"}
	ErrTokenInvalidated     = AppError{Code: "TOKEN_INVALIDATED", Message: "Token đã bị vô hiệu hóa"}
	ErrMissingAuthHeader    = AppError{Code: "MISSING_AUTH_HEADER", Message: "Thiếu Authorization header"}
	ErrInvalidAuthHeader    = AppError{Code: "INVALID_AUTH_HEADER", Message: "Định dạng Authorization header không hợp lệ"}
	ErrMissingToken         = AppError{Code: "MISSING_TOKEN", Message: "Thiếu token xác thực"}
	ErrNotAuthenticated     = AppError{Code: "NOT_AUTHENTICATED", Message: "Chưa xác thực"}
	ErrInvalidAuthContext   = AppError{Code: "INVALID_AUTH_CONTEXT", Message: "Ngữ cảnh xác thực không hợp lệ"}
	ErrAdminAccessRequired  = AppError{Code: "ADMIN_ACCESS_REQUIRED", Message: "Cần quyền quản trị viên"}
	ErrForbidden            = AppError{Code: "FORBIDDEN", Message: "Bạn không có quyền thực hiện hành động này"}
	ErrBadRequest           = AppError{Code: "BAD_REQUEST", Message: "Yêu cầu không hợp lệ"}
	ErrEmailNotVerified     = AppError{Code: "EMAIL_NOT_VERIFIED", Message: "Email chưa được xác thực"}
	ErrEmailAlreadyVerified = AppError{Code: "EMAIL_ALREADY_VERIFIED", Message: "Email đã được xác thực"}
	ErrInvalidOTP           = AppError{Code: "INVALID_OTP", Message: "Mã xác thực không đúng"}
	ErrOTPExpired           = AppError{Code: "OTP_EXPIRED", Message: "Mã xác thực đã hết hạn"}
	ErrLoginMethodMismatch  = AppError{Code: "LOGIN_METHOD_MISMATCH", Message: "Tài khoản này đã đăng ký bằng Google. Vui lòng đăng nhập bằng Google hoặc liên hệ quản trị viên để hỗ trợ."}
	ErrEmailNotRegistered   = AppError{Code: "EMAIL_NOT_REGISTERED", Message: "Email không tồn tại trong hệ thống"}
	ErrInvalidUsername      = AppError{Code: "INVALID_USERNAME", Message: "Username phải từ 3-20 ký tự"}
	ErrPasswordTooShort     = AppError{Code: "PASSWORD_TOO_SHORT", Message: "Mật khẩu phải có ít nhất 8 ký tự"}
	ErrPasswordTooWeak      = AppError{Code: "PASSWORD_TOO_WEAK", Message: "Mật khẩu phải chứa ít nhất 1 chữ hoa, 1 chữ thường, 1 số và 1 ký tự đặc biệt (@$!%*?&)"}

	// Generic
	ErrInternal          = AppError{Code: "INTERNAL_ERROR", Message: "Lỗi hệ thống"}
	ErrNoFieldsToUpdate  = AppError{Code: "NO_FIELDS_TO_UPDATE", Message: "Không có trường nào để cập nhật"}
	ErrInvalidID         = AppError{Code: "INVALID_ID", Message: "Định dạng ID không hợp lệ"}
	ErrPaginationInvalid = AppError{Code: "PAGINATION_INVALID", Message: "Số trang hoặc kích thước trang không hợp lệ. Kích thước trang phải nhỏ hơn 500."}

	// User-related
	ErrUserNotFound   = AppError{Code: "USER_NOT_FOUND", Message: "Không tìm thấy người dùng"}
	ErrUsernameExists = AppError{Code: "USERNAME_EXISTS", Message: "Tên người dùng đã tồn tại"}
	ErrEmailExists    = AppError{Code: "EMAIL_EXISTS", Message: "Email đã được sử dụng"}
	ErrUserInactive   = AppError{Code: "USER_INACTIVE", Message: "Tài khoản người dùng đã bị vô hiệu hóa"}

	// Profile validation
	ErrInvalidGender     = AppError{Code: "INVALID_GENDER", Message: "Giá trị giới tính không hợp lệ"}
	ErrInvalidDateFormat = AppError{Code: "INVALID_DATE_FORMAT", Message: "Định dạng ngày không hợp lệ, sử dụng YYYY-MM-DD"}
	ErrAgeTooYoung       = AppError{Code: "AGE_TOO_YOUNG", Message: "Phải từ 13 tuổi trở lên"}
	ErrInvalidBirthDate  = AppError{Code: "INVALID_BIRTH_DATE", Message: "Ngày sinh không hợp lệ"}
	ErrInvalidProvince   = AppError{Code: "INVALID_PROVINCE", Message: "Tỉnh/thành phố không hợp lệ"}
	ErrTooManyInterests  = AppError{Code: "TOO_MANY_INTERESTS", Message: "Tối đa 10 sở thích"}
	ErrInvalidInterest   = AppError{Code: "INVALID_INTEREST", Message: "Sở thích không hợp lệ"}

	// Community-related
	ErrCommunityNotFound         = AppError{Code: "COMMUNITY_NOT_FOUND", Message: "Không tìm thấy cộng đồng"}
	ErrCommunityNameExists       = AppError{Code: "COMMUNITY_NAME_EXISTS", Message: "Tên cộng đồng đã tồn tại"}
	ErrCommunityDeleted          = AppError{Code: "COMMUNITY_DELETED", Message: "Cộng đồng đã bị xóa"}
	ErrModeratorAlreadyExists    = AppError{Code: "MODERATOR_ALREADY_EXISTS", Message: "Người dùng đã là moderator của cộng đồng."}
	ErrCannotRemoveModerator     = AppError{Code: "CANNOT_REMOVE_MODERATOR", Message: "Không thể xóa moderator này."}
	ErrCannotRemoveCreator       = AppError{Code: "CANNOT_REMOVE_CREATOR", Message: "Không thể xóa người tạo cộng đồng khỏi danh sách moderator."}
	ErrUserIsBannedFromCommunity = AppError{Code: "BANNED_COMMUNITY", Message: "Người dùng đã bị ban khỏi cộng đồng."}

	// Membership-related
	ErrUserNotMember              = AppError{Code: "USER_NOT_MEMBER", Message: "Bạn chưa tham gia cộng đồng này"}
	ErrMembershipNotFound         = AppError{Code: "MEMBERSHIP_NOT_FOUND", Message: "Không tìm thấy thành viên"}
	ErrAlreadyMember              = AppError{Code: "ALREADY_MEMBER", Message: "Bạn đã là thành viên của cộng đồng này"}
	ErrMembershipCreateFailed     = AppError{Code: "MEMBERSHIP_CREATE_FAILED", Message: "Không thể tạo thành viên"}
	ErrMembershipDeleteFailed     = AppError{Code: "MEMBERSHIP_DELETE_FAILED", Message: "Không thể xóa thành viên"}
	ErrInvalidMembershipData      = AppError{Code: "INVALID_MEMBERSHIP_DATA", Message: "Dữ liệu thành viên không hợp lệ"}
	ErrMembershipAlreadyProcessed = AppError{Code: "MEMBERSHIP_ALREADY_PROCESSED", Message: "Yêu cầu tham gia đã được xử lý"}

	// Post-related
	ErrPostNotFound    = AppError{Code: "POST_NOT_FOUND", Message: "Không tìm thấy bài viết"}
	ErrVoteNotFound    = AppError{Code: "VOTE_NOT_FOUND", Message: "Không tìm thấy bình chọn"}
	ErrPollVoted       = AppError{Code: "POLL_ALREADY_VOTED", Message: "Bạn đã bình chọn lựa chọn này rồi"}
	ErrPollCannotEdit  = AppError{Code: "POLL_CANNOT_EDIT", Message: "Không thể chỉnh sửa bình chọn sau khi đã có người bình chọn"}
	ErrAlreadyReported = AppError{Code: "ALREADY_REPORTED", Message: "Bạn đã báo cáo nội dung này rồi"}
	ErrDraftNotFound   = AppError{Code: "DRAFT_NOT_FOUND", Message: "Không tìm thấy bản nháp"}

	// Comment-related
	ErrCommentNotFound = AppError{Code: "COMMENT_NOT_FOUND", Message: "Không tìm thấy bình luận"}
	ErrDepthInvalid    = AppError{Code: "DEPTH_TOO_HIGH", Message: "Độ sâu phải từ 0 đến 2"}
	ErrUserIsMuted     = AppError{Code: "USER_MUTED", Message: "Người dùng đã bị cấm bình luận"}

	// Messaging-related
	ErrChannelNotFound = AppError{Code: "CHANNEL_NOT_FOUND", Message: "Không tìm thấy kênh"}
	ErrNoMessageFound  = AppError{Code: "NO_MESSAGE_FOUND", Message: "Không tìm thấy tin nhắn"}

	// Report-related
	ErrReportNotFound = AppError{Code: "REPORT_NOT_FOUND", Message: "Không tìm thấy report"}
)
