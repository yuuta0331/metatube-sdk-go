package dlsite

import (
	"os"
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
)

// TestWorkPageHTMLStructure validates that work_page.html contains all required elements
func TestWorkPageHTMLStructure(t *testing.T) {
	// Read the test HTML file
	file, err := os.Open("testdata/work_page.html")
	if err != nil {
		t.Fatalf("Failed to open work_page.html: %v", err)
	}
	defer file.Close()

	// Parse HTML with goquery
	doc, err := goquery.NewDocumentFromReader(file)
	if err != nil {
		t.Fatalf("Failed to parse HTML: %v", err)
	}

	// Test Title extraction
	t.Run("Title", func(t *testing.T) {
		title := strings.TrimSpace(doc.Find("#work_name a, h1#work_name").First().Text())
		if title == "" {
			t.Error("Title not found or empty")
		}
		if title != "サンプル音声作品" {
			t.Errorf("Expected title 'サンプル音声作品', got '%s'", title)
		}
	})

	// Test Maker (Circle) extraction
	t.Run("Maker", func(t *testing.T) {
		maker := strings.TrimSpace(doc.Find("span[class*='maker_name'] a").First().Text())
		if maker == "" {
			t.Error("Maker not found or empty")
		}
		if maker != "テストサークル" {
			t.Errorf("Expected maker 'テストサークル', got '%s'", maker)
		}
	})

	// Test Release Date extraction
	t.Run("ReleaseDate", func(t *testing.T) {
		var releaseDate string
		doc.Find("table#work_outline tr").Each(func(i int, s *goquery.Selection) {
			th := strings.TrimSpace(s.Find("th").Text())
			if strings.Contains(th, "販売日") {
				releaseDate = strings.TrimSpace(s.Find("td").Text())
			}
		})
		if releaseDate == "" {
			t.Error("Release date not found or empty")
		}
		if releaseDate != "2024年01月15日" {
			t.Errorf("Expected release date '2024年01月15日', got '%s'", releaseDate)
		}
	})

	// Test Summary extraction
	t.Run("Summary", func(t *testing.T) {
		summary := strings.TrimSpace(doc.Find("div.work_parts_area div[itemprop='description']").Text())
		if summary == "" {
			t.Error("Summary not found or empty")
		}
		if !strings.Contains(summary, "これはテスト用のサンプル作品です") {
			t.Error("Summary doesn't contain expected text")
		}
	})

	// Test Cover Image extraction
	t.Run("CoverImage", func(t *testing.T) {
		coverSrc, exists := doc.Find("div.product-slider-data div.product-slider-item img").First().Attr("src")
		if !exists || coverSrc == "" {
			t.Error("Cover image not found or empty")
		}
		expectedURL := "https://img.dlsite.jp/modpub/images2/work/doujin/RJ123000/RJ123456_img_main.jpg"
		if coverSrc != expectedURL {
			t.Errorf("Expected cover URL '%s', got '%s'", expectedURL, coverSrc)
		}
	})

	// Test Thumbnail extraction
	t.Run("Thumbnail", func(t *testing.T) {
		thumbSrc, exists := doc.Find("li.slider_item img").First().Attr("src")
		if !exists || thumbSrc == "" {
			t.Error("Thumbnail not found or empty")
		}
		expectedURL := "https://img.dlsite.jp/modpub/images2/work/doujin/RJ123000/RJ123456_img_smp1.jpg"
		if thumbSrc != expectedURL {
			t.Errorf("Expected thumbnail URL '%s', got '%s'", expectedURL, thumbSrc)
		}
	})

	// Test Genres/Tags extraction
	t.Run("Genres", func(t *testing.T) {
		var genres []string
		doc.Find("div.main_genre a, div.genre a").Each(func(i int, s *goquery.Selection) {
			genre := strings.TrimSpace(s.Text())
			if genre != "" {
				genres = append(genres, genre)
			}
		})
		if len(genres) == 0 {
			t.Error("No genres found")
		}
		expectedGenres := []string{"ボイス・ASMR", "癒し", "バイノーラル", "耳かき", "囁き"}
		if len(genres) != len(expectedGenres) {
			t.Errorf("Expected %d genres, got %d", len(expectedGenres), len(genres))
		}
	})
}

// TestSearchPageHTMLStructure validates that search_page.html contains all required elements
func TestSearchPageHTMLStructure(t *testing.T) {
	// Read the test HTML file
	file, err := os.Open("testdata/search_page.html")
	if err != nil {
		t.Fatalf("Failed to open search_page.html: %v", err)
	}
	defer file.Close()

	// Parse HTML with goquery
	doc, err := goquery.NewDocumentFromReader(file)
	if err != nil {
		t.Fatalf("Failed to parse HTML: %v", err)
	}

	// Test search results extraction
	t.Run("SearchResults", func(t *testing.T) {
		results := doc.Find("li.search_result_img_box_inner")
		if results.Length() == 0 {
			t.Fatal("No search results found")
		}
		
		expectedCount := 4 // RJ123456, VJ234567, RJ789012, BJ345678
		if results.Length() != expectedCount {
			t.Errorf("Expected %d search results, got %d", expectedCount, results.Length())
		}
	})

	// Test Work ID extraction and filtering
	t.Run("WorkIDExtraction", func(t *testing.T) {
		var workIDs []string
		doc.Find("li.search_result_img_box_inner").Each(func(i int, s *goquery.Selection) {
			href, exists := s.Find("dd.work_name a").Attr("href")
			if exists {
				// Extract work ID from href
				dlsite := New()
				id := dlsite.NormalizeMovieID(href)
				workIDs = append(workIDs, id)
			}
		})

		if len(workIDs) != 4 {
			t.Errorf("Expected 4 work IDs, got %d", len(workIDs))
		}

		// Check that we have RJ work IDs
		hasRJ := false
		for _, id := range workIDs {
			if strings.HasPrefix(id, "RJ") {
				hasRJ = true
				break
			}
		}

		if !hasRJ {
			t.Error("Expected at least one RJ work ID")
		}
		// Note: VJ and BJ will be empty strings after normalization since they're filtered
	})

	// Test RJ-only filtering
	t.Run("RJOnlyFiltering", func(t *testing.T) {
		dlsite := New()
		var rjCount int
		
		doc.Find("li.search_result_img_box_inner").Each(func(i int, s *goquery.Selection) {
			href, exists := s.Find("dd.work_name a").Attr("href")
			if exists {
				id := dlsite.NormalizeMovieID(href)
				// Only RJ IDs should be non-empty after normalization
				if id != "" && strings.HasPrefix(id, "RJ") {
					rjCount++
				}
			}
		})

		expectedRJCount := 2 // RJ123456 and RJ789012
		if rjCount != expectedRJCount {
			t.Errorf("Expected %d RJ works, got %d", expectedRJCount, rjCount)
		}
	})

	// Test Title extraction
	t.Run("TitleExtraction", func(t *testing.T) {
		firstTitle := strings.TrimSpace(doc.Find("li.search_result_img_box_inner").First().Find("dd.work_name a").Text())
		if firstTitle == "" {
			t.Error("First result title not found or empty")
		}
		if firstTitle != "サンプル音声作品" {
			t.Errorf("Expected first title 'サンプル音声作品', got '%s'", firstTitle)
		}
	})

	// Test Thumbnail extraction
	t.Run("ThumbnailExtraction", func(t *testing.T) {
		thumbSrc, exists := doc.Find("li.search_result_img_box_inner").First().Find("dt.search_img img").Attr("src")
		if !exists || thumbSrc == "" {
			t.Error("First result thumbnail not found or empty")
		}
		expectedURL := "https://img.dlsite.jp/modpub/images2/work/doujin/RJ123000/RJ123456_img_main.jpg"
		if thumbSrc != expectedURL {
			t.Errorf("Expected thumbnail URL '%s', got '%s'", expectedURL, thumbSrc)
		}
	})

	// Test Homepage extraction
	t.Run("HomepageExtraction", func(t *testing.T) {
		href, exists := doc.Find("li.search_result_img_box_inner").First().Find("dd.work_name a").Attr("href")
		if !exists || href == "" {
			t.Error("First result homepage not found or empty")
		}
		if !strings.Contains(href, "RJ123456") {
			t.Errorf("Expected homepage to contain 'RJ123456', got '%s'", href)
		}
	})
}
