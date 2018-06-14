package apperror

import "errors"

var (
	InternalServerError = errors.New("Internal Server Error")
	InvalidAuthToken    = errors.New("Invalid auth token")
	TokenIsExpired      = errors.New("Token is expired")
	LoginTypeNotExists  = errors.New("Login type not exists")
	AccountNotExists    = errors.New("Account not exists")
)

type ErrorCodes struct {
	HTTPcode   int
	StatusCode int
}

var (
	DefaultErrorCode  = ErrorCodes{400, 100000}
	NotFoundErrorCode = ErrorCodes{0, 0}
)

var GetErrorCodes = map[error]ErrorCodes{
	InternalServerError: ErrorCodes{400, 100100},
	InvalidAuthToken:    ErrorCodes{400, 100111},
	TokenIsExpired:      ErrorCodes{400, 100112},
	LoginTypeNotExists:  ErrorCodes{400, 100120},
	AccountNotExists:    ErrorCodes{400, 200010},
}
