package custom

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/Eriyc/wails-helper/pkg/providers/license"
)

// Config holds the configuration for the custom server license provider
type Config struct {
	BaseURL    string // Base URL of the custom license server
	APIKey     string // API key for authentication
	HTTPClient *http.Client
}

// Provider implements the license.Provider interface for custom license servers
type Provider struct {
	config Config
}

// NewProvider creates a new custom server license provider
func NewProvider(config Config) *Provider {
	if config.HTTPClient == nil {
		config.HTTPClient = &http.Client{
			Timeout: 30 * time.Second,
		}
	}
	return &Provider{config: config}
}

// Validate validates a license key
func (p *Provider) Validate(ctx context.Context, licenseKey string) (*license.ValidationResult, error) {
	url := fmt.Sprintf("%s/api/licenses/validate", p.config.BaseURL)

	reqBody := map[string]string{"license_key": licenseKey}
	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("marshaling request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	if p.config.APIKey != "" {
		req.Header.Set("X-API-Key", p.config.APIKey)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := p.config.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("validating license: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}

	var result license.ValidationResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	return &result, nil
}

// Activate activates a license for a specific device
func (p *Provider) Activate(ctx context.Context, req *license.ActivationRequest) (*license.ValidationResult, error) {
	url := fmt.Sprintf("%s/api/licenses/activate", p.config.BaseURL)

	bodyBytes, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshaling request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	if p.config.APIKey != "" {
		httpReq.Header.Set("X-API-Key", p.config.APIKey)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")

	resp, err := p.config.HTTPClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("activating license: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}

	var result license.ValidationResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	return &result, nil
}

// Deactivate deactivates a license from a device
func (p *Provider) Deactivate(ctx context.Context, licenseKey, deviceID string) error {
	url := fmt.Sprintf("%s/api/licenses/deactivate", p.config.BaseURL)

	reqBody := map[string]string{
		"license_key": licenseKey,
		"device_id":   deviceID,
	}
	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("marshaling request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(bodyBytes))
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}

	if p.config.APIKey != "" {
		req.Header.Set("X-API-Key", p.config.APIKey)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := p.config.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("deactivating license: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}

	return nil
}

// CheckFeature checks if a specific feature is enabled for a license
func (p *Provider) CheckFeature(ctx context.Context, licenseKey, feature string) (bool, error) {
	url := fmt.Sprintf("%s/api/licenses/features?license_key=%s&feature=%s", p.config.BaseURL, licenseKey, feature)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return false, fmt.Errorf("creating request: %w", err)
	}

	if p.config.APIKey != "" {
		req.Header.Set("X-API-Key", p.config.APIKey)
	}
	req.Header.Set("Accept", "application/json")

	resp, err := p.config.HTTPClient.Do(req)
	if err != nil {
		return false, fmt.Errorf("checking feature: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, nil
	}

	var result struct {
		Enabled bool `json:"enabled"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return false, fmt.Errorf("decoding response: %w", err)
	}

	return result.Enabled, nil
}

// RefreshLicense refreshes license information from the server
func (p *Provider) RefreshLicense(ctx context.Context, licenseKey string) (*license.License, error) {
	url := fmt.Sprintf("%s/api/licenses/%s", p.config.BaseURL, licenseKey)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	if p.config.APIKey != "" {
		req.Header.Set("X-API-Key", p.config.APIKey)
	}
	req.Header.Set("Accept", "application/json")

	resp, err := p.config.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("refreshing license: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}

	var lic license.License
	if err := json.NewDecoder(resp.Body).Decode(&lic); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	return &lic, nil
}
