package dlsite

import (
	"net/http"
	"testing"
	"time"
)

// TestImageURLAccessibility tests that the image URLs are accessible
func TestImageURLAccessibility(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping image URL accessibility test in short mode")
	}

	testCases := []struct {
		name string
		url  string
	}{
		{
			name: "Thumbnail Image (_img_sam.jpg)",
			url:  "https://img.dlsite.jp/modpub/images2/work/doujin/RJ01228000/RJ01227569_img_sam.jpg",
		},
		{
			name: "Main Image (_img_main.jpg)",
			url:  "https://img.dlsite.jp/modpub/images2/work/doujin/RJ01228000/RJ01227569_img_main.jpg",
		},
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Test 1: Without Referer
			t.Log("--- Test without Referer ---")
			req1, err := http.NewRequest("GET", tc.url, nil)
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}
			req1.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")

			resp1, err := client.Do(req1)
			if err != nil {
				t.Fatalf("Failed to fetch image without Referer: %v", err)
			}
			defer resp1.Body.Close()

			t.Logf("URL: %s", tc.url)
			t.Logf("Status Code (no Referer): %d", resp1.StatusCode)
			t.Logf("Content-Type (no Referer): %s", resp1.Header.Get("Content-Type"))

			// Test 2: With Referer
			t.Log("--- Test with Referer ---")
			req2, err := http.NewRequest("GET", tc.url, nil)
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}
			req2.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")
			req2.Header.Set("Referer", "https://www.dlsite.com/")

			resp2, err := client.Do(req2)
			if err != nil {
				t.Fatalf("Failed to fetch image with Referer: %v", err)
			}
			defer resp2.Body.Close()

			t.Logf("Status Code (with Referer): %d", resp2.StatusCode)
			t.Logf("Content-Type (with Referer): %s", resp2.Header.Get("Content-Type"))
			t.Logf("Content-Length (with Referer): %s", resp2.Header.Get("Content-Length"))

			if resp2.StatusCode != http.StatusOK {
				t.Errorf("Expected status 200, got %d", resp2.StatusCode)
			}

			contentType := resp2.Header.Get("Content-Type")
			if contentType != "image/jpeg" && contentType != "image/jpg" {
				t.Errorf("Expected image content type, got %s", contentType)
			}
		})
	}
}


// TestLargerImageVariants tests if larger image variants are available
func TestLargerImageVariants(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping larger image variants test in short mode")
	}

	baseURL := "https://img.dlsite.jp/modpub/images2/work/doujin/RJ01228000/RJ01227569"
	
	variants := []string{
		"_img_main.jpg",      // Standard main image
		"_img_main.webp",     // WebP version
		"_img_smp.jpg",       // Sample image
		"_img_sam.jpg",       // Thumbnail
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	for _, variant := range variants {
		url := baseURL + variant
		t.Run(variant, func(t *testing.T) {
			req, err := http.NewRequest("GET", url, nil)
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}

			req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")
			req.Header.Set("Referer", "https://www.dlsite.com/")

			resp, err := client.Do(req)
			if err != nil {
				t.Logf("Failed to fetch %s: %v", variant, err)
				return
			}
			defer resp.Body.Close()

			t.Logf("URL: %s", url)
			t.Logf("Status Code: %d", resp.StatusCode)
			t.Logf("Content-Type: %s", resp.Header.Get("Content-Type"))
			t.Logf("Content-Length: %s bytes", resp.Header.Get("Content-Length"))

			if resp.StatusCode == http.StatusOK {
				t.Logf("✓ %s is available", variant)
			} else {
				t.Logf("✗ %s is not available (status: %d)", variant, resp.StatusCode)
			}
		})
	}
}
