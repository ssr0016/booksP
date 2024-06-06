package errors

import (
	"fmt"
	"net/http"
	"practice/api/meta"
)

type StatusError struct {
	ErrStatusInvalid meta.Status
}

func (e *StatusError) Error() string {
	return e.ErrStatusInvalid.Message
}

func (e *StatusError) Status() meta.Status {
	return e.ErrStatusInvalid
}

func NewUnauthorized(reason string) *StatusError {
	message := reason
	if len(message) == 0 {
		message = "Not authorized"
	}
	return &StatusError{
		meta.Status{
			Status:  meta.StatusFailure,
			Code:    http.StatusUnauthorized,
			Reason:  meta.StatusReasonUnauthorized,
			Message: message,
		},
	}
}

func NewBadRequest(message string) *StatusError {
	if len(message) == 0 {
		message = "bad request"
	}
	return &StatusError{
		meta.Status{
			Status:  meta.StatusFailure,
			Code:    http.StatusOK,
			Reason:  meta.StatusReasonBadRequest,
			Message: message,
		},
	}
}

func NewNotFound(name string) *StatusError {
	return &StatusError{
		meta.Status{
			Status:  meta.StatusFailure,
			Code:    http.StatusOK,
			Reason:  meta.StatusReasonNotFound,
			Message: fmt.Sprintf("%s not found", name),
		},
	}
}

func NotFound() *StatusError {
	return &StatusError{
		meta.Status{
			Status:  meta.StatusFailure,
			Code:    http.StatusOK,
			Reason:  meta.StatusReasonNotFound,
			Message: fmt.Sprintf("Not found"),
		},
	}
}

func NewAlreadyExists(name string) *StatusError {
	return &StatusError{
		meta.Status{
			Status:  meta.StatusFailure,
			Code:    http.StatusOK,
			Reason:  meta.StatusReasonAlreadyExists,
			Message: fmt.Sprintf("%s already exists", name),
		},
	}
}

func NewInternalError(err error) *StatusError {
	return &StatusError{
		meta.Status{
			Status:  meta.StatusFailure,
			Code:    http.StatusOK,
			Reason:  meta.StatusReasonInternalError,
			Message: fmt.Sprintf("Internal error occurred: %v", err),
		},
	}
}

// func New(code int32, message string) *StatusError {
// 	return &StatusError{
// 		meta.Status{
// 			Status:  meta.StatusFailure,
// 			Code:    code,
// 			Reason:  meta.StatusReasonInternalError,
// 			Message: message,
// 		},
// 	}
// }

func NewInvalid(name string) *StatusError {
	return &StatusError{
		meta.Status{
			Status:  meta.StatusFailure,
			Code:    http.StatusOK,
			Reason:  meta.StatusReasonInvalid,
			Message: fmt.Sprintf("%s is invalid", name),
		},
	}
}

func NewStrError(str string) *StatusError {
	return &StatusError{
		meta.Status{
			Status:  meta.StatusFailure,
			Code:    http.StatusOK,
			Reason:  meta.StatusReasonUnknown,
			Message: str,
		},
	}
}

func NewStrErrorCode(str string, errorCode string) *StatusError {
	return &StatusError{
		meta.Status{
			Status:    meta.StatusFailure,
			Code:      http.StatusOK,
			Reason:    meta.StatusReasonUnknown,
			Message:   str,
			ErrorCode: errorCode,
		},
	}
}

// TO DO

type ErrorStatus struct {
	Code    string `json:"code"`
	Message string `json:"message,omitempty"`
	Err     string `json:"-"`
}

func (e ErrorStatus) Error() string {
	return e.Err
}

const (
	ErrGeneral string = "G0001"
)

var (
	ErrBadRequest     = New("G0002", "Bad request")
	ErrUnauthorized   = New("G0003", "Unauthorized")
	ErrForbidden      = New("G0004", "Permission denied")
	ErrSessionExpired = New("G0005", "Session has expired")
)

func New(code string, message string) ErrorStatus {
	return ErrorStatus{
		Code:    code,
		Message: message,
		Err:     message,
	}
}
