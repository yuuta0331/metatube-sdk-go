package dlsite

import (
	"encoding/json"
	"testing"
)

// TestSearchResultJSONOutput tests the JSON output of search results
func TestSearchResultJSONOutput(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping JSON output test in short mode")
	}

	provider := New()
	
	searchQuery := "RJ01227569"
	t.Logf("=== Testing JSON Output for Search Results ===")
	t.Logf("Search Query: %s", searchQuery)
	
	results, err := provider.SearchMovie(searchQuery)
	if err != nil {
		t.Fatalf("SearchMovie failed: %v", err)
	}
	
	if len(results) == 0 {
		t.Fatal("SearchMovie returned no results")
	}
	
	// Convert to JSON
	jsonData, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal to JSON: %v", err)
	}
	
	t.Logf("\n=== JSON Output ===\n%s", string(jsonData))
	
	// Verify JSON structure
	var decoded []map[string]interface{}
	if err := json.Unmarshal(jsonData, &decoded); err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}
	
	if len(decoded) == 0 {
		t.Fatal("Decoded JSON is empty")
	}
	
	firstResult := decoded[0]
	
	// Check required fields
	requiredFields := []string{"id", "number", "title", "provider", "homepage", "thumb_url", "cover_url"}
	for _, field := range requiredFields {
		if _, ok := firstResult[field]; !ok {
			t.Errorf("Missing required field: %s", field)
		} else {
			t.Logf("✓ Field '%s': %v", field, firstResult[field])
		}
	}
	
	// Verify image URLs are not empty
	if thumbURL, ok := firstResult["thumb_url"].(string); !ok || thumbURL == "" {
		t.Error("thumb_url is empty or not a string")
	}
	
	if coverURL, ok := firstResult["cover_url"].(string); !ok || coverURL == "" {
		t.Error("cover_url is empty or not a string")
	}
}
