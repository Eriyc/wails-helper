# Architecture Documentation

## Overview

The Wails Helper library is designed with a provider pattern architecture that enables pluggable implementations for various features. This document explains the architectural decisions and design patterns used.

## Design Principles

1. **Separation of Concerns**: Each package has a single, well-defined responsibility
2. **Interface-Driven Design**: Core functionality is defined through interfaces
3. **Provider Pattern**: Allows swapping implementations without changing client code
4. **Extensibility**: Easy to add new providers or extend existing functionality
5. **Type Safety**: Strong typing throughout the API

## Package Structure

### Core Packages

- **`helper`**: Main entry point providing unified access to all functionality
- **`errors`**: Common error types used across the library

### Feature Packages

- **`pkg/window`**: Window management abstractions
- **`pkg/dialog`**: Dialog system interfaces
- **`pkg/events`**: Event handling and keyboard shortcuts

### Provider Packages

- **`pkg/providers/autoupdate`**: Auto-update provider interface and implementations
  - `github`: GitHub Releases-based updates
  - `custom`: Custom server-based updates
  
- **`pkg/providers/license`**: License management provider interface and implementations
  - `github`: GitHub-based license validation
  - `custom`: Custom server-based license validation

## Provider Pattern

The provider pattern allows users to choose or implement their own backend for features like auto-updates and license management.

### Benefits

1. **Flexibility**: Switch between different implementations at runtime
2. **Testability**: Easy to mock providers for testing
3. **Extensibility**: Add new providers without modifying existing code
4. **Separation**: Business logic is separate from infrastructure concerns

### Example

```go
// Define the interface
type Provider interface {
    CheckForUpdates(ctx context.Context, currentVersion string) (*UpdateInfo, error)
}

// Implement for GitHub
type GitHubProvider struct { ... }
func (p *GitHubProvider) CheckForUpdates(...) { ... }

// Implement for custom server
type CustomProvider struct { ... }
func (p *CustomProvider) CheckForUpdates(...) { ... }

// Use interchangeably
var provider Provider
provider = NewGitHubProvider(...)
// or
provider = NewCustomProvider(...)

updateInfo, err := provider.CheckForUpdates(ctx, "1.0.0")
```

## Interface Design

All major features are defined as interfaces:

1. **window.Manager**: Creates and manages windows
2. **dialog.Manager**: Shows dialogs and file pickers
3. **events.Manager**: Handles keyboard and mouse events
4. **autoupdate.Provider**: Checks and downloads updates
5. **license.Provider**: Validates and manages licenses

This allows:
- Easy mocking for tests
- Multiple implementations
- Future extensibility without breaking changes

## Context Usage

All API methods that perform I/O or long-running operations accept a `context.Context` parameter. This enables:

- Request cancellation
- Timeout control
- Request-scoped values
- Graceful shutdown

Example:
```go
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

updateInfo, err := provider.CheckForUpdates(ctx, "1.0.0")
```

## Error Handling

The library uses Go's standard error handling with wrapped errors:

```go
return nil, fmt.Errorf("checking for updates: %w", err)
```

This allows error inspection with `errors.Is()` and `errors.As()`.

Common errors are defined in the `errors.go` file:
- `ErrProviderNotConfigured`
- `ErrDialogManagerNotConfigured`
- `ErrWindowManagerNotConfigured`
- `ErrEventManagerNotConfigured`

## Extensibility Points

### Custom Providers

Users can implement their own providers for any feature:

```go
type MyUpdateProvider struct {
    // Custom fields
}

func (p *MyUpdateProvider) CheckForUpdates(ctx context.Context, currentVersion string) (*autoupdate.UpdateInfo, error) {
    // Custom implementation
}

// Use it
h := helper.New(helper.Config{
    UpdateProvider: &MyUpdateProvider{},
})
```

### Custom Implementations

While the library provides interfaces for window, dialog, and event management, users must provide concrete implementations that wrap Wails3's APIs:

```go
type Wails3WindowManager struct {
    app *wails.App
}

func (w *Wails3WindowManager) Create(ctx context.Context, opts window.Options) (window.Window, error) {
    // Wrap Wails3 window creation
}
```

## Future Enhancements

The architecture allows for easy addition of:

1. **New Provider Types**:
   - Analytics providers
   - Telemetry providers
   - Crash reporting providers
   - Authentication providers

2. **New Feature Packages**:
   - System tray management
   - Notification system
   - File system watcher
   - IPC/RPC abstractions

3. **Provider Composition**:
   - Multi-provider fallback chains
   - Provider middleware/interceptors
   - Provider load balancing

## Integration with Wails3

The library is designed to complement Wails3, not replace it. Integration points:

1. **Window Management**: Wrap Wails3's window API
2. **Dialogs**: Wrap Wails3's dialog system
3. **Events**: Bridge Wails3 events to the helper's event system
4. **Bindings**: Use the helper from Wails3's Go bindings

Example integration:
```go
type App struct {
    helper *helper.Helper
    wails  *wails.App
}

func (a *App) Initialize() {
    // Create Wails3 implementations
    windowMgr := NewWails3WindowManager(a.wails)
    dialogMgr := NewWails3DialogManager(a.wails)
    eventMgr := NewWails3EventManager(a.wails)
    
    // Configure helper
    a.helper = helper.New(helper.Config{
        WindowManager:   windowMgr,
        DialogManager:   dialogMgr,
        EventManager:    eventMgr,
        UpdateProvider:  github.NewProvider(...),
        LicenseProvider: custom.NewProvider(...),
    })
}
```

## Thread Safety

Provider implementations should be thread-safe when possible. The interfaces allow concurrent calls, so implementations should handle synchronization internally if needed.

## Best Practices

1. **Use contexts**: Always pass and respect context for cancellation
2. **Handle errors**: Don't ignore errors; wrap them with context
3. **Close resources**: Use `defer` to close io.ReadCloser and other resources
4. **Validate inputs**: Check for nil pointers and invalid values
5. **Document behavior**: Use comments to explain non-obvious behavior

## Testing Strategy

1. **Interface mocking**: Easy to create mock implementations
2. **Provider testing**: Test each provider implementation independently
3. **Integration testing**: Test the helper with real providers
4. **Example testing**: Verify examples compile and run

Example mock:
```go
type MockUpdateProvider struct {
    CheckForUpdatesFunc func(ctx context.Context, currentVersion string) (*autoupdate.UpdateInfo, error)
}

func (m *MockUpdateProvider) CheckForUpdates(ctx context.Context, currentVersion string) (*autoupdate.UpdateInfo, error) {
    if m.CheckForUpdatesFunc != nil {
        return m.CheckForUpdatesFunc(ctx, currentVersion)
    }
    return nil, nil
}
```
