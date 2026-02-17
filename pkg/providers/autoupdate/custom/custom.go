package custom

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/Eriyc/wails-helper/pkg/providers/autoupdate"
)

// Config holds the configuration for the custom server auto-update provider
type Config struct {
	BaseURL    string // Base URL of the custom update server
	APIKey     string // API key for authentication (optional)
	HTTPClient *http.Client
}

// Provider implements the autoupdate.Provider interface for custom update servers
type Provider struct {
	config Config
}

// NewProvider creates a new custom server auto-update provider
func NewProvider(config Config) *Provider {
	if config.HTTPClient == nil {
		config.HTTPClient = &http.Client{
			Timeout: 30 * time.Second,
		}
	}
	return &Provider{config: config}
}

type updateResponse struct {
	Available      bool                   `json:"available"`
	CurrentVersion string                 `json:"current_version"`
	LatestVersion  *autoupdate.Version    `json:"latest_version,omitempty"`
}

// CheckForUpdates checks for updates from the custom server
func (p *Provider) CheckForUpdates(ctx context.Context, currentVersion string) (*autoupdate.UpdateInfo, error) {
	url := fmt.Sprintf("%s/api/updates/check?version=%s", p.config.BaseURL, currentVersion)

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
		return nil, fmt.Errorf("checking for updates: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var updateResp updateResponse
	if err := json.NewDecoder(resp.Body).Decode(&updateResp); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	return &autoupdate.UpdateInfo{
		Available:      updateResp.Available,
		CurrentVersion: currentVersion,
		LatestVersion:  updateResp.LatestVersion,
	}, nil
}

// DownloadUpdate downloads the update from the custom server
func (p *Provider) DownloadUpdate(ctx context.Context, version *autoupdate.Version) (io.ReadCloser, error) {
	if version.DownloadURL == "" {
		return nil, fmt.Errorf("no download URL available")
	}

	req, err := http.NewRequestWithContext(ctx, "GET", version.DownloadURL, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	if p.config.APIKey != "" {
		req.Header.Set("X-API-Key", p.config.APIKey)
	}

	resp, err := p.config.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("downloading update: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return resp.Body, nil
}

// VerifyUpdate verifies the checksum of the downloaded update
func (p *Provider) VerifyUpdate(ctx context.Context, data io.Reader, version *autoupdate.Version) error {
	// Basic implementation - in production, implement actual verification
	if version.Checksum == "" {
		return nil
	}
	// TODO: Implement checksum verification
	// For now, return an error to indicate this feature is not yet implemented
	return fmt.Errorf("checksum verification not yet implemented")
}

// GetVersionHistory retrieves the version history from the custom server
func (p *Provider) GetVersionHistory(ctx context.Context, limit int) ([]*autoupdate.Version, error) {
	if limit <= 0 {
		limit = 10
	}

	url := fmt.Sprintf("%s/api/updates/history?limit=%d", p.config.BaseURL, limit)

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
		return nil, fmt.Errorf("fetching version history: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var versions []*autoupdate.Version
	if err := json.NewDecoder(resp.Body).Decode(&versions); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	return versions, nil
}
