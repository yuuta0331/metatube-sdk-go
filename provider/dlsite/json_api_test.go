package dlsite

import (
	"testing"
)

// TestJSONAPIImageURLs tests that JSON API returns correct image URLs
func TestJSONAPIImageURLs(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping JSON API test in short mode")
	}

	provider := New()
	
	workID := "RJ01227569"
	t.Logf("=== Testing JSON API Image URLs ===")
	t.Logf("Work ID: %s", workID)
	
	info, err := provider.GetMovieInfoByID(workID)
	if err != nil {
		t.Fatalf("GetMovieInfoByID failed: %v", err)
	}
	
	t.Logf("\n=== MovieInfo Image URLs ===")
	t.Logf("ThumbURL: %s", info.ThumbURL)
	t.Logf("BigThumbURL: %s", info.BigThumbURL)
	t.Logf("CoverURL: %s", info.CoverURL)
	t.Logf("BigCoverURL: %s", info.BigCoverURL)
	
	// Verify URLs start with https://
	if info.ThumbURL != "" && !hasValidProtocol(info.ThumbURL) {
		t.Errorf("ThumbURL does not have valid protocol: %s", info.ThumbURL)
	}
	if info.BigThumbURL != "" && !hasValidProtocol(info.BigThumbURL) {
		t.Errorf("BigThumbURL does not have valid protocol: %s", info.BigThumbURL)
	}
	if info.CoverURL != "" && !hasValidProtocol(info.CoverURL) {
		t.Errorf("CoverURL does not have valid protocol: %s", info.CoverURL)
	}
	if info.BigCoverURL != "" && !hasValidProtocol(info.BigCoverURL) {
		t.Errorf("BigCoverURL does not have valid protocol: %s", info.BigCoverURL)
	}
	
	// Convert to search result
	searchResult := info.ToSearchResult()
	
	t.Logf("\n=== MovieSearchResult Image URLs ===")
	t.Logf("ThumbURL: %s", searchResult.ThumbURL)
	t.Logf("CoverURL: %s", searchResult.CoverURL)
	
	// Verify search result URLs
	if searchResult.ThumbURL != "" && !hasValidProtocol(searchResult.ThumbURL) {
		t.Errorf("SearchResult ThumbURL does not have valid protocol: %s", searchResult.ThumbURL)
	}
	if searchResult.CoverURL != "" && !hasValidProtocol(searchResult.CoverURL) {
		t.Errorf("SearchResult CoverURL does not have valid protocol: %s", searchResult.CoverURL)
	}
}

func hasValidProtocol(url string) bool {
	return len(url) > 8 && (url[:7] == "http://" || url[:8] == "https://")
}
