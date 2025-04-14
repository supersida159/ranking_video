package apperror

import (
	"net/http"
)

type ErrorCode int

const (
	// Bad Request Errors (400)
	CodeInvalidMessageKey ErrorCode = -100 + iota
	CodeUsernameTooShort
	CodePasswordMustContainNumber
	CodeInvalidEmail
	CodeUserNotExist
	CodeJSONBindingError

	// Database Errors (500)
	CodeDBInternal
	CodeRecordNotFound
)

type AppError struct {
	StatusCode   int       `json:"status_code"`
	RootErr      error     `json:"-"`
	MessageEn    string    `json:"message_en"`
	MessageVi    string    `json:"message_vi"`
	Key          string    `json:"error_key"`
	InternalCode ErrorCode `json:"code"`
}

func (a *AppError) Error() string {
	if a.RootErr != nil {
		return a.RootErr.Error()
	}
	return a.MessageEn
}

func NewErrorResponse(rootErr error, messageEn, messageVi, key string, internalCode ErrorCode, statusCode int) *AppError {
	return &AppError{
		RootErr:      rootErr,
		MessageEn:    messageEn,
		MessageVi:    messageVi,
		Key:          key,
		StatusCode:   statusCode,
		InternalCode: internalCode,
	}
}

// Bad Request Errors (400)
func ErrInvalidMessageKey(err ...error) *AppError {
	var rootErr error
	if len(err) > 0 {
		rootErr = err[0]
	}
	return NewErrorResponse(
		rootErr,
		"Invalid message key",
		"Message Key không hợp lệ",
		"INVALID_MESSAGE_KEY",
		CodeInvalidMessageKey,
		http.StatusBadRequest,
	)
}

func ErrUsernameTooShort(err ...error) *AppError {
	var rootErr error
	if len(err) > 0 {
		rootErr = err[0]
	}
	return NewErrorResponse(
		rootErr,
		"Username must be at least 8 characters",
		"Tên người dùng phải có ít nhất 8 ký tự",
		"USERNAME_TOO_SHORT",
		CodeUsernameTooShort,
		http.StatusBadRequest,
	)
}

func ErrPasswordMustContainNumber(err ...error) *AppError {
	var rootErr error
	if len(err) > 0 {
		rootErr = err[0]
	}
	return NewErrorResponse(
		rootErr,
		"Password must contain at least one number",
		"Mật khẩu phải chứa ít nhất một số",
		"PASSWORD_MUST_CONTAIN_NUMBER",
		CodePasswordMustContainNumber,
		http.StatusBadRequest,
	)
}

func ErrInvalidEmail(err ...error) *AppError {
	var rootErr error
	if len(err) > 0 {
		rootErr = err[0]
	}
	return NewErrorResponse(
		rootErr,
		"Email is invalid",
		"Email không hợp lệ",
		"INVALID_EMAIL",
		CodeInvalidEmail,
		http.StatusBadRequest,
	)
}

func ErrUserNotExist(err ...error) *AppError {
	var rootErr error
	if len(err) > 0 {
		rootErr = err[0]
	}
	return NewErrorResponse(
		rootErr,
		"User does not exist",
		"Người dùng không tồn tại",
		"USER_NOT_EXIST",
		CodeUserNotExist,
		http.StatusBadRequest,
	)
}

func ErrJSONBinding(err ...error) *AppError {
	var rootErr error
	if len(err) > 0 {
		rootErr = err[0]
	}
	return NewErrorResponse(
		rootErr,
		"Invalid JSON format",
		"Định dạng JSON không hợp lệ",
		"INVALID_JSON_FORMAT",
		CodeJSONBindingError,
		http.StatusBadRequest,
	)
}

// Database Errors
func ErrDBInternal(err ...error) *AppError {
	var rootErr error
	if len(err) > 0 {
		rootErr = err[0]
	}
	return NewErrorResponse(
		rootErr,
		"Internal database error",
		"Lỗi cơ sở dữ liệu nội bộ",
		"DB_INTERNAL_ERROR",
		CodeDBInternal,
		http.StatusInternalServerError,
	)
}

func ErrRecordNotFound(err ...error) *AppError {
	var rootErr error
	if len(err) > 0 {
		rootErr = err[0]
	}
	return NewErrorResponse(
		rootErr,
		"Record not found",
		"Không tìm thấy bản ghi",
		"RECORD_NOT_FOUND",
		CodeRecordNotFound,
		http.StatusNotFound,
	)
}
