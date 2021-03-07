package enums

import (
	"errors"
	"net/http"
)

var (
	ErrFileNotExist              = errors.New("file not exist")
	ErrInValidConfigTimeDuration = errors.New("invalid config time duration")
)

var (
	ErrMutexNotFound = errors.New("mutex no found")
	ErrMutexExisted = errors.New("mutex already exist")
)

// ResponseError

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (r ErrorResponse) GetCode() int {
	return r.Code
}

func (r ErrorResponse) GetMessage() string {
	return r.Message
}

var (
	ErrorSystem = ErrorResponse{
		Code:    http.StatusInternalServerError,
		Message: "System error, Please contact admin to get more information",
	}
	ErrorNoContent = ErrorResponse{
		Code: http.StatusNoContent,
	}
)
