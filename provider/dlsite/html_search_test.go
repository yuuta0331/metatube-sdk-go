package dlsite

import (
	"testing"
)

// TestHTMLSearchImageURLs tests that HTML search returns correct image URLs
func TestHTMLSearchImageURLs(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping HTML search test in short mode")
	}

	provider := New()
	
	// Use a keyword that will trigger HTML search (not an RJ ID)
	searchQuery := "VR"
	t.Logf("=== Testing HTML Search Image URLs ===")
	t.Logf("Search Query: %s", searchQuery)
	
	results, err := provider.SearchMovie(searchQuery)
	if err != nil {
		t.Fatalf("SearchMovie failed: %v", err)
	}
	
	if len(results) == 0 {
		t.Fatal("SearchMovie returned no results")
	}
	
	t.Logf("Found %d result(s)", len(results))
	
	// Check first few results
	maxResults := 3
	if len(results) < maxResults {
		maxResults = len(results)
	}
	
	for i := 0; i < maxResults; i++ {
		result := results[i]
		t.Logf("\n=== Result %d ===", i+1)
		t.Logf("ID: %s", result.ID)
		t.Logf("Title: %s", result.Title)
		t.Logf("ThumbURL: %s", result.ThumbURL)
		t.Logf("CoverURL: %s", result.CoverURL)
		
		// Verify URLs have valid protocol
		if result.ThumbURL != "" {
			if !hasValidProtocol(result.ThumbURL) {
				t.Errorf("Result %d: ThumbURL does not have valid protocol: %s", i+1, result.ThumbURL)
			}
		} else {
			t.Errorf("Result %d: ThumbURL is empty", i+1)
		}
		
		if result.CoverURL != "" {
			if !hasValidProtocol(result.CoverURL) {
				t.Errorf("Result %d: CoverURL does not have valid protocol: %s", i+1, result.CoverURL)
			}
		} else {
			t.Errorf("Result %d: CoverURL is empty", i+1)
		}
	}
}
