package window

import (
	"context"
)

// Position represents the position of a window
type Position struct {
	X int `json:"x"`
	Y int `json:"y"`
}

// Size represents the size of a window
type Size struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

// Options represents window creation options
type Options struct {
	Title           string   `json:"title"`
	Width           int      `json:"width"`
	Height          int      `json:"height"`
	MinWidth        int      `json:"min_width,omitempty"`
	MinHeight       int      `json:"min_height,omitempty"`
	MaxWidth        int      `json:"max_width,omitempty"`
	MaxHeight       int      `json:"max_height,omitempty"`
	Resizable       bool     `json:"resizable"`
	Frameless       bool     `json:"frameless"`
	AlwaysOnTop     bool     `json:"always_on_top"`
	Hidden          bool     `json:"hidden"`
	Fullscreen      bool     `json:"fullscreen"`
	StartState      string   `json:"start_state,omitempty"` // "normal", "maximized", "minimized"
	URL             string   `json:"url,omitempty"`
	HTML            string   `json:"html,omitempty"`
	BackgroundColor string   `json:"background_color,omitempty"`
}

// Manager defines the interface for window management
type Manager interface {
	// Create creates a new window with the given options
	Create(ctx context.Context, opts Options) (Window, error)

	// Get retrieves a window by its ID
	Get(ctx context.Context, id string) (Window, error)

	// GetAll retrieves all windows
	GetAll(ctx context.Context) ([]Window, error)

	// GetMain retrieves the main window
	GetMain(ctx context.Context) (Window, error)

	// Close closes a window by its ID
	Close(ctx context.Context, id string) error
}

// Window represents a window instance
type Window interface {
	// ID returns the unique identifier of the window
	ID() string

	// Show shows the window
	Show(ctx context.Context) error

	// Hide hides the window
	Hide(ctx context.Context) error

	// Close closes the window
	Close(ctx context.Context) error

	// Minimize minimizes the window
	Minimize(ctx context.Context) error

	// Maximize maximizes the window
	Maximize(ctx context.Context) error

	// Restore restores the window to its normal state
	Restore(ctx context.Context) error

	// Fullscreen makes the window fullscreen
	Fullscreen(ctx context.Context) error

	// UnFullscreen exits fullscreen mode
	UnFullscreen(ctx context.Context) error

	// SetTitle sets the window title
	SetTitle(ctx context.Context, title string) error

	// GetTitle gets the window title
	GetTitle(ctx context.Context) (string, error)

	// SetSize sets the window size
	SetSize(ctx context.Context, size Size) error

	// GetSize gets the window size
	GetSize(ctx context.Context) (Size, error)

	// SetPosition sets the window position
	SetPosition(ctx context.Context, pos Position) error

	// GetPosition gets the window position
	GetPosition(ctx context.Context) (Position, error)

	// SetAlwaysOnTop sets whether the window should always be on top
	SetAlwaysOnTop(ctx context.Context, onTop bool) error

	// Center centers the window on the screen
	Center(ctx context.Context) error

	// SetResizable sets whether the window is resizable
	SetResizable(ctx context.Context, resizable bool) error

	// SetBackgroundColor sets the background color of the window
	SetBackgroundColor(ctx context.Context, color string) error

	// ExecJS executes JavaScript in the window
	ExecJS(ctx context.Context, js string) error

	// On registers an event handler for the window
	On(event string, handler func(data interface{})) error

	// Emit emits an event to the window
	Emit(ctx context.Context, event string, data interface{}) error
}
