package dlsite

import (
	"testing"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
)

// Property 1: Fault Condition - CSS Selector Mismatch Detection
// **Validates: Requirements 1.1, 1.2**
//
// CRITICAL: This test MUST FAIL on unfixed code - failure confirms the bug exists
// DO NOT attempt to fix the test or the code when it fails
// NOTE: This test encodes the expected behavior - it will validate the fix when it passes after implementation
// GOAL: Surface counterexamples that demonstrate the bug exists
//
// This property tests that the parser successfully extracts metadata from HTML pages
// with the new Vue.js CSS selectors. On unfixed code, this test will FAIL because
// the current selectors (e.g., #work_name, .maker_name) don't match the new HTML structure.
//
// The test uses testdata files that contain the new Vue.js HTML structure.
// When run on unfixed code, it will demonstrate that:
// - Title field is empty (selector #work_name doesn't match)
// - Maker field is empty (selector span[class*='maker_name'] doesn't match)
// - CoverURL field is empty (selector .product-slider-item img doesn't match)
// - Validation fails with "missing required fields"
func TestProperty1_FaultCondition_CSSSelectorMismatchDetection(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 10 // Small number since we're testing with specific work IDs
	properties := gopter.NewProperties(parameters)

	// Generator for known work IDs that exist on DLsite with Vue.js structure
	// These are real work IDs that we know exist and have the new HTML structure
	genKnownWorkID := func() gopter.Gen {
		return gen.OneConstOf(
			"RJ01463741", // Known work: むっつり彼女VRその4 by NLVR
			// Add more known work IDs here as needed for testing
		)
	}

	properties.Property("Feature: dlsite-metadata-retrieval-bug, Property 1: Fault Condition - Metadata extraction from Vue.js pages succeeds",
		prop.ForAll(
			func(workID string) bool {
				provider := New()

				// Call GetMovieInfoByID with a known work ID
				info, err := provider.GetMovieInfoByID(workID)

				// On unfixed code, this will return an error due to validation failure
				// On fixed code, this should succeed
				if err != nil {
					// Log the error for debugging
					t.Logf("GetMovieInfoByID(%s) failed: %v", workID, err)
					return false
				}

				// Verify that all required fields are populated
				// These assertions match the Expected Behavior Properties from design:
				// - Parser successfully extracts all metadata fields
				// - Extracted values match expected values from HTML
				// - No fields are empty or missing

				if info.Title == "" {
					t.Logf("Title is empty for work ID %s", workID)
					return false
				}

				if info.Maker == "" {
					t.Logf("Maker is empty for work ID %s", workID)
					return false
				}

				if info.CoverURL == "" {
					t.Logf("CoverURL is empty for work ID %s", workID)
					return false
				}

				// Verify that the info passes validation
				if !info.IsValid() {
					t.Logf("MovieInfo validation failed for work ID %s", workID)
					return false
				}

				// Verify that ID and Number are set correctly
				if info.ID != workID {
					t.Logf("ID mismatch: expected %s, got %s", workID, info.ID)
					return false
				}

				if info.Number != workID {
					t.Logf("Number mismatch: expected %s, got %s", workID, info.Number)
					return false
				}

				// Verify that Provider and Homepage are set
				if info.Provider != "dlsite" {
					t.Logf("Provider mismatch: expected 'dlsite', got '%s'", info.Provider)
					return false
				}

				if info.Homepage == "" {
					t.Logf("Homepage is empty for work ID %s", workID)
					return false
				}

				// All checks passed - metadata extraction succeeded
				return true
			},
			genKnownWorkID(),
		))

	properties.TestingRun(t)
}
