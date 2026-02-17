package helper

import (
	"errors"
	"fmt"
	"runtime/debug"
)

type ErrorCode string

const (
	ErrCodeUnknown      ErrorCode = "unknown"
	ErrCodeInvalidInput ErrorCode = "invalid_input"
	ErrCodeUnavailable  ErrorCode = "unavailable"
	ErrCodeUnauthorized ErrorCode = "unauthorized"
	ErrCodeConflict     ErrorCode = "conflict"
)

type AppError struct {
	Code      ErrorCode
	Message   string
	Cause     error
	Retryable bool
	Metadata  map[string]string
}

func (e *AppError) Error() string {
	if e == nil {
		return ""
	}

	if e.Message != "" {
		return e.Message
	}

	if e.Code != "" {
		return string(e.Code)
	}

	return string(ErrCodeUnknown)
}

func (e *AppError) Unwrap() error {
	if e == nil {
		return nil
	}

	return e.Cause
}

func Wrap(code ErrorCode, message string, cause error) *AppError {
	if code == "" {
		code = ErrCodeUnknown
	}

	return &AppError{
		Code:    code,
		Message: message,
		Cause:   cause,
	}
}

func IsCode(err error, code ErrorCode) bool {
	if err == nil || code == "" {
		return false
	}

	var appErr *AppError
	if !errors.As(err, &appErr) {
		return false
	}

	return appErr.Code == code
}

func NewUnknown(message string, cause error) *AppError {
	return Wrap(ErrCodeUnknown, message, cause)
}

func NewInvalidInput(message string, cause error) *AppError {
	return Wrap(ErrCodeInvalidInput, message, cause)
}

func NewUnavailable(message string, cause error) *AppError {
	return Wrap(ErrCodeUnavailable, message, cause)
}

func NewUnauthorized(message string, cause error) *AppError {
	return Wrap(ErrCodeUnauthorized, message, cause)
}

func NewConflict(message string, cause error) *AppError {
	return Wrap(ErrCodeConflict, message, cause)
}

func PanicToError(recovered any, metadata map[string]string) *AppError {
	if recovered == nil {
		return nil
	}

	resultMetadata := map[string]string{
		"panic.type":  fmt.Sprintf("%T", recovered),
		"panic.value": fmt.Sprint(recovered),
		"panic.stack": string(debug.Stack()),
	}

	for key, value := range metadata {
		resultMetadata[key] = value
	}

	return &AppError{
		Code:     ErrCodeUnknown,
		Message:  "panic recovered",
		Cause:    panicCause(recovered),
		Metadata: resultMetadata,
	}
}

func RunWithRecovery(fn func(), metadata map[string]string) (err error) {
	defer func() {
		recovered := recover()
		if recovered == nil {
			return
		}

		md := cloneMetadata(metadata)
		if _, ok := md["panic.context"]; !ok {
			md["panic.context"] = "callback"
		}
		err = PanicToError(recovered, md)
	}()

	fn()
	return nil
}

func GoWithRecovery(fn func(), onError func(error), metadata map[string]string) {
	go func() {
		if err := RunWithRecovery(fn, mergeMetadata(metadata, map[string]string{"panic.context": "goroutine"})); err != nil && onError != nil {
			onError(err)
		}
	}()
}

func panicCause(recovered any) error {
	if err, ok := recovered.(error); ok {
		return err
	}

	return errors.New(fmt.Sprint(recovered))
}

func cloneMetadata(metadata map[string]string) map[string]string {
	result := make(map[string]string, len(metadata))
	for key, value := range metadata {
		result[key] = value
	}

	return result
}

func mergeMetadata(base map[string]string, overlay map[string]string) map[string]string {
	result := cloneMetadata(base)
	for key, value := range overlay {
		result[key] = value
	}

	return result
}
