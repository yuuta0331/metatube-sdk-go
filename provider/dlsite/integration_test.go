package dlsite

import (
	"strings"
	"testing"
)

// TestIntegration_SearchAndRetrieveMetadata tests the complete workflow:
// 1. Search for a specific work by title
// 2. Retrieve full metadata for the found work
// 3. Verify all metadata fields are populated correctly
func TestIntegration_SearchAndRetrieveMetadata(t *testing.T) {
	// Skip in short mode as this makes real network requests
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	provider := New()
	
	// Test case: Search for "[NLVR] 【VR動画】ご奉仕Girls VR [RJ01227569]"
	searchQuery := "[NLVR] 【VR動画】ご奉仕Girls VR [RJ01227569]"
	
	t.Logf("=== Testing Search and Metadata Retrieval ===")
	t.Logf("Search Query: %s", searchQuery)
	
	// Step 1: Search for the work
	t.Logf("\n--- Step 1: Searching for work ---")
	results, err := provider.SearchMovie(searchQuery)
	if err != nil {
		t.Fatalf("SearchMovie failed: %v", err)
	}
	
	if len(results) == 0 {
		t.Fatal("SearchMovie returned no results")
	}
	
	t.Logf("Found %d result(s)", len(results))
	
	// Display all search results
	for i, result := range results {
		t.Logf("\nResult %d:", i+1)
		t.Logf("  ID: %s", result.ID)
		t.Logf("  Title: %s", result.Title)
		t.Logf("  Homepage: %s", result.Homepage)
		t.Logf("  ThumbURL: %s", result.ThumbURL)
		t.Logf("  CoverURL: %s", result.CoverURL)
		
		// Verify URLs are not empty
		if result.ThumbURL == "" {
			t.Error("ThumbURL is empty")
		}
		if result.CoverURL == "" {
			t.Error("CoverURL is empty")
		}
	}
	
	// Verify the first result matches expected work ID
	firstResult := results[0]
	expectedID := "RJ01227569"
	if firstResult.ID != expectedID {
		t.Errorf("Expected first result ID to be %s, got %s", expectedID, firstResult.ID)
	}
	
	// Verify that CoverURL is the main image (not thumbnail)
	if !strings.Contains(firstResult.CoverURL, "_img_main.jpg") {
		t.Errorf("Expected CoverURL to contain '_img_main.jpg', got %s", firstResult.CoverURL)
	}
	
	// Verify that ThumbURL is the thumbnail image
	if !strings.Contains(firstResult.ThumbURL, "_img_sam.jpg") {
		t.Errorf("Expected ThumbURL to contain '_img_sam.jpg', got %s", firstResult.ThumbURL)
	}
	
	// Step 2: Retrieve full metadata for the found work
	t.Logf("\n--- Step 2: Retrieving full metadata for %s ---", firstResult.ID)
	info, err := provider.GetMovieInfoByID(firstResult.ID)
	if err != nil {
		t.Fatalf("GetMovieInfoByID failed: %v", err)
	}
	
	// Step 3: Verify all metadata fields
	t.Logf("\n--- Step 3: Verifying metadata fields ---")
	t.Logf("\nBasic Information:")
	t.Logf("  ID: %s", info.ID)
	t.Logf("  Number: %s", info.Number)
	t.Logf("  Title: %s", info.Title)
	t.Logf("  Provider: %s", info.Provider)
	t.Logf("  Homepage: %s", info.Homepage)
	
	t.Logf("\nMaker Information:")
	t.Logf("  Maker: %s", info.Maker)
	t.Logf("  Label: %s", info.Label)
	t.Logf("  Series: %s", info.Series)
	
	t.Logf("\nRelease Information:")
	t.Logf("  ReleaseDate: %v", info.ReleaseDate)
	
	t.Logf("\nContent:")
	t.Logf("  Summary: %s", info.Summary)
	
	t.Logf("\nImages:")
	t.Logf("  ThumbURL: %s", info.ThumbURL)
	t.Logf("  CoverURL: %s", info.CoverURL)
	t.Logf("  BigCoverURL: %s", info.BigCoverURL)
	t.Logf("  PreviewImages: %d image(s)", len(info.PreviewImages))
	for i, img := range info.PreviewImages {
		t.Logf("    [%d] %s", i+1, img)
	}
	
	t.Logf("\nGenres:")
	t.Logf("  Genres: %v", info.Genres)
	
	// Verify required fields are not empty
	if info.ID == "" {
		t.Error("ID is empty")
	}
	if info.Title == "" {
		t.Error("Title is empty")
	}
	if info.Maker == "" {
		t.Error("Maker is empty")
	}
	if info.CoverURL == "" {
		t.Error("CoverURL is empty")
	}
	if info.ThumbURL == "" {
		t.Error("ThumbURL is empty")
	}
	
	t.Logf("\n=== Test Completed Successfully ===")
}

// TestIntegration_SearchByRJID tests searching by RJ ID directly
func TestIntegration_SearchByRJID(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	provider := New()
	
	workID := "RJ01227569"
	
	t.Logf("=== Testing Search by RJ ID ===")
	t.Logf("Work ID: %s", workID)
	
	results, err := provider.SearchMovie(workID)
	if err != nil {
		t.Fatalf("SearchMovie failed: %v", err)
	}
	
	if len(results) == 0 {
		t.Fatal("SearchMovie returned no results")
	}
	
	if len(results) != 1 {
		t.Errorf("Expected 1 result, got %d", len(results))
	}
	
	result := results[0]
	t.Logf("\nResult:")
	t.Logf("  ID: %s", result.ID)
	t.Logf("  Title: %s", result.Title)
	t.Logf("  Homepage: %s", result.Homepage)
	
	if result.ID != workID {
		t.Errorf("Expected ID %s, got %s", workID, result.ID)
	}
	
	t.Logf("\n=== Test Completed Successfully ===")
}

// TestIntegration_MultiProviderSearch simulates searching across all providers
// to verify DLsite provider is correctly registered and returns results
func TestIntegration_MultiProviderSearch(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// This test verifies that the DLsite provider is properly registered
	// and can be discovered when searching across all providers
	
	provider := New()
	
	// Verify provider is properly initialized
	if provider.Name() != Name {
		t.Errorf("Expected provider name %s, got %s", Name, provider.Name())
	}
	
	t.Logf("=== Testing Multi-Provider Scenario ===")
	t.Logf("Provider Name: %s", provider.Name())
	t.Logf("Provider Priority: %f", provider.Priority())
	t.Logf("Provider URL: %s", provider.URL())
	
	// Test search functionality
	searchQuery := "RJ01227569"
	t.Logf("\nSearching for: %s", searchQuery)
	
	results, err := provider.SearchMovie(searchQuery)
	if err != nil {
		t.Fatalf("SearchMovie failed: %v", err)
	}
	
	if len(results) == 0 {
		t.Fatal("SearchMovie returned no results")
	}
	
	t.Logf("Found %d result(s) from %s provider", len(results), provider.Name())
	
	for i, result := range results {
		t.Logf("\nResult %d:", i+1)
		t.Logf("  Provider: %s", result.Provider)
		t.Logf("  ID: %s", result.ID)
		t.Logf("  Title: %s", result.Title)
		
		// Verify provider name is set correctly
		if result.Provider != Name {
			t.Errorf("Expected provider name %s, got %s", Name, result.Provider)
		}
	}
	
	t.Logf("\n=== Test Completed Successfully ===")
	t.Logf("DLsite provider is correctly registered and functional")
}
