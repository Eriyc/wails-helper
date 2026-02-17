package github

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"time"

	"github.com/Eriyc/wails-helper/pkg/providers/license"
)

// Config holds the configuration for the GitHub license provider
type Config struct {
	Owner      string // GitHub repository owner
	Repo       string // GitHub repository name
	Token      string // GitHub API token (required for private repos)
	HTTPClient *http.Client
}

// Provider implements the license.Provider interface using GitHub as a backend
// This can store license data in GitHub issues, releases, or a separate data store
type Provider struct {
	config Config
}

// NewProvider creates a new GitHub license provider
func NewProvider(config Config) *Provider {
	if config.HTTPClient == nil {
		config.HTTPClient = &http.Client{
			Timeout: 30 * time.Second,
		}
	}
	return &Provider{config: config}
}

type licenseData struct {
	License   *license.License `json:"license"`
	Valid     bool             `json:"valid"`
	Message   string           `json:"message,omitempty"`
}

// Validate validates a license key
// This implementation uses GitHub Gists or a custom endpoint to validate licenses
func (p *Provider) Validate(ctx context.Context, licenseKey string) (*license.ValidationResult, error) {
	// For GitHub implementation, we could use GitHub Gists or a custom API endpoint
	// This is a simplified mock implementation
	_ = p.getBaseURL()

	req, err := http.NewRequestWithContext(ctx, "GET", "https://api.github.com", nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	_ = req

	// For demonstration purposes, we'll return a mock validation
	// In a real implementation, this would call your license validation API
	result := &license.ValidationResult{
		Valid: p.isValidLicenseFormat(licenseKey),
		ValidUntil: time.Now().AddDate(1, 0, 0), // 1 year from now
	}

	if result.Valid {
		result.License = &license.License{
			Key:       licenseKey,
			Type:      "standard",
			IssuedAt:  time.Now().AddDate(-1, 0, 0),
			Features:  []string{"core", "updates"},
		}
	} else {
		result.Error = "Invalid license key format"
	}

	return result, nil
}

// Activate activates a license for a specific device
func (p *Provider) Activate(ctx context.Context, req *license.ActivationRequest) (*license.ValidationResult, error) {
	// Hash the device ID for privacy
	deviceHash := p.hashDeviceID(req.DeviceID)
	
	// In a real implementation, this would register the activation with your backend
	// For now, we'll validate and return a result
	validationResult, err := p.Validate(ctx, req.LicenseKey)
	if err != nil {
		return nil, fmt.Errorf("validating license: %w", err)
	}

	if !validationResult.Valid {
		return validationResult, nil
	}

	// Store activation info (in real implementation)
	_ = deviceHash

	return validationResult, nil
}

// Deactivate deactivates a license from a device
func (p *Provider) Deactivate(ctx context.Context, licenseKey, deviceID string) error {
	deviceHash := p.hashDeviceID(deviceID)
	_ = deviceHash
	
	// In a real implementation, this would remove the activation from your backend
	return nil
}

// CheckFeature checks if a specific feature is enabled for a license
func (p *Provider) CheckFeature(ctx context.Context, licenseKey, feature string) (bool, error) {
	result, err := p.Validate(ctx, licenseKey)
	if err != nil {
		return false, err
	}

	if !result.Valid || result.License == nil {
		return false, nil
	}

	for _, f := range result.License.Features {
		if f == feature {
			return true, nil
		}
	}

	return false, nil
}

// RefreshLicense refreshes license information from the server
func (p *Provider) RefreshLicense(ctx context.Context, licenseKey string) (*license.License, error) {
	result, err := p.Validate(ctx, licenseKey)
	if err != nil {
		return nil, err
	}

	if !result.Valid {
		return nil, fmt.Errorf("invalid license")
	}

	return result.License, nil
}

func (p *Provider) getBaseURL() string {
	// In a real implementation, this would point to your license validation API
	// which could be hosted on GitHub Pages, a separate server, etc.
	return fmt.Sprintf("https://api.github.com/repos/%s/%s", p.config.Owner, p.config.Repo)
}

func (p *Provider) isValidLicenseFormat(licenseKey string) bool {
	// Basic validation - check if key has expected format
	// In a real implementation, you'd have more sophisticated validation
	return len(licenseKey) >= 16
}

func (p *Provider) hashDeviceID(deviceID string) string {
	hash := sha256.Sum256([]byte(deviceID))
	return hex.EncodeToString(hash[:])
}
