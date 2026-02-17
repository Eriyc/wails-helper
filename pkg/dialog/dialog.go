package dialog

import (
	"context"
)

// MessageType represents the type of message dialog
type MessageType string

const (
	// MessageTypeInfo represents an info message
	MessageTypeInfo MessageType = "info"
	// MessageTypeWarning represents a warning message
	MessageTypeWarning MessageType = "warning"
	// MessageTypeError represents an error message
	MessageTypeError MessageType = "error"
	// MessageTypeQuestion represents a question message
	MessageTypeQuestion MessageType = "question"
)

// MessageOptions represents options for a message dialog
type MessageOptions struct {
	Type          MessageType `json:"type"`
	Title         string      `json:"title"`
	Message       string      `json:"message"`
	Buttons       []string    `json:"buttons,omitempty"`       // Custom button labels
	DefaultButton int         `json:"default_button,omitempty"` // Index of default button
	CancelButton  int         `json:"cancel_button,omitempty"`  // Index of cancel button
}

// OpenDialogOptions represents options for an open file/folder dialog
type OpenDialogOptions struct {
	Title                string   `json:"title,omitempty"`
	DefaultDirectory     string   `json:"default_directory,omitempty"`
	DefaultFilename      string   `json:"default_filename,omitempty"`
	Filters              []Filter `json:"filters,omitempty"`
	ShowHiddenFiles      bool     `json:"show_hidden_files,omitempty"`
	CanCreateDirectories bool     `json:"can_create_directories,omitempty"`
	ResolvesAliases      bool     `json:"resolves_aliases,omitempty"`
	AllowsMultipleSelection bool  `json:"allows_multiple_selection,omitempty"`
	CanChooseDirectories bool     `json:"can_choose_directories,omitempty"`
	CanChooseFiles       bool     `json:"can_choose_files,omitempty"`
}

// SaveDialogOptions represents options for a save file dialog
type SaveDialogOptions struct {
	Title                string   `json:"title,omitempty"`
	DefaultDirectory     string   `json:"default_directory,omitempty"`
	DefaultFilename      string   `json:"default_filename,omitempty"`
	Filters              []Filter `json:"filters,omitempty"`
	ShowHiddenFiles      bool     `json:"show_hidden_files,omitempty"`
	CanCreateDirectories bool     `json:"can_create_directories,omitempty"`
	TreatPackagesAsDirectories bool `json:"treat_packages_as_directories,omitempty"`
}

// Filter represents a file filter for dialogs
type Filter struct {
	DisplayName string   `json:"display_name"` // e.g. "Images"
	Pattern     string   `json:"pattern"`       // e.g. "*.png;*.jpg"
}

// Manager defines the interface for dialog management
type Manager interface {
	// ShowMessage shows a message dialog
	ShowMessage(ctx context.Context, opts MessageOptions) (string, error)

	// ShowInfo shows an info message dialog
	ShowInfo(ctx context.Context, title, message string) error

	// ShowWarning shows a warning message dialog
	ShowWarning(ctx context.Context, title, message string) error

	// ShowError shows an error message dialog
	ShowError(ctx context.Context, title, message string) error

	// ShowQuestion shows a question dialog with Yes/No buttons
	ShowQuestion(ctx context.Context, title, message string) (bool, error)

	// OpenFile shows an open file dialog
	OpenFile(ctx context.Context, opts OpenDialogOptions) (string, error)

	// OpenFiles shows an open multiple files dialog
	OpenFiles(ctx context.Context, opts OpenDialogOptions) ([]string, error)

	// OpenDirectory shows an open directory dialog
	OpenDirectory(ctx context.Context, opts OpenDialogOptions) (string, error)

	// SaveFile shows a save file dialog
	SaveFile(ctx context.Context, opts SaveDialogOptions) (string, error)
}
