# Contributing to Wails Helper

Thank you for your interest in contributing to Wails Helper! This document provides guidelines and information for contributors.

## Getting Started

### Prerequisites

- Go 1.20 or later
- Git
- (Optional) Nix with flakes enabled for the development environment

### Setting Up Development Environment

1. **Clone the repository**:
   ```bash
   git clone https://github.com/Eriyc/wails-helper.git
   cd wails-helper
   ```

2. **With Nix (recommended)**:
   ```bash
   nix develop
   ```

3. **Without Nix**:
   Ensure you have Go installed:
   ```bash
   go version
   ```

4. **Install dependencies**:
   ```bash
   go mod download
   ```

## Project Structure

```
wails-helper/
├── helper.go              # Main helper entry point
├── errors.go              # Common error types
├── pkg/
│   ├── window/           # Window management interfaces
│   ├── dialog/           # Dialog interfaces
│   ├── events/           # Event system interfaces
│   └── providers/
│       ├── autoupdate/   # Auto-update provider interface
│       │   ├── github/   # GitHub implementation
│       │   └── custom/   # Custom server implementation
│       └── license/      # License provider interface
│           ├── github/   # GitHub implementation
│           └── custom/   # Custom server implementation
├── examples/             # Usage examples
├── README.md            # User documentation
└── ARCHITECTURE.md      # Architecture documentation
```

## Development Workflow

### Making Changes

1. **Create a branch**:
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. **Make your changes** following the code style guidelines

3. **Build and test**:
   ```bash
   go build ./...
   go test ./...
   ```

4. **Format your code**:
   ```bash
   go fmt ./...
   ```

5. **Commit your changes**:
   ```bash
   git commit -m "Add your descriptive commit message"
   ```

6. **Push and create a pull request**:
   ```bash
   git push origin feature/your-feature-name
   ```

## Code Style Guidelines

### Go Code Style

- Follow standard Go conventions and idioms
- Use `gofmt` to format code
- Run `go vet` to catch common mistakes
- Keep functions small and focused
- Use meaningful variable and function names
- Add comments for exported types and functions

### Interface Design

When adding new interfaces:

1. Keep interfaces small and focused (Interface Segregation Principle)
2. Use `context.Context` as the first parameter for I/O operations
3. Return errors as the last return value
4. Document all exported interfaces, types, and functions

Example:
```go
// Manager defines the interface for managing widgets
type Manager interface {
    // Create creates a new widget with the given options
    Create(ctx context.Context, opts Options) (Widget, error)
    
    // Get retrieves a widget by its ID
    Get(ctx context.Context, id string) (Widget, error)
}
```

### Provider Implementation

When implementing a provider:

1. Implement all interface methods
2. Handle context cancellation properly
3. Validate input parameters
4. Return descriptive errors with context
5. Add comprehensive documentation

Example:
```go
// Provider implements the example.Provider interface
type Provider struct {
    config Config
    client *http.Client
}

// NewProvider creates a new provider instance
func NewProvider(config Config) *Provider {
    if config.HTTPClient == nil {
        config.HTTPClient = &http.Client{
            Timeout: 30 * time.Second,
        }
    }
    return &Provider{
        config: config,
        client: config.HTTPClient,
    }
}

func (p *Provider) DoSomething(ctx context.Context, input string) error {
    // Validate input
    if input == "" {
        return fmt.Errorf("input cannot be empty")
    }
    
    // Check context
    select {
    case <-ctx.Done():
        return ctx.Err()
    default:
    }
    
    // Do work...
    
    return nil
}
```

## Adding New Features

### Adding a New Provider Type

1. Define the interface in `pkg/providers/<feature>/provider.go`
2. Implement providers in subdirectories (e.g., `github/`, `custom/`)
3. Update the main `helper.Config` to include the new provider
4. Add convenience methods to the main `Helper` struct if appropriate
5. Document the provider in README.md
6. Add examples in the `examples/` directory

### Adding New Abstractions

1. Create a new package under `pkg/`
2. Define interfaces for the abstraction
3. Document expected behavior thoroughly
4. Add integration points in the main helper
5. Provide examples of implementation

## Testing

### Unit Tests

Write unit tests for your code:

```go
func TestProviderDoSomething(t *testing.T) {
    p := NewProvider(Config{})
    
    ctx := context.Background()
    err := p.DoSomething(ctx, "test")
    
    if err != nil {
        t.Errorf("expected no error, got %v", err)
    }
}
```

### Integration Tests

For integration tests that require external services, use build tags:

```go
//go:build integration

func TestGitHubProvider_Integration(t *testing.T) {
    // Test with real GitHub API
}
```

Run integration tests with:
```bash
go test -tags=integration ./...
```

## Documentation

### Code Documentation

- Document all exported types, functions, and constants
- Use complete sentences in comments
- Explain the "why" not just the "what"
- Include examples in doc comments when helpful

Example:
```go
// CheckForUpdates checks if a new version is available.
// It compares the currentVersion with the latest version from the provider.
// Returns an UpdateInfo struct indicating if an update is available.
//
// Example:
//   info, err := provider.CheckForUpdates(ctx, "1.0.0")
//   if err != nil {
//       return err
//   }
//   if info.Available {
//       fmt.Printf("New version %s available\n", info.LatestVersion.Version)
//   }
func (p *Provider) CheckForUpdates(ctx context.Context, currentVersion string) (*UpdateInfo, error) {
    // Implementation
}
```

### README Updates

When adding features, update the README.md with:
- Usage examples
- API documentation
- Configuration options

## Pull Request Process

1. **Create a clear PR description** explaining:
   - What changes you made
   - Why you made them
   - How to test them

2. **Ensure CI passes**:
   - All tests pass
   - Code is formatted
   - No linting errors

3. **Request review** from maintainers

4. **Address feedback** from reviewers

5. **Keep your PR updated** with the main branch

## Commit Message Guidelines

Follow conventional commits format:

```
<type>(<scope>): <subject>

<body>

<footer>
```

Types:
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style changes (formatting, etc.)
- `refactor`: Code refactoring
- `test`: Adding or updating tests
- `chore`: Maintenance tasks

Example:
```
feat(autoupdate): add retry logic for failed downloads

Add exponential backoff retry mechanism when download fails.
This improves reliability when network is unstable.

Fixes #123
```

## Community Guidelines

- Be respectful and inclusive
- Help others learn
- Give constructive feedback
- Focus on the code, not the person

## Questions?

If you have questions:
- Open an issue for discussion
- Check existing issues and pull requests
- Review the architecture documentation

## License

By contributing, you agree that your contributions will be licensed under the Unlicense, the same license as the project.
