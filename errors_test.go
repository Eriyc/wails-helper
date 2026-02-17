package helper

import (
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestWrap(t *testing.T) {
	cause := errors.New("network timeout")
	err := Wrap(ErrCodeUnavailable, "unable to fetch release", cause)

	if err == nil {
		t.Fatal("expected non-nil AppError")
	}

	if err.Code != ErrCodeUnavailable {
		t.Fatalf("expected code %q, got %q", ErrCodeUnavailable, err.Code)
	}

	if err.Message != "unable to fetch release" {
		t.Fatalf("unexpected message: %q", err.Message)
	}

	if !errors.Is(err, cause) {
		t.Fatal("expected wrapped cause to be discoverable with errors.Is")
	}
}

func TestWrapDefaultsToUnknownCode(t *testing.T) {
	err := Wrap("", "fallback", nil)
	if err.Code != ErrCodeUnknown {
		t.Fatalf("expected default code %q, got %q", ErrCodeUnknown, err.Code)
	}
}

func TestIsCode(t *testing.T) {
	base := NewInvalidInput("bad payload", errors.New("missing field"))
	wrapped := fmt.Errorf("request failed: %w", base)

	if !IsCode(base, ErrCodeInvalidInput) {
		t.Fatal("expected base AppError to match invalid input code")
	}

	if !IsCode(wrapped, ErrCodeInvalidInput) {
		t.Fatal("expected wrapped AppError to match invalid input code")
	}

	if IsCode(wrapped, ErrCodeUnavailable) {
		t.Fatal("did not expect unavailable code to match")
	}

	if IsCode(nil, ErrCodeInvalidInput) {
		t.Fatal("nil error must not match any code")
	}
}

func TestPanicToError(t *testing.T) {
	err := PanicToError("boom", map[string]string{"source": "test"})
	if err == nil {
		t.Fatal("expected non-nil error")
	}

	if err.Code != ErrCodeUnknown {
		t.Fatalf("expected unknown code, got %q", err.Code)
	}

	if err.Message != "panic recovered" {
		t.Fatalf("unexpected message: %q", err.Message)
	}

	if err.Metadata["source"] != "test" {
		t.Fatal("expected caller metadata to be preserved")
	}

	if err.Metadata["panic.type"] != "string" {
		t.Fatalf("unexpected panic type: %q", err.Metadata["panic.type"])
	}

	if err.Metadata["panic.value"] != "boom" {
		t.Fatalf("unexpected panic value: %q", err.Metadata["panic.value"])
	}

	if strings.TrimSpace(err.Metadata["panic.stack"]) == "" {
		t.Fatal("expected non-empty panic stack metadata")
	}
}

func TestRunWithRecovery(t *testing.T) {
	err := RunWithRecovery(func() {
		panic("callback panic")
	}, map[string]string{"source": "callback"})

	if err == nil {
		t.Fatal("expected panic to be converted to error")
	}

	if !IsCode(err, ErrCodeUnknown) {
		t.Fatal("expected recovered panic to have unknown code")
	}

	var appErr *AppError
	if !errors.As(err, &appErr) {
		t.Fatal("expected AppError from RunWithRecovery")
	}

	if appErr.Metadata["panic.context"] != "callback" {
		t.Fatalf("expected callback context, got %q", appErr.Metadata["panic.context"])
	}

	if appErr.Metadata["source"] != "callback" {
		t.Fatal("expected metadata to include source")
	}
}

func TestGoWithRecovery(t *testing.T) {
	errCh := make(chan error, 1)

	GoWithRecovery(func() {
		panic("goroutine panic")
	}, func(err error) {
		errCh <- err
	}, map[string]string{"source": "goroutine"})

	select {
	case err := <-errCh:
		if err == nil {
			t.Fatal("expected panic error")
		}

		var appErr *AppError
		if !errors.As(err, &appErr) {
			t.Fatal("expected AppError from goroutine recovery")
		}

		if appErr.Metadata["panic.context"] != "goroutine" {
			t.Fatalf("expected goroutine context, got %q", appErr.Metadata["panic.context"])
		}

		if appErr.Metadata["source"] != "goroutine" {
			t.Fatal("expected metadata to include source")
		}
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for recovered goroutine panic")
	}
}
