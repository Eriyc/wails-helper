package helper

import "errors"

var (
	// ErrProviderNotConfigured is returned when a provider is not configured
	ErrProviderNotConfigured = errors.New("provider not configured")

	// ErrDialogManagerNotConfigured is returned when the dialog manager is not configured
	ErrDialogManagerNotConfigured = errors.New("dialog manager not configured")

	// ErrWindowManagerNotConfigured is returned when the window manager is not configured
	ErrWindowManagerNotConfigured = errors.New("window manager not configured")

	// ErrEventManagerNotConfigured is returned when the event manager is not configured
	ErrEventManagerNotConfigured = errors.New("event manager not configured")
)
