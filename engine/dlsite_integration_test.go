package engine

import (
	"testing"

	"github.com/metatube-community/metatube-sdk-go/provider"
	"github.com/metatube-community/metatube-sdk-go/provider/dlsite"
)

// TestDLsiteIntegration_AllProvidersSearch tests that DLsite provider
// is correctly registered in the engine and can be discovered when
// searching across all registered providers
func TestDLsiteIntegration_AllProvidersSearch(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	t.Logf("=== Testing DLsite Provider Registration in Engine ===")

	// Verify DLsite is registered
	dlsiteFound := false
	provider.RangeMovieFactory(func(name string, factory provider.MovieFactory) bool {
		if name == dlsite.Name {
			dlsiteFound = true
			t.Logf("✓ DLsite provider found in registry: %s", name)
			return false // Stop iteration
		}
		return true // Continue iteration
	})

	if !dlsiteFound {
		t.Fatal("DLsite provider is not registered in the engine")
	}

	// List all registered movie providers
	t.Logf("\n--- All Registered Movie Providers ---")
	providerCount := 0
	provider.RangeMovieFactory(func(name string, factory provider.MovieFactory) bool {
		providerCount++
		t.Logf("  [%d] %s", providerCount, name)
		return true
	})
	t.Logf("Total: %d provider(s)", providerCount)

	// Test searching with DLsite provider
	t.Logf("\n--- Testing Search with DLsite Provider ---")
	searchQuery := "RJ01227569"
	t.Logf("Search Query: %s", searchQuery)

	// Get DLsite provider instance
	var dlsiteProvider provider.MovieProvider
	provider.RangeMovieFactory(func(name string, factory provider.MovieFactory) bool {
		if name == dlsite.Name {
			dlsiteProvider = factory()
			return false
		}
		return true
	})

	if dlsiteProvider == nil {
		t.Fatal("Failed to create DLsite provider instance")
	}

	// Check if provider implements MovieSearcher
	searcher, ok := dlsiteProvider.(provider.MovieSearcher)
	if !ok {
		t.Fatal("DLsite provider does not implement MovieSearcher interface")
	}

	// Perform search
	results, err := searcher.SearchMovie(searchQuery)
	if err != nil {
		t.Fatalf("SearchMovie failed: %v", err)
	}

	if len(results) == 0 {
		t.Fatal("SearchMovie returned no results")
	}

	t.Logf("Found %d result(s)", len(results))

	for i, result := range results {
		t.Logf("\nResult %d:", i+1)
		t.Logf("  Provider: %s", result.Provider)
		t.Logf("  ID: %s", result.ID)
		t.Logf("  Title: %s", result.Title)
		t.Logf("  Homepage: %s", result.Homepage)

		// Verify provider name
		if result.Provider != dlsite.Name {
			t.Errorf("Expected provider name %s, got %s", dlsite.Name, result.Provider)
		}
	}

	t.Logf("\n=== Test Completed Successfully ===")
	t.Logf("DLsite provider is correctly registered and functional in the engine")
}

// TestDLsiteIntegration_MultiProviderSearch simulates searching across
// multiple providers to verify DLsite returns results for RJ-prefixed IDs
func TestDLsiteIntegration_MultiProviderSearch(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	t.Logf("=== Testing Multi-Provider Search Scenario ===")

	searchQuery := "[NLVR] 【VR動画】ご奉仕Girls VR [RJ01227569]"
	t.Logf("Search Query: %s", searchQuery)

	// Simulate searching across all providers
	t.Logf("\n--- Searching Across All Providers ---")

	foundResults := false
	var dlsiteResults int

	provider.RangeMovieFactory(func(name string, factory provider.MovieFactory) bool {
		p := factory()

		// Check if provider implements MovieSearcher
		searcher, ok := p.(provider.MovieSearcher)
		if !ok {
			t.Logf("  [%s] Does not implement MovieSearcher, skipping", name)
			return true
		}

		// Perform search
		results, err := searcher.SearchMovie(searchQuery)
		if err != nil {
			t.Logf("  [%s] Search failed: %v", name, err)
			return true
		}

		if len(results) > 0 {
			foundResults = true
			t.Logf("  [%s] Found %d result(s)", name, len(results))

			if name == dlsite.Name {
				dlsiteResults = len(results)
				// Display DLsite results
				for i, result := range results {
					t.Logf("    Result %d:", i+1)
					t.Logf("      ID: %s", result.ID)
					t.Logf("      Title: %s", result.Title)
				}
			}
		} else {
			t.Logf("  [%s] No results", name)
		}

		return true
	})

	if !foundResults {
		t.Error("No provider returned results for the search query")
	}

	if dlsiteResults == 0 {
		t.Error("DLsite provider did not return any results")
	} else {
		t.Logf("\n✓ DLsite provider successfully returned %d result(s)", dlsiteResults)
	}

	t.Logf("\n=== Test Completed Successfully ===")
}
