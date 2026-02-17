package autoupdate

import (
	"context"
	"io"
)

// Version represents a software version with metadata
type Version struct {
	Version     string            `json:"version"`
	ReleaseDate string            `json:"release_date,omitempty"`
	Notes       string            `json:"notes,omitempty"`
	DownloadURL string            `json:"download_url"`
	Checksum    string            `json:"checksum,omitempty"`
	Metadata    map[string]string `json:"metadata,omitempty"`
}

// UpdateInfo contains information about an available update
type UpdateInfo struct {
	Available      bool     `json:"available"`
	CurrentVersion string   `json:"current_version"`
	LatestVersion  *Version `json:"latest_version,omitempty"`
}

// Provider defines the interface for auto-update providers
type Provider interface {
	// CheckForUpdates checks if a new version is available
	CheckForUpdates(ctx context.Context, currentVersion string) (*UpdateInfo, error)

	// DownloadUpdate downloads the update package
	DownloadUpdate(ctx context.Context, version *Version) (io.ReadCloser, error)

	// VerifyUpdate verifies the integrity of a downloaded update
	VerifyUpdate(ctx context.Context, data io.Reader, version *Version) error

	// GetVersionHistory retrieves the version history
	GetVersionHistory(ctx context.Context, limit int) ([]*Version, error)
}

// DownloadProgress represents download progress information
type DownloadProgress struct {
	BytesDownloaded int64
	TotalBytes      int64
	Percentage      float64
}

// ProgressCallback is called during download to report progress
type ProgressCallback func(progress *DownloadProgress)
