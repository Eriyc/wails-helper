package events

import (
	"context"
)

// KeyboardModifier represents keyboard modifier keys
type KeyboardModifier uint32

const (
	// ModifierShift represents the Shift key
	ModifierShift KeyboardModifier = 1 << iota
	// ModifierControl represents the Control/Ctrl key
	ModifierControl
	// ModifierAlt represents the Alt/Option key
	ModifierAlt
	// ModifierSuper represents the Super/Command/Windows key
	ModifierSuper
)

// KeyEvent represents a keyboard event
type KeyEvent struct {
	Key       string           `json:"key"`        // The key that was pressed (e.g., "a", "Enter", "F1")
	Code      string           `json:"code"`       // The physical key code (e.g., "KeyA", "Enter")
	Modifiers KeyboardModifier `json:"modifiers"`  // Modifier keys held during the event
	Repeat    bool             `json:"repeat"`     // Whether this is a repeated key event
}

// MouseButton represents mouse buttons
type MouseButton int

const (
	// MouseButtonLeft represents the left mouse button
	MouseButtonLeft MouseButton = iota
	// MouseButtonMiddle represents the middle mouse button
	MouseButtonMiddle
	// MouseButtonRight represents the right mouse button
	MouseButtonRight
)

// MouseEvent represents a mouse event
type MouseEvent struct {
	Button    MouseButton      `json:"button"`     // The mouse button
	X         int              `json:"x"`          // X coordinate
	Y         int              `json:"y"`          // Y coordinate
	Modifiers KeyboardModifier `json:"modifiers"`  // Modifier keys held during the event
}

// EventType represents the type of event
type EventType string

const (
	// EventTypeKeyDown represents a key down event
	EventTypeKeyDown EventType = "keydown"
	// EventTypeKeyUp represents a key up event
	EventTypeKeyUp EventType = "keyup"
	// EventTypeMouseDown represents a mouse down event
	EventTypeMouseDown EventType = "mousedown"
	// EventTypeMouseUp represents a mouse up event
	EventTypeMouseUp EventType = "mouseup"
	// EventTypeMouseMove represents a mouse move event
	EventTypeMouseMove EventType = "mousemove"
	// EventTypeMouseWheel represents a mouse wheel event
	EventTypeMouseWheel EventType = "mousewheel"
)

// Handler represents an event handler function
type Handler func(event interface{}) error

// Shortcut represents a keyboard shortcut
type Shortcut struct {
	Key       string           `json:"key"`
	Modifiers KeyboardModifier `json:"modifiers"`
	Handler   Handler          `json:"-"`
}

// Manager defines the interface for event management
type Manager interface {
	// On registers an event handler
	On(eventType EventType, handler Handler) error

	// Off removes an event handler
	Off(eventType EventType, handler Handler) error

	// Emit emits an event
	Emit(ctx context.Context, eventType EventType, event interface{}) error

	// RegisterShortcut registers a keyboard shortcut
	RegisterShortcut(ctx context.Context, shortcut Shortcut) error

	// UnregisterShortcut unregisters a keyboard shortcut
	UnregisterShortcut(ctx context.Context, key string, modifiers KeyboardModifier) error

	// RegisterGlobalShortcut registers a global keyboard shortcut (system-wide)
	RegisterGlobalShortcut(ctx context.Context, shortcut Shortcut) error

	// UnregisterGlobalShortcut unregisters a global keyboard shortcut
	UnregisterGlobalShortcut(ctx context.Context, key string, modifiers KeyboardModifier) error
}

// CustomEvent represents a custom application event
type CustomEvent struct {
	Name string                 `json:"name"`
	Data map[string]interface{} `json:"data,omitempty"`
}

// Emitter defines the interface for emitting custom events
type Emitter interface {
	// Emit emits a custom event
	Emit(ctx context.Context, event CustomEvent) error

	// EmitSync emits a custom event synchronously
	EmitSync(ctx context.Context, event CustomEvent) error

	// On registers a handler for custom events
	On(eventName string, handler func(event CustomEvent) error) error

	// Once registers a one-time handler for custom events
	Once(eventName string, handler func(event CustomEvent) error) error

	// Off removes a handler for custom events
	Off(eventName string, handler func(event CustomEvent) error) error
}
