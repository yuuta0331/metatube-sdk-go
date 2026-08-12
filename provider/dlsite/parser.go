package dlsite

import (
	"encoding/json"
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/metatube-community/metatube-sdk-go/model"
	"gorm.io/datatypes"
)

var (
	// workIDPattern matches RJ-prefixed work IDs with 6-8 digits (voice/video works only).
	// Pattern: RJ followed by 6, 7, or 8 digits
	// Examples: RJ123456, RJ01234567, RJ12345678
	workIDPattern = regexp.MustCompile(`RJ\d{6,8}`)
	
	// nonRJPattern matches non-RJ prefixed work IDs (games, comics, etc.).
	// This pattern is used to detect unsupported content types.
	// Examples: VJ012345 (games), BJ234567 (comics)
	nonRJPattern = regexp.MustCompile(`^[A-Z]{2}\d{6,8}$`)
)

// NormalizeMovieID normalizes a work ID to uppercase RJ format.
//
// This method extracts and normalizes RJ-prefixed work IDs (voice/video works).
// Only RJ-prefixed IDs with 6-8 digits are supported. Other content types
// (VJ for games, BJ for comics, etc.) are rejected.
//
// The method performs the following operations:
//  1. Converts the input to uppercase
//  2. Searches for the RJ work ID pattern (RJ followed by 6-8 digits)
//  3. Returns the extracted ID or empty string if not found
//
// Parameters:
//   - id: The work ID to normalize (can be in any case, with or without surrounding text)
//
// Returns:
//   - The normalized work ID in uppercase (e.g., "RJ123456")
//   - Empty string if the ID is invalid or has a non-RJ prefix
//
// Examples:
//
//	NormalizeMovieID("rj123456")                    // Returns: "RJ123456"
//	NormalizeMovieID("RJ123456")                    // Returns: "RJ123456"
//	NormalizeMovieID("https://dlsite.com/RJ123456") // Returns: "RJ123456"
//	NormalizeMovieID("VJ012345")                    // Returns: "" (games not supported)
//	NormalizeMovieID("invalid")                     // Returns: ""
//
// Note: This method does not return an error. Callers should check for empty string
// to determine if the ID is invalid.
func (d *DLSite) NormalizeMovieID(id string) string {
	// Convert to uppercase
	id = strings.ToUpper(id)
	
	// Extract RJ work ID pattern only
	match := workIDPattern.FindString(id)
	
	return match // Returns empty string if no match or non-RJ prefix
}

// ParseMovieIDFromURL extracts the work ID from a DLSite URL.
//
// This method parses a DLSite URL and extracts the RJ-prefixed work ID.
// The URL must be from a dlsite.com domain (including subdomains like www.dlsite.com
// and play.dlsite.com).
//
// Parameters:
//   - rawURL: The DLSite URL to parse
//
// Returns:
//   - The extracted and normalized work ID (e.g., "RJ123456")
//   - An error if the URL is invalid, not from DLSite, or doesn't contain a valid work ID
//
// Supported URL formats:
//   - https://www.dlsite.com/maniax/work/=/product_id/RJ123456.html
//   - https://www.dlsite.com/home/work/=/product_id/RJ123456.html
//   - https://play.dlsite.com/csr/=/product_id/RJ123456
//
// Examples:
//
//	ParseMovieIDFromURL("https://www.dlsite.com/maniax/work/=/product_id/RJ123456.html")
//	// Returns: "RJ123456", nil
//
//	ParseMovieIDFromURL("https://example.com/RJ123456")
//	// Returns: "", error (not a DLSite URL)
//
//	ParseMovieIDFromURL("https://www.dlsite.com/invalid")
//	// Returns: "", error (work ID not found)
//
// Errors:
//   - ErrInvalidWorkID: If the URL is malformed, not from DLSite, or doesn't contain a work ID
func (d *DLSite) ParseMovieIDFromURL(rawURL string) (string, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", NewDetailedError(ErrInvalidWorkID, map[string]interface{}{
			"url":     rawURL,
			"error":   err.Error(),
			"message": "invalid URL",
		})
	}
	
	// Check if domain is dlsite.com
	if !strings.Contains(u.Host, "dlsite.com") {
		return "", NewDetailedError(ErrInvalidWorkID, map[string]interface{}{
			"url":     rawURL,
			"host":    u.Host,
			"message": "not a DLSite URL",
		})
	}
	
	// Extract ID from URL
	id := d.NormalizeMovieID(u.String())
	if id == "" {
		return "", NewDetailedError(ErrInvalidWorkID, map[string]interface{}{
			"url":     rawURL,
			"message": "work ID not found in URL",
		})
	}
	
	return id, nil
}

// NormalizeMovieKeyword normalizes a search keyword.
//
// If the keyword contains an RJ work ID pattern, this method extracts and returns
// the work ID. Otherwise, it returns the keyword as-is for text-based search.
//
// This allows users to search by either work ID or keyword seamlessly.
//
// Parameters:
//   - keyword: The search keyword to normalize
//
// Returns:
//   - The extracted work ID if the keyword contains one, otherwise the original keyword
//
// Examples:
//
//	NormalizeMovieKeyword("RJ123456")           // Returns: "RJ123456"
//	NormalizeMovieKeyword("voice work title")   // Returns: "voice work title"
//	NormalizeMovieKeyword("check out rj123456") // Returns: "RJ123456"
func (d *DLSite) NormalizeMovieKeyword(keyword string) string {
	// If keyword contains RJ work ID pattern, extract it
	if id := d.NormalizeMovieID(keyword); id != "" {
		return id
	}
	return keyword
}

// parseJapaneseDate parses Japanese date format "YYYY年MM月DD日" to datatypes.Date.
//
// This helper function extracts year, month, and day from Japanese date strings
// commonly found on DLSite pages.
//
// Parameters:
//   - dateStr: The Japanese date string to parse (e.g., "2023年12月25日")
//
// Returns:
//   - A datatypes.Date representing the parsed date
//   - Zero value (datatypes.Date{}) if parsing fails
//
// Examples:
//
//	parseJapaneseDate("2023年12月25日") // Returns: Date(2023-12-25)
//	parseJapaneseDate("2024年1月5日")   // Returns: Date(2024-01-05)
//	parseJapaneseDate("invalid")        // Returns: Date{} (zero value)
func parseJapaneseDate(dateStr string) datatypes.Date {
	// Remove Japanese characters and parse
	re := regexp.MustCompile(`(\d{4})年(\d{1,2})月(\d{1,2})日`)
	matches := re.FindStringSubmatch(dateStr)
	if len(matches) != 4 {
		return datatypes.Date{}
	}
	
	year, _ := strconv.Atoi(matches[1])
	month, _ := strconv.Atoi(matches[2])
	day, _ := strconv.Atoi(matches[3])
	
	t := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
	return datatypes.Date(t)
}

// GetMovieInfoByID retrieves detailed metadata for a work by its ID.
//
// This method fetches comprehensive information about a DLSite voice or video work,
// including title, maker (circle name), release date, summary, cover image,
// thumbnail, and genres.
//
// Only RJ-prefixed work IDs are supported. Other content types (VJ for games,
// BJ for comics, etc.) will return an error.
//
// Parameters:
//   - id: The work ID to retrieve (e.g., "RJ123456", "rj123456")
//
// Returns:
//   - A *model.MovieInfo containing the work's metadata
//   - An error if the ID is invalid, unsupported, or the work is not found
//
// The returned MovieInfo includes:
//   - ID, Number: The normalized work ID
//   - Title: The work title
//   - Maker: The circle/maker name
//   - ReleaseDate: The release/sale date
//   - Summary: The work description
//   - CoverURL: The cover image URL
//   - ThumbURL: The thumbnail image URL
//   - Genres: Array of genre/tag strings
//   - Provider: "dlsite"
//   - Homepage: The work page URL
//
// Examples:
//
//	info, err := provider.GetMovieInfoByID("RJ123456")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Printf("Title: %s\nMaker: %s\n", info.Title, info.Maker)
//
// Errors:
//   - ErrInvalidWorkID: If the ID format is invalid
//   - ErrUnsupportedContentType: If the ID is not RJ-prefixed (e.g., VJ, BJ)
//   - ErrWorkNotFound: If the work doesn't exist (HTTP 404)
//   - ErrNetworkError: If the HTTP request fails
//   - ErrValidationError: If required metadata fields are missing
func (d *DLSite) GetMovieInfoByID(id string) (info *model.MovieInfo, err error) {
	// Store original ID for error messages
	originalID := id
	
	// Normalize ID
	id = d.NormalizeMovieID(id)
	if id == "" {
		// Check if it's a non-RJ work ID
		upperID := strings.ToUpper(originalID)
		if nonRJPattern.MatchString(upperID) {
			return nil, NewDetailedError(ErrUnsupportedContentType, map[string]interface{}{
				"work_id": originalID,
				"message": "only RJ-prefixed voice/video works are supported (VJ, BJ, etc. are not supported)",
			})
		}
		return nil, NewDetailedError(ErrInvalidWorkID, map[string]interface{}{
			"work_id": originalID,
		})
	}
	
	// Construct work page URL using buildURL helper
	workURL := buildURL(baseURL, "/maniax/work/=/product_id/", id, ".html")
	
	// Fetch and parse
	return d.GetMovieInfoByURL(workURL)
}

// GetMovieInfoByURL retrieves detailed metadata for a work by its URL.
//
// This method fetches metadata from DLSite using the JSON API as the primary method,
// with HTML parsing as a fallback. The JSON API provides more reliable data extraction
// compared to HTML parsing which is subject to page structure changes.
//
// The method automatically:
//   - Extracts the work ID from the URL
//   - Attempts to fetch metadata from the JSON API first
//   - Falls back to HTML parsing if JSON API fails
//   - Validates that required fields are present
//
// Parameters:
//   - rawURL: The DLSite work page URL
//
// Returns:
//   - A *model.MovieInfo containing the work's metadata
//   - An error if the request fails, parsing fails, or validation fails
//
// Extracted fields:
//   - Title: Work name
//   - Maker: Circle/maker name
//   - ReleaseDate: Registration/sale date
//   - Summary: Work description
//   - CoverURL: Main cover image URL
//   - ThumbURL: Thumbnail image URL
//   - Genres: Array of genre/tag strings
//
// Examples:
//
//	info, err := provider.GetMovieInfoByURL("https://www.dlsite.com/maniax/work/=/product_id/RJ123456.html")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Printf("Title: %s\n", info.Title)
//
// Errors:
//   - ErrWorkNotFound: If the work doesn't exist (HTTP 404)
//   - ErrNetworkError: If the HTTP request fails or returns an error status
//   - ErrValidationError: If required metadata fields are missing
//
// Note: This method is also called internally by GetMovieInfoByID.
func (d *DLSite) GetMovieInfoByURL(rawURL string) (info *model.MovieInfo, err error) {
	// Initialize info with pre-allocated genres slice
	info = &model.MovieInfo{
		Provider: d.Name(),
		Homepage: rawURL,
		Genres:   make([]string, 0, 10), // Pre-allocate with reasonable capacity
	}
	
	// Extract ID from URL
	info.ID = d.NormalizeMovieID(rawURL)
	info.Number = info.ID
	
	// Try JSON API first (more reliable)
	jsonInfo, jsonErr := d.getMovieInfoFromJSON(info.ID)
	if jsonErr == nil && jsonInfo != nil {
		// Successfully got data from JSON API
		jsonInfo.Provider = d.Name()
		jsonInfo.Homepage = rawURL
		return jsonInfo, nil
	}
	
	// Fallback to HTML parsing if JSON API fails
	return d.getMovieInfoFromHTML(rawURL, info)
}

// getMovieInfoFromJSON fetches metadata from DLSite's JSON API.
// This function does not retain references to large temporary objects.
// The jsonData variable is local and will be garbage collected after the function returns.
func (d *DLSite) getMovieInfoFromJSON(workID string) (*model.MovieInfo, error) {
	// Use buildURL helper for efficient URL construction
	apiURL := buildURL(baseURL, "/maniax/api/=/product.json?workno=", workID)
	
	c := d.ClonedCollector()
	
	var jsonData []map[string]interface{}
	var fetchErr error
	
	c.OnResponse(func(r *colly.Response) {
		// Parse JSON response - jsonData is local and will be GC'd
		if err := json.Unmarshal(r.Body, &jsonData); err != nil {
			fetchErr = NewDetailedError(ErrParseError, map[string]interface{}{
				"error":   err.Error(),
				"url":     apiURL,
				"message": "failed to parse JSON response",
			})
			return
		}
		
		if len(jsonData) == 0 {
			fetchErr = NewDetailedError(ErrWorkNotFound, map[string]interface{}{
				"url":     apiURL,
				"message": "empty JSON response",
			})
			return
		}
	})
	
	c.OnError(func(r *colly.Response, e error) {
		switch r.StatusCode {
		case 404:
			fetchErr = NewDetailedErrorWithStatus(ErrWorkNotFound, r.StatusCode, apiURL, map[string]interface{}{
				"work_id": workID,
			})
		case 500, 502, 503, 504:
			fetchErr = NewDetailedErrorWithStatus(ErrNetworkError, r.StatusCode, apiURL, map[string]interface{}{
				"work_id": workID,
				"message": "server error",
			})
		default:
			if r.StatusCode >= 400 {
				fetchErr = NewDetailedErrorWithStatus(ErrNetworkError, r.StatusCode, apiURL, map[string]interface{}{
					"work_id": workID,
				})
			} else if e != nil {
				fetchErr = fmt.Errorf("%w: %v", ErrNetworkError, e)
			}
		}
	})
	
	if e := c.Visit(apiURL); e != nil {
		return nil, fmt.Errorf("%w: %v", ErrNetworkError, e)
	}
	
	if fetchErr != nil {
		return nil, fetchErr
	}
	
	if len(jsonData) == 0 {
		return nil, fmt.Errorf("%w: no data returned", ErrWorkNotFound)
	}
	
	// Extract first element - we don't retain the full jsonData slice
	data := jsonData[0]
	// Clear jsonData reference to allow GC (optional, but explicit)
	jsonData = nil
	
	info := &model.MovieInfo{
		ID:     workID,
		Number: workID,
		Genres: make([]string, 0, 10), // Pre-allocate with reasonable capacity
	}
	
	// Extract title
	if workName, ok := data["work_name"].(string); ok {
		info.Title = workName
	}
	
	// Extract maker
	if makerName, ok := data["maker_name"].(string); ok {
		info.Maker = makerName
	}
	
	// Extract release date
	if registDate, ok := data["regist_date"].(string); ok {
		// Parse date format: "2025-09-04 00:00:00"
		if t, err := time.Parse("2006-01-02 15:04:05", registDate); err == nil {
			info.ReleaseDate = datatypes.Date(t)
		}
	}
	
	// Extract summary
	if introS, ok := data["intro_s"].(string); ok {
		info.Summary = introS
	}
	
	// Extract cover image
	if imageMain, ok := data["image_main"].(map[string]interface{}); ok {
		if url, ok := imageMain["url"].(string); ok {
			// Use normalizeImageURL helper
			normalizedURL := normalizeImageURL(url)
			info.CoverURL = normalizedURL
			info.BigCoverURL = normalizedURL  // Use same URL for BigCoverURL
		}
	}
	
	// Extract thumbnail
	if imageThum, ok := data["image_thum"].(map[string]interface{}); ok {
		if url, ok := imageThum["url"].(string); ok {
			// Use normalizeImageURL helper
			normalizedURL := normalizeImageURL(url)
			info.ThumbURL = normalizedURL
			info.BigThumbURL = normalizedURL  // Use same URL for BigThumbURL
		}
	}
	
	// Extract genres with pre-allocated slice
	if genres, ok := data["genres"].([]interface{}); ok {
		// Pre-allocate slice with appropriate capacity
		info.Genres = make([]string, 0, len(genres))
		for _, g := range genres {
			if genreMap, ok := g.(map[string]interface{}); ok {
				if name, ok := genreMap["name"].(string); ok && name != "" {
					info.Genres = append(info.Genres, name)
				}
			}
		}
	}
	
	return info, nil
}

// getMovieInfoFromHTML fetches metadata by parsing HTML (fallback method).
// This function does not retain references to large temporary objects.
// The colly collector is local and will be garbage collected after the function returns.
func (d *DLSite) getMovieInfoFromHTML(rawURL string, info *model.MovieInfo) (*model.MovieInfo, error) {
	var err error
	
	// Create a local collector - will be GC'd after function returns
	c := d.ClonedCollector()
	
	// Extract title
	c.OnHTML("#work_name a, h1#work_name", func(e *colly.HTMLElement) {
		if info.Title == "" {
			info.Title = strings.TrimSpace(e.Text)
		}
	})
	
	// Extract maker (circle)
	c.OnHTML("span[class*='maker_name'] a", func(e *colly.HTMLElement) {
		if info.Maker == "" {
			info.Maker = strings.TrimSpace(e.Text)
		}
	})
	
	// Extract release date
	c.OnHTML("table#work_outline tr", func(e *colly.HTMLElement) {
		if strings.Contains(e.ChildText("th"), "販売日") {
			releaseDateStr := strings.TrimSpace(e.ChildText("td"))
			if releaseDateStr != "" {
				info.ReleaseDate = parseJapaneseDate(releaseDateStr)
			}
		}
	})
	
	// Extract summary
	c.OnHTML("div.work_parts_area div[itemprop='description']", func(e *colly.HTMLElement) {
		if info.Summary == "" {
			info.Summary = strings.TrimSpace(e.Text)
		}
	})
	
	// Extract cover image
	c.OnHTML("div.product-slider-data div.product-slider-item img", func(e *colly.HTMLElement) {
		if info.CoverURL == "" {
			info.CoverURL = e.Attr("src")
			info.BigCoverURL = e.Attr("src")  // Use same URL for BigCoverURL
		}
	})
	
	// Extract thumbnail
	c.OnHTML("li.slider_item img", func(e *colly.HTMLElement) {
		if info.ThumbURL == "" {
			info.ThumbURL = e.Attr("src")
			info.BigThumbURL = e.Attr("src")  // Use same URL for BigThumbURL
		}
	})
	
	// Extract genres
	c.OnHTML("div.main_genre a, div.genre a", func(e *colly.HTMLElement) {
		genre := strings.TrimSpace(e.Text)
		if genre != "" {
			info.Genres = append(info.Genres, genre)
		}
	})
	
	// Handle errors
	c.OnError(func(r *colly.Response, e error) {
		switch r.StatusCode {
		case 404:
			err = NewDetailedErrorWithStatus(ErrWorkNotFound, r.StatusCode, rawURL, map[string]interface{}{
				"work_id": info.ID,
			})
		case 500, 502, 503, 504:
			err = NewDetailedErrorWithStatus(ErrNetworkError, r.StatusCode, rawURL, map[string]interface{}{
				"work_id": info.ID,
				"message": "server error",
			})
		default:
			if r.StatusCode >= 400 {
				err = NewDetailedErrorWithStatus(ErrNetworkError, r.StatusCode, rawURL, map[string]interface{}{
					"work_id": info.ID,
				})
			} else if e != nil {
				err = fmt.Errorf("%w: %v", ErrNetworkError, e)
			}
		}
	})
	
	// Visit the URL
	if e := c.Visit(rawURL); e != nil {
		return nil, fmt.Errorf("%w: %v", ErrNetworkError, e)
	}
	
	// Return early if there was an error during scraping
	if err != nil {
		return nil, err
	}
	
	// Validate required fields
	if !info.IsValid() {
		var missingFields []string
		if info.ID == "" {
			missingFields = append(missingFields, "ID")
		}
		if info.Number == "" {
			missingFields = append(missingFields, "Number")
		}
		if info.Title == "" {
			missingFields = append(missingFields, "Title")
		}
		if info.CoverURL == "" {
			missingFields = append(missingFields, "CoverURL")
		}
		if info.Provider == "" {
			missingFields = append(missingFields, "Provider")
		}
		if info.Homepage == "" {
			missingFields = append(missingFields, "Homepage")
		}
		
		return nil, NewDetailedError(ErrValidationError, map[string]interface{}{
			"missing_fields": missingFields,
			"work_id":        info.ID,
			"url":            rawURL,
		})
	}
	
	return info, nil
}

// SearchMovie searches for works by keyword.
//
// This method performs a keyword-based search on DLSite and returns matching works.
// Only RJ-prefixed voice/video works are included in the results; other content
// types (VJ for games, BJ for comics, etc.) are automatically filtered out.
//
// Search behavior:
//   - If the keyword contains an RJ work ID, it performs a direct ID lookup
//   - Otherwise, it performs a keyword search on DLSite
//   - Results are filtered to include only voice/video works (RJ-prefixed)
//   - Age verification cookie is automatically included in the request
//
// Parameters:
//   - keyword: The search keyword or work ID
//
// Returns:
//   - A slice of *model.MovieSearchResult containing matching works
//   - An error if the search request fails
//
// Each MovieSearchResult includes:
//   - ID, Number: The work ID
//   - Title: The work title
//   - ThumbURL: The thumbnail image URL
//   - Homepage: The work page URL
//   - Provider: "dlsite"
//
// Examples:
//
//	// Search by keyword
//	results, err := provider.SearchMovie("voice work")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	for _, result := range results {
//	    fmt.Printf("%s: %s\n", result.ID, result.Title)
//	}
//
//	// Search by work ID
//	results, err := provider.SearchMovie("RJ123456")
//	// Returns detailed info for that specific work
//
// Errors:
//   - ErrNetworkError: If the search request fails
//   - Other errors from GetMovieInfoByID if searching by work ID
//
// Note: Empty results (no matches) are returned as an empty slice, not an error.
func (d *DLSite) SearchMovie(keyword string) (results []*model.MovieSearchResult, err error) {
	// Normalize keyword
	keyword = d.NormalizeMovieKeyword(keyword)
	
	// If keyword is an RJ work ID, search by ID directly
	if regexp.MustCompile(`^RJ\d{6,8}$`).MatchString(keyword) {
		info, err := d.GetMovieInfoByID(keyword)
		if err != nil {
			return nil, err
		}
		return []*model.MovieSearchResult{info.ToSearchResult()}, nil
	}
	
	// Construct search URL using buildURL helper
	escapedKeyword := url.QueryEscape(keyword)
	searchURL := buildURL(baseURL, "/maniax/fsr/=/language/jp/keyword/", escapedKeyword, "/per_page/30/")
	
	// Pre-allocate results slice with reasonable capacity (30 is the per_page limit)
	results = make([]*model.MovieSearchResult, 0, 30)
	
	// Create a local collector - will be GC'd after function returns
	c := d.ClonedCollector()
	
	// Extract results
	c.OnHTML("li.search_result_img_box_inner", func(e *colly.HTMLElement) {
		result := &model.MovieSearchResult{
			Provider: d.Name(),
		}
		
		// Extract work link and ID
		href := e.ChildAttr("dd.work_name a", "href")
		if href != "" {
			result.Homepage = href
			result.ID = d.NormalizeMovieID(href)
			result.Number = result.ID
		}
		
		// Only process if ID is RJ-prefixed (filter out VJ, BJ, etc.)
		if result.ID == "" || !strings.HasPrefix(result.ID, "RJ") {
			return // Skip non-voice/video works
		}
		
		// Extract title
		result.Title = strings.TrimSpace(e.ChildText("dd.work_name a"))
		
		// Extract thumbnail and cover
		thumbSrc := e.ChildAttr("dt.search_img img", "src")
		if thumbSrc != "" {
			// Use normalizeImageURL helper
			thumbSrc = normalizeImageURL(thumbSrc)
			result.ThumbURL = thumbSrc
			// Convert thumbnail URL to main image URL for better poster display in Jellyfin
			// Pattern: *_img_sam.jpg -> *_img_main.jpg
			result.CoverURL = strings.Replace(thumbSrc, "_img_sam.jpg", "_img_main.jpg", 1)
		}
		
		// Only add if valid
		if result.IsValid() {
			results = append(results, result)
		}
	})
	
	// Handle errors
	c.OnError(func(r *colly.Response, e error) {
		switch r.StatusCode {
		case 404:
			err = NewDetailedErrorWithStatus(ErrNetworkError, r.StatusCode, searchURL, map[string]interface{}{
				"keyword": keyword,
				"message": "search endpoint not found",
			})
		case 500, 502, 503, 504:
			err = NewDetailedErrorWithStatus(ErrNetworkError, r.StatusCode, searchURL, map[string]interface{}{
				"keyword": keyword,
				"message": "server error during search",
			})
		default:
			if r.StatusCode >= 400 {
				err = NewDetailedErrorWithStatus(ErrNetworkError, r.StatusCode, searchURL, map[string]interface{}{
					"keyword": keyword,
				})
			} else if e != nil {
				err = fmt.Errorf("%w: search failed: %v", ErrNetworkError, e)
			}
		}
	})
	
	// Visit the search URL
	if e := c.Visit(searchURL); e != nil {
		return nil, e
	}
	
	return results, err
}

