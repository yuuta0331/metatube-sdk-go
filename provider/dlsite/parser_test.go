package dlsite

import (
	"strings"
	"testing"
)

func TestNormalizeMovieID_Examples(t *testing.T) {
	dlsite := New()
	
	tests := []struct {
		input    string
		expected string
		desc     string
	}{
		{"RJ123456", "RJ123456", "valid RJ ID"},
		{"rj123456", "RJ123456", "lowercase RJ ID"},
		{"RJ01234567", "RJ01234567", "7-digit RJ ID"},
		{"RJ12345678", "RJ12345678", "8-digit RJ ID"},
		{"https://www.dlsite.com/maniax/work/=/product_id/RJ123456.html", "RJ123456", "RJ ID from URL"},
		{"VJ123456", "", "VJ ID (game) - unsupported"},
		{"BJ123456", "", "BJ ID (comic) - unsupported"},
		{"RE123456", "", "RE ID - unsupported"},
		{"invalid", "", "invalid string"},
		{"", "", "empty string"},
		{"RJ12345", "", "too few digits"},
		{"RJ123456789", "RJ12345678", "9 digits - extracts first 8"},
	}
	
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			result := dlsite.NormalizeMovieID(tt.input)
			if result != tt.expected {
				t.Errorf("NormalizeMovieID(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestParseMovieIDFromURL(t *testing.T) {
	dlsite := New()
	
	tests := []struct {
		url      string
		expected string
		wantErr  bool
		desc     string
	}{
		{
			"https://www.dlsite.com/maniax/work/=/product_id/RJ123456.html",
			"RJ123456",
			false,
			"valid maniax URL",
		},
		{
			"https://www.dlsite.com/home/work/=/product_id/RJ789012.html",
			"RJ789012",
			false,
			"valid home URL",
		},
		{
			"https://play.dlsite.com/csr/=/product_id/RJ456789",
			"RJ456789",
			false,
			"valid play URL",
		},
		{
			"https://example.com/RJ123456",
			"",
			true,
			"non-DLSite domain",
		},
		{
			"https://www.dlsite.com/page/without/id",
			"",
			true,
			"URL without work ID",
		},
		{
			"invalid-url",
			"",
			true,
			"invalid URL format",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			result, err := dlsite.ParseMovieIDFromURL(tt.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseMovieIDFromURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if result != tt.expected {
				t.Errorf("ParseMovieIDFromURL() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestNormalizeMovieKeyword(t *testing.T) {
	dlsite := New()
	
	tests := []struct {
		keyword  string
		expected string
		desc     string
	}{
		{"RJ123456", "RJ123456", "work ID as keyword"},
		{"rj123456", "RJ123456", "lowercase work ID"},
		{"some keyword", "some keyword", "regular keyword"},
		{"keyword with RJ123456 in it", "RJ123456", "keyword containing work ID"},
	}
	
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			result := dlsite.NormalizeMovieKeyword(tt.keyword)
			if result != tt.expected {
				t.Errorf("NormalizeMovieKeyword(%q) = %q, want %q", tt.keyword, result, tt.expected)
			}
		})
	}
}

func TestGetMovieInfoByID_UnsupportedContentType(t *testing.T) {
	dlsite := New()
	
	unsupportedIDs := []string{
		"VJ123456", // Game
		"BJ123456", // Comic
		"RE123456", // Other
	}
	
	for _, id := range unsupportedIDs {
		t.Run(id, func(t *testing.T) {
			_, err := dlsite.GetMovieInfoByID(id)
			if err == nil {
				t.Errorf("Expected error for unsupported content type %s, got nil", id)
			}
			
			errMsg := err.Error()
			if !strings.Contains(errMsg, "unsupported") && !strings.Contains(errMsg, "invalid") {
				t.Errorf("Error message should indicate unsupported content type, got: %s", errMsg)
			}
		})
	}
}
