package github

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/Eriyc/wails-helper/pkg/providers/autoupdate"
)

// Config holds the configuration for the GitHub auto-update provider
type Config struct {
	Owner      string // GitHub repository owner
	Repo       string // GitHub repository name
	Token      string // GitHub API token (optional, for private repos or higher rate limits)
	HTTPClient *http.Client
}

// Provider implements the autoupdate.Provider interface for GitHub releases
type Provider struct {
	config Config
}

// NewProvider creates a new GitHub auto-update provider
func NewProvider(config Config) *Provider {
	if config.HTTPClient == nil {
		config.HTTPClient = &http.Client{
			Timeout: 30 * time.Second,
		}
	}
	return &Provider{config: config}
}

type githubRelease struct {
	TagName     string `json:"tag_name"`
	Name        string `json:"name"`
	PublishedAt string `json:"published_at"`
	Body        string `json:"body"`
	Assets      []struct {
		Name               string `json:"name"`
		BrowserDownloadURL string `json:"browser_download_url"`
	} `json:"assets"`
	Prerelease bool `json:"prerelease"`
	Draft      bool `json:"draft"`
}

// CheckForUpdates checks for updates from GitHub releases
func (p *Provider) CheckForUpdates(ctx context.Context, currentVersion string) (*autoupdate.UpdateInfo, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", p.config.Owner, p.config.Repo)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	if p.config.Token != "" {
		req.Header.Set("Authorization", "Bearer "+p.config.Token)
	}
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := p.config.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetching releases: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var release githubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	latestVersion := strings.TrimPrefix(release.TagName, "v")
	currentVersion = strings.TrimPrefix(currentVersion, "v")

	updateInfo := &autoupdate.UpdateInfo{
		Available:      latestVersion != currentVersion,
		CurrentVersion: currentVersion,
	}

	if updateInfo.Available {
		downloadURL := ""
		if len(release.Assets) > 0 {
			downloadURL = release.Assets[0].BrowserDownloadURL
		}

		updateInfo.LatestVersion = &autoupdate.Version{
			Version:     latestVersion,
			ReleaseDate: release.PublishedAt,
			Notes:       release.Body,
			DownloadURL: downloadURL,
		}
	}

	return updateInfo, nil
}

// DownloadUpdate downloads the update from GitHub
func (p *Provider) DownloadUpdate(ctx context.Context, version *autoupdate.Version) (io.ReadCloser, error) {
	if version.DownloadURL == "" {
		return nil, fmt.Errorf("no download URL available")
	}

	req, err := http.NewRequestWithContext(ctx, "GET", version.DownloadURL, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	if p.config.Token != "" {
		req.Header.Set("Authorization", "Bearer "+p.config.Token)
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
	// Basic implementation - in production, you'd verify checksums
	if version.Checksum == "" {
		return nil // No checksum to verify
	}
	// TODO: Implement checksum verification
	// For now, return an error to indicate this feature is not yet implemented
	return fmt.Errorf("checksum verification not yet implemented")
}

// GetVersionHistory retrieves the version history from GitHub releases
func (p *Provider) GetVersionHistory(ctx context.Context, limit int) ([]*autoupdate.Version, error) {
	if limit <= 0 {
		limit = 10
	}

	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases?per_page=%d", p.config.Owner, p.config.Repo, limit)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	if p.config.Token != "" {
		req.Header.Set("Authorization", "Bearer "+p.config.Token)
	}
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := p.config.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetching releases: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var releases []githubRelease
	if err := json.NewDecoder(resp.Body).Decode(&releases); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	versions := make([]*autoupdate.Version, 0, len(releases))
	for _, release := range releases {
		if release.Draft {
			continue
		}

		downloadURL := ""
		if len(release.Assets) > 0 {
			downloadURL = release.Assets[0].BrowserDownloadURL
		}

		versions = append(versions, &autoupdate.Version{
			Version:     strings.TrimPrefix(release.TagName, "v"),
			ReleaseDate: release.PublishedAt,
			Notes:       release.Body,
			DownloadURL: downloadURL,
		})
	}

	return versions, nil
}
