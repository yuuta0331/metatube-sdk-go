package dlsite

import (
	"testing"

	"golang.org/x/text/language"
)

func TestProviderInterface(t *testing.T) {
	dlsite := New()
	
	// Test provider name
	if dlsite.Name() != "dlsite" {
		t.Errorf("expected provider name 'dlsite', got '%s'", dlsite.Name())
	}
	
	// Test base URL
	expectedURL := "https://www.dlsite.com"
	if dlsite.URL().String() != expectedURL {
		t.Errorf("expected URL '%s', got '%s'", expectedURL, dlsite.URL().String())
	}
	
	// Test language
	if dlsite.Language() != language.Japanese {
		t.Errorf("expected Japanese language, got %v", dlsite.Language())
	}
}

func TestProviderRegistration(t *testing.T) {
	// Test that provider is registered in factory
	// This will be implemented after checking the factory interface
}
