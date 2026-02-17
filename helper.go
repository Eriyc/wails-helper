// Package helper provides a simplified, opinionated framework for building Wails3 applications.
// It abstracts common patterns like window management, dialogs, keyboard events, auto-updates,
// and license management into easy-to-use interfaces with pluggable providers.
package helper

import (
	"context"

	"github.com/Eriyc/wails-helper/pkg/dialog"
	"github.com/Eriyc/wails-helper/pkg/events"
	"github.com/Eriyc/wails-helper/pkg/providers/autoupdate"
	"github.com/Eriyc/wails-helper/pkg/providers/license"
	"github.com/Eriyc/wails-helper/pkg/window"
)

// Config holds the configuration for the Wails helper
type Config struct {
	// AppName is the name of the application
	AppName string

	// AppVersion is the current version of the application
	AppVersion string

	// Window configuration
	WindowManager window.Manager

	// Dialog configuration
	DialogManager dialog.Manager

	// Event configuration
	EventManager events.Manager

	// Auto-update provider (optional)
	UpdateProvider autoupdate.Provider

	// License provider (optional)
	LicenseProvider license.Provider
}

// Helper is the main helper struct that provides access to all functionality
type Helper struct {
	config Config
}

// New creates a new Wails helper instance
func New(config Config) *Helper {
	return &Helper{
		config: config,
	}
}

// Windows returns the window manager
func (h *Helper) Windows() window.Manager {
	return h.config.WindowManager
}

// Dialogs returns the dialog manager
func (h *Helper) Dialogs() dialog.Manager {
	return h.config.DialogManager
}

// Events returns the event manager
func (h *Helper) Events() events.Manager {
	return h.config.EventManager
}

// Updates returns the auto-update provider
func (h *Helper) Updates() autoupdate.Provider {
	return h.config.UpdateProvider
}

// License returns the license provider
func (h *Helper) License() license.Provider {
	return h.config.LicenseProvider
}

// AppName returns the application name
func (h *Helper) AppName() string {
	return h.config.AppName
}

// AppVersion returns the application version
func (h *Helper) AppVersion() string {
	return h.config.AppVersion
}

// CheckForUpdates is a convenience method to check for updates
func (h *Helper) CheckForUpdates(ctx context.Context) (*autoupdate.UpdateInfo, error) {
	if h.config.UpdateProvider == nil {
		return nil, ErrProviderNotConfigured
	}
	return h.config.UpdateProvider.CheckForUpdates(ctx, h.config.AppVersion)
}

// ValidateLicense is a convenience method to validate a license
func (h *Helper) ValidateLicense(ctx context.Context, licenseKey string) (*license.ValidationResult, error) {
	if h.config.LicenseProvider == nil {
		return nil, ErrProviderNotConfigured
	}
	return h.config.LicenseProvider.Validate(ctx, licenseKey)
}

// ShowInfo is a convenience method to show an info dialog
func (h *Helper) ShowInfo(ctx context.Context, title, message string) error {
	if h.config.DialogManager == nil {
		return ErrDialogManagerNotConfigured
	}
	return h.config.DialogManager.ShowInfo(ctx, title, message)
}

// ShowError is a convenience method to show an error dialog
func (h *Helper) ShowError(ctx context.Context, title, message string) error {
	if h.config.DialogManager == nil {
		return ErrDialogManagerNotConfigured
	}
	return h.config.DialogManager.ShowError(ctx, title, message)
}

// ShowQuestion is a convenience method to show a question dialog
func (h *Helper) ShowQuestion(ctx context.Context, title, message string) (bool, error) {
	if h.config.DialogManager == nil {
		return false, ErrDialogManagerNotConfigured
	}
	return h.config.DialogManager.ShowQuestion(ctx, title, message)
}
