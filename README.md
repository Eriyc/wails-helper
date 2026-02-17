# Wails Helper

A helper Go library that abstracts [Wails3](https://v3alpha.wails.io/) functionality into simple, opinionated interfaces. This library reduces boilerplate and provides a clean API for common tasks like window management, dialogs, keyboard events, auto-updates, and license management.

## Features

- 🪟 **Window Management**: Simple API for creating and managing windows
- 💬 **Dialogs**: Easy-to-use dialog system for messages, file pickers, etc.
- ⌨️ **Keyboard Events**: Handle keyboard shortcuts and events
- 🔄 **Auto-Updates**: Pluggable providers for checking and downloading updates
  - GitHub Releases provider
  - Custom server provider
- 🔐 **License Management**: Flexible license validation and activation
  - GitHub-based provider
  - Custom server provider
- 🔌 **Provider Pattern**: Easy to extend with custom providers

## Installation

```bash
go get github.com/Eriyc/wails-helper
```

## Architecture

The library is built around a provider pattern that allows you to plug in different implementations for various features. This makes it easy to:

- Switch between different update sources (GitHub, custom server, etc.)
- Implement custom license validation logic
- Extend functionality with your own providers

### Roadmap

- [v0.1 Proposal](docs/v0.1-proposal.md)

### Package Structure

```
github.com/Eriyc/wails-helper/
├── helper.go                    # Main helper with convenience methods
├── errors.go                    # Common error types
├── pkg/
│   ├── window/                  # Window management interfaces
│   ├── dialog/                  # Dialog interfaces
│   ├── events/                  # Event system interfaces
│   └── providers/
│       ├── autoupdate/          # Auto-update provider interface
│       │   ├── github/          # GitHub releases provider
│       │   └── custom/          # Custom server provider
│       └── license/             # License provider interface
│           ├── github/          # GitHub-based license provider
│           └── custom/          # Custom server license provider
```

## Usage

### Basic Setup

```go
package main

import (
    "context"
    "github.com/Eriyc/wails-helper"
    "github.com/Eriyc/wails-helper/pkg/providers/autoupdate/github"
    licensegithub "github.com/Eriyc/wails-helper/pkg/providers/license/github"
)

func main() {
    // Create auto-update provider
    updateProvider := github.NewProvider(github.Config{
        Owner: "your-username",
        Repo:  "your-repo",
        Token: "your-github-token", // optional
    })

    // Create license provider
    licenseProvider := licensegithub.NewProvider(licensegithub.Config{
        Owner: "your-username",
        Repo:  "your-licenses-repo",
        Token: "your-github-token",
    })

    // Create helper instance
    h := helper.New(helper.Config{
        AppName:         "My App",
        AppVersion:      "1.0.0",
        UpdateProvider:  updateProvider,
        LicenseProvider: licenseProvider,
        // WindowManager, DialogManager, EventManager would be set here
        // These would typically wrap Wails3's native implementations
    })

    // Use the helper
    ctx := context.Background()

    // Check for updates
    updateInfo, err := h.CheckForUpdates(ctx)
    if err != nil {
        // handle error
    }

    if updateInfo.Available {
        // New version available!
    }
}
```

### Auto-Updates

#### Using GitHub Releases

```go
import (
    "github.com/Eriyc/wails-helper/pkg/providers/autoupdate/github"
)

provider := github.NewProvider(github.Config{
    Owner: "username",
    Repo:  "repo-name",
    Token: "optional-token",
})

// Check for updates
updateInfo, err := provider.CheckForUpdates(ctx, "1.0.0")
if err != nil {
    // handle error
}

if updateInfo.Available {
    // Download update
    reader, err := provider.DownloadUpdate(ctx, updateInfo.LatestVersion)
    if err != nil {
        // handle error
    }
    defer reader.Close()

    // Process the update...
}

// Get version history
versions, err := provider.GetVersionHistory(ctx, 10)
```

#### Using Custom Server

```go
import (
    "github.com/Eriyc/wails-helper/pkg/providers/autoupdate/custom"
)

provider := custom.NewProvider(custom.Config{
    BaseURL: "https://updates.example.com",
    APIKey:  "your-api-key",
})

// Same API as GitHub provider
updateInfo, err := provider.CheckForUpdates(ctx, "1.0.0")
```

### License Management

#### Using GitHub-based Provider

```go
import (
    "github.com/Eriyc/wails-helper/pkg/providers/license/github"
)

provider := github.NewProvider(github.Config{
    Owner: "username",
    Repo:  "licenses",
    Token: "github-token",
})

// Validate license
result, err := provider.Validate(ctx, "LICENSE-KEY-HERE")
if err != nil {
    // handle error
}

if result.Valid {
    // License is valid
    features := result.License.Features
}

// Activate license for a device
activationReq := &license.ActivationRequest{
    LicenseKey: "LICENSE-KEY",
    DeviceID:   "device-unique-id",
    Email:      "user@example.com",
}
result, err = provider.Activate(ctx, activationReq)

// Check specific feature
hasFeature, err := provider.CheckFeature(ctx, "LICENSE-KEY", "premium")
```

#### Using Custom Server

```go
import (
    "github.com/Eriyc/wails-helper/pkg/providers/license/custom"
)

provider := custom.NewProvider(custom.Config{
    BaseURL: "https://license.example.com",
    APIKey:  "your-api-key",
})

// Same API as GitHub provider
result, err := provider.Validate(ctx, "LICENSE-KEY")
```

### Window Management

```go
// Note: Window manager needs to be implemented to wrap Wails3 APIs
// This shows the interface usage

// Create a new window
window, err := h.Windows().Create(ctx, window.Options{
    Title:      "My Window",
    Width:      800,
    Height:     600,
    Resizable:  true,
})

// Manipulate window
window.SetTitle(ctx, "New Title")
window.Maximize(ctx)
window.Center(ctx)

// Execute JavaScript
window.ExecJS(ctx, "console.log('Hello from Go!')")
```

### Dialogs

```go
// Note: Dialog manager needs to be implemented to wrap Wails3 APIs

// Show info dialog
err := h.ShowInfo(ctx, "Success", "Operation completed successfully")

// Show error dialog
err := h.ShowError(ctx, "Error", "Something went wrong")

// Show question dialog
yes, err := h.ShowQuestion(ctx, "Confirm", "Are you sure?")

// Open file dialog
file, err := h.Dialogs().OpenFile(ctx, dialog.OpenDialogOptions{
    Title: "Select a file",
    Filters: []dialog.Filter{
        {DisplayName: "Text Files", Pattern: "*.txt"},
        {DisplayName: "All Files", Pattern: "*.*"},
    },
})
```

### Keyboard Events

```go
// Note: Event manager needs to be implemented to wrap Wails3 APIs

// Register keyboard shortcut
h.Events().RegisterShortcut(ctx, events.Shortcut{
    Key:       "s",
    Modifiers: events.ModifierControl,
    Handler: func(event interface{}) error {
        // Handle Ctrl+S
        return nil
    },
})

// Register global shortcut (system-wide)
h.Events().RegisterGlobalShortcut(ctx, events.Shortcut{
    Key:       "F12",
    Modifiers: events.ModifierNone,
    Handler: func(event interface{}) error {
        // Handle F12
        return nil
    },
})
```

## Implementing Custom Providers

You can implement your own providers by satisfying the provider interfaces:

### Custom Auto-Update Provider

```go
import "github.com/Eriyc/wails-helper/pkg/providers/autoupdate"

type MyUpdateProvider struct {
    // your fields
}

func (p *MyUpdateProvider) CheckForUpdates(ctx context.Context, currentVersion string) (*autoupdate.UpdateInfo, error) {
    // your implementation
}

func (p *MyUpdateProvider) DownloadUpdate(ctx context.Context, version *autoupdate.Version) (io.ReadCloser, error) {
    // your implementation
}

func (p *MyUpdateProvider) VerifyUpdate(ctx context.Context, data io.Reader, version *autoupdate.Version) error {
    // your implementation
}

func (p *MyUpdateProvider) GetVersionHistory(ctx context.Context, limit int) ([]*autoupdate.Version, error) {
    // your implementation
}
```

### Custom License Provider

```go
import "github.com/Eriyc/wails-helper/pkg/providers/license"

type MyLicenseProvider struct {
    // your fields
}

func (p *MyLicenseProvider) Validate(ctx context.Context, licenseKey string) (*license.ValidationResult, error) {
    // your implementation
}

func (p *MyLicenseProvider) Activate(ctx context.Context, req *license.ActivationRequest) (*license.ValidationResult, error) {
    // your implementation
}

func (p *MyLicenseProvider) Deactivate(ctx context.Context, licenseKey, deviceID string) error {
    // your implementation
}

func (p *MyLicenseProvider) CheckFeature(ctx context.Context, licenseKey, feature string) (bool, error) {
    // your implementation
}

func (p *MyLicenseProvider) RefreshLicense(ctx context.Context, licenseKey string) (*license.License, error) {
    // your implementation
}
```

## Custom Server API Specifications

### Auto-Update Server API

Your custom update server should implement these endpoints:

#### Check for Updates

```
GET /api/updates/check?version=1.0.0
Response: {
  "available": true,
  "current_version": "1.0.0",
  "latest_version": {
    "version": "1.1.0",
    "release_date": "2024-01-01T00:00:00Z",
    "notes": "Release notes",
    "download_url": "https://...",
    "checksum": "sha256:..."
  }
}
```

#### Get Version History

```
GET /api/updates/history?limit=10
Response: [
  {
    "version": "1.1.0",
    "release_date": "2024-01-01T00:00:00Z",
    "notes": "Release notes",
    "download_url": "https://..."
  }
]
```

### License Server API

Your custom license server should implement these endpoints:

#### Validate License

```
POST /api/licenses/validate
Body: {"license_key": "..."}
Response: {
  "valid": true,
  "license": {...},
  "valid_until": "2025-01-01T00:00:00Z"
}
```

#### Activate License

```
POST /api/licenses/activate
Body: {
  "license_key": "...",
  "device_id": "...",
  "email": "..."
}
```

#### Deactivate License

```
POST /api/licenses/deactivate
Body: {
  "license_key": "...",
  "device_id": "..."
}
```

#### Check Feature

```
GET /api/licenses/features?license_key=...&feature=premium
Response: {"enabled": true}
```

#### Refresh License

```
GET /api/licenses/{license_key}
Response: {
  "key": "...",
  "type": "premium",
  "features": ["core", "premium"]
}
```

## Integrating with Wails3

To fully integrate with Wails3, you'll need to implement the `window.Manager`, `dialog.Manager`, and `events.Manager` interfaces using Wails3's native APIs. These implementations will wrap Wails3 functionality and expose it through the helper's clean interface.

Pin the Wails3 CLI to the supported version:

```bash
go install github.com/wailsapp/wails/v3/cmd/wails3@v3.0.0-alpha.72
```

Example skeleton:

```go
type Wails3WindowManager struct {
    app *application.App
}

func (w *Wails3WindowManager) Create(ctx context.Context, opts window.Options) (window.Window, error) {
    // Use Wails3 API to create window
}

// ... implement other methods
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is released into the public domain under the Unlicense. See LICENSE file for details.
