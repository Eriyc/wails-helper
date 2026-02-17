package license

import (
	"context"
	"time"
)

// License represents a software license
type License struct {
	Key        string            `json:"key"`
	Type       string            `json:"type"` // e.g., "trial", "standard", "premium"
	Email      string            `json:"email,omitempty"`
	Name       string            `json:"name,omitempty"`
	IssuedAt   time.Time         `json:"issued_at"`
	ExpiresAt  *time.Time        `json:"expires_at,omitempty"`
	Features   []string          `json:"features,omitempty"`
	Metadata   map[string]string `json:"metadata,omitempty"`
	MaxDevices int               `json:"max_devices,omitempty"`
}

// ValidationResult contains the result of license validation
type ValidationResult struct {
	Valid      bool      `json:"valid"`
	License    *License  `json:"license,omitempty"`
	Error      string    `json:"error,omitempty"`
	ValidUntil time.Time `json:"valid_until,omitempty"`
}

// ActivationRequest represents a license activation request
type ActivationRequest struct {
	LicenseKey string            `json:"license_key"`
	DeviceID   string            `json:"device_id"`
	Email      string            `json:"email,omitempty"`
	Metadata   map[string]string `json:"metadata,omitempty"`
}

// Provider defines the interface for license management providers
type Provider interface {
	// Validate validates a license key
	Validate(ctx context.Context, licenseKey string) (*ValidationResult, error)

	// Activate activates a license for a specific device
	Activate(ctx context.Context, req *ActivationRequest) (*ValidationResult, error)

	// Deactivate deactivates a license from a device
	Deactivate(ctx context.Context, licenseKey, deviceID string) error

	// CheckFeature checks if a specific feature is enabled for a license
	CheckFeature(ctx context.Context, licenseKey, feature string) (bool, error)

	// RefreshLicense refreshes license information from the server
	RefreshLicense(ctx context.Context, licenseKey string) (*License, error)
}
