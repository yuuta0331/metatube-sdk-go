package dlsite

import (
	"testing"
	"time"
)

// TestWithTimeout verifies that WithTimeout option correctly sets the timeout.
func TestWithTimeout(t *testing.T) {
	timeout := 60 * time.Second
	provider := New(WithTimeout(timeout))

	if provider.options.Timeout != timeout {
		t.Errorf("expected timeout %v, got %v", timeout, provider.options.Timeout)
	}
}

// TestWithRequestTimeout verifies that WithRequestTimeout option correctly sets the request timeout.
func TestWithRequestTimeout(t *testing.T) {
	timeout := 45 * time.Second
	provider := New(WithRequestTimeout(timeout))

	if provider.options.RequestTimeout != timeout {
		t.Errorf("expected request timeout %v, got %v", timeout, provider.options.RequestTimeout)
	}
}

// TestDefaultOptions verifies that default options are set correctly.
func TestDefaultOptions(t *testing.T) {
	provider := New()

	if provider.options.Timeout != 30*time.Second {
		t.Errorf("expected default timeout 30s, got %v", provider.options.Timeout)
	}

	if provider.options.RequestTimeout != 30*time.Second {
		t.Errorf("expected default request timeout 30s, got %v", provider.options.RequestTimeout)
	}
}

// TestMultipleOptions verifies that multiple options can be applied together.
func TestMultipleOptions(t *testing.T) {
	timeout := 60 * time.Second
	requestTimeout := 45 * time.Second

	provider := New(
		WithTimeout(timeout),
		WithRequestTimeout(requestTimeout),
	)

	if provider.options.Timeout != timeout {
		t.Errorf("expected timeout %v, got %v", timeout, provider.options.Timeout)
	}

	if provider.options.RequestTimeout != requestTimeout {
		t.Errorf("expected request timeout %v, got %v", requestTimeout, provider.options.RequestTimeout)
	}
}

// TestBackwardCompatibility verifies that New() without options still works.
func TestBackwardCompatibility(t *testing.T) {
	// This should not panic and should create a valid provider
	provider := New()

	if provider == nil {
		t.Fatal("expected non-nil provider")
	}

	if provider.options == nil {
		t.Fatal("expected non-nil options")
	}

	// Verify provider methods work
	name := provider.Name()
	if name != "dlsite" {
		t.Errorf("expected name 'dlsite', got %s", name)
	}
}
