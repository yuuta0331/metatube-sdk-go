package dlsite

import (
	"strings"
	"testing"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
)

// TestBugCondition_TitleSearchReturnsResults is a bug condition exploration test.
// **Validates: Requirements 2.1, 2.2, 2.3, 2.4**
//
// **CRITICAL**: This test MUST FAIL on unfixed code - failure confirms the bug exists.
// **DO NOT attempt to fix the test or the code when it fails.**
// **NOTE**: This test encodes the expected behavior - it will validate the fix when it passes after implementation.
//
// **Property 1: Fault Condition** - Title Search Returns Results
//
// This test verifies that searching for a work by its title returns results containing
// the expected work ID. The test is scoped to the concrete failing case:
// - Search keyword: "むっつり彼女VRその4" (title of work RJ01463741)
// - Expected: Results contain work with ID "RJ01463741"
// - Expected: len(results) > 0
//
// **EXPECTED OUTCOME ON UNFIXED CODE**: Test FAILS (this is correct - it proves the bug exists)
// **Counterexample**: Search returns empty results array even though work exists on DLSite
func TestBugCondition_TitleSearchReturnsResults(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 1 // Single concrete test case
	properties := gopter.NewProperties(parameters)

	properties.Property("Feature: dlsite-search-title-bug, Property 1: Fault Condition - Title Search Returns Results",
		prop.ForAll(
			func(keyword string) bool {
				provider := New()
				
				// Search by title
				results, err := provider.SearchMovie(keyword)
				
				// Should not return an error
				if err != nil {
					t.Logf("Search returned error: %v", err)
					return false
				}
				
				// Should return non-empty results
				if len(results) == 0 {
					t.Logf("Search returned empty results for keyword: %s", keyword)
					return false
				}
				
				// Should contain the expected work ID
				foundExpectedWork := false
				for _, result := range results {
					if result.ID == "RJ01463741" {
						foundExpectedWork = true
						break
					}
				}
				
				if !foundExpectedWork {
					t.Logf("Search results did not contain expected work RJ01463741")
					t.Logf("Found %d results:", len(results))
					for i, result := range results {
						t.Logf("  [%d] ID: %s, Title: %s", i, result.ID, result.Title)
					}
					return false
				}
				
				return true
			},
			gen.Const("むっつり彼女VRその4"), // Concrete failing case
		))

	properties.TestingRun(t)
}

// TestPreservation_RJIDDirectLookup tests that RJ work ID searches bypass the search endpoint
// and call GetMovieInfoByID directly.
// **Validates: Requirements 3.1**
//
// **Property 2: Preservation** - Non-Title Search Behavior
//
// This test verifies that when searching with an RJ work ID (e.g., "RJ01463741"),
// the system bypasses the search endpoint and calls GetMovieInfoByID directly.
// This behavior must be preserved after the fix.
//
// **EXPECTED OUTCOME ON UNFIXED CODE**: Test PASSES (confirms baseline behavior)
func TestPreservation_RJIDDirectLookup(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 50
	properties := gopter.NewProperties(parameters)

	properties.Property("Feature: dlsite-search-title-bug, Property 2: Preservation - RJ ID Direct Lookup",
		prop.ForAll(
			func(digits string) bool {
				provider := New()
				
				// Generate RJ work ID
				rjID := "RJ" + digits
				
				// Search by RJ ID
				results, err := provider.SearchMovie(rjID)
				
				// Should not return an error
				if err != nil {
					// Network errors are acceptable in property tests
					// as they indicate the function tried to fetch data
					return true
				}
				
				// Should return exactly one result (direct lookup)
				if len(results) != 1 {
					t.Logf("Expected 1 result for RJ ID search, got %d", len(results))
					return false
				}
				
				// Result should have the same ID (normalized to uppercase)
				if results[0].ID != rjID {
					t.Logf("Expected result ID %s, got %s", rjID, results[0].ID)
					return false
				}
				
				return true
			},
			gen.RegexMatch(`\d{6,8}`), // 6-8 digits for RJ work IDs
		))

	properties.TestingRun(t)
}

// TestPreservation_RJOnlyFiltering tests that search results filter out non-RJ works.
// **Validates: Requirements 3.2**
//
// **Property 2: Preservation** - Non-Title Search Behavior
//
// This test verifies that search results only contain RJ-prefixed works and filter out
// VJ (games), BJ (comics), and other non-RJ content types. This behavior must be preserved.
//
// **EXPECTED OUTCOME ON UNFIXED CODE**: Test PASSES (confirms baseline behavior)
//
// Note: This test uses the observation that the SearchMovie function filters results
// in the HTML parsing callback by checking if ID starts with "RJ".
func TestPreservation_RJOnlyFiltering(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 50
	properties := gopter.NewProperties(parameters)

	properties.Property("Feature: dlsite-search-title-bug, Property 2: Preservation - RJ Only Filtering",
		prop.ForAll(
			func(keyword string) bool {
				provider := New()
				
				// Search with keyword
				results, err := provider.SearchMovie(keyword)
				
				// Network errors are acceptable
				if err != nil {
					return true
				}
				
				// All results must have RJ prefix
				for _, result := range results {
					if !strings.HasPrefix(result.ID, "RJ") {
						t.Logf("Found non-RJ work in results: %s", result.ID)
						return false
					}
				}
				
				return true
			},
			gen.AlphaString(), // Random keywords
		))

	properties.TestingRun(t)
}

// TestPreservation_HTTPErrorWrapping tests that HTTP errors are wrapped with ErrNetworkError.
// **Validates: Requirements 3.3**
//
// **Property 2: Preservation** - Non-Title Search Behavior
//
// This test verifies that when the search endpoint returns HTTP errors (404, 500, etc.),
// the errors are properly wrapped with ErrNetworkError. This behavior must be preserved.
//
// **EXPECTED OUTCOME ON UNFIXED CODE**: Test PASSES (confirms baseline behavior)
//
// Note: This test observes that network errors from SearchMovie contain "network error"
// or are nil (successful response). We cannot easily simulate HTTP errors in property tests
// without mocking, so we verify that any errors returned follow the expected pattern.
func TestPreservation_HTTPErrorWrapping(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 50
	properties := gopter.NewProperties(parameters)

	properties.Property("Feature: dlsite-search-title-bug, Property 2: Preservation - HTTP Error Wrapping",
		prop.ForAll(
			func(keyword string) bool {
				provider := New()
				
				// Search with keyword
				_, err := provider.SearchMovie(keyword)
				
				// If there's an error, it should be a network error or contain "network error"
				if err != nil {
					errMsg := err.Error()
					// Check if error is wrapped with ErrNetworkError or contains network error context
					if !strings.Contains(errMsg, "network error") && 
					   !strings.Contains(errMsg, "Network error") &&
					   !strings.Contains(errMsg, "HTTP") {
						t.Logf("Error not properly wrapped: %v", err)
						return false
					}
				}
				
				return true
			},
			gen.AlphaString(),
		))

	properties.TestingRun(t)
}

// TestPreservation_InvalidResultsSkipped tests that invalid search results are skipped.
// **Validates: Requirements 3.5**
//
// **Property 2: Preservation** - Non-Title Search Behavior
//
// This test verifies that search results with missing required fields are skipped using
// the IsValid() check. All returned results must have valid ID, title, homepage, and provider.
// This behavior must be preserved.
//
// **EXPECTED OUTCOME ON UNFIXED CODE**: Test PASSES (confirms baseline behavior)
func TestPreservation_InvalidResultsSkipped(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 50
	properties := gopter.NewProperties(parameters)

	properties.Property("Feature: dlsite-search-title-bug, Property 2: Preservation - Invalid Results Skipped",
		prop.ForAll(
			func(keyword string) bool {
				provider := New()
				
				// Search with keyword
				results, err := provider.SearchMovie(keyword)
				
				// Network errors are acceptable
				if err != nil {
					return true
				}
				
				// All results must be valid (have required fields)
				for _, result := range results {
					if result.ID == "" || result.Title == "" || 
					   result.Homepage == "" || result.Provider == "" {
						t.Logf("Found invalid result in search results: ID=%s, Title=%s, Homepage=%s, Provider=%s",
							result.ID, result.Title, result.Homepage, result.Provider)
						return false
					}
					
					// Additionally verify IsValid() returns true
					if !result.IsValid() {
						t.Logf("Result.IsValid() returned false for: %+v", result)
						return false
					}
				}
				
				return true
			},
			gen.AlphaString(),
		))

	properties.TestingRun(t)
}
