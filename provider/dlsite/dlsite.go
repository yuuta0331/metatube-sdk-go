// Package dlsite provides a metadata provider for DLSite (https://www.dlsite.com/).
//
// This provider implements the metatube-sdk-go provider interfaces to retrieve
// metadata for voice and video works (RJ-prefixed IDs) from DLSite.
//
// Supported Content Types:
//   - Voice works (RJ-prefixed IDs)
//   - Video works (RJ-prefixed IDs)
//
// Unsupported Content Types:
//   - Games (VJ-prefixed IDs)
//   - Comics/Manga (BJ-prefixed IDs)
//   - Other content types
//
// Features:
//   - Work ID normalization and validation (RJ-prefixed only)
//   - URL parsing to extract work IDs
//   - Keyword-based search (filtered to voice/video works only)
//   - Detailed metadata extraction (title, maker, release date, summary, images, genres)
//   - Automatic age verification cookie handling
//   - Content type filtering to exclude non-voice/video works
//
// Example Usage:
//
//	provider := dlsite.New()
//
//	// Get metadata by work ID
//	info, err := provider.GetMovieInfoByID("RJ123456")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Printf("Title: %s\n", info.Title)
//
//	// Search for works
//	results, err := provider.SearchMovie("keyword")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	for _, result := range results {
//	    fmt.Printf("Found: %s (%s)\n", result.Title, result.ID)
//	}
package dlsite

import (
	"net/http"
	"net/url"
	"time"

	"github.com/metatube-community/metatube-sdk-go/provider"
	"github.com/metatube-community/metatube-sdk-go/provider/internal/scraper"
	"golang.org/x/text/language"
)

const (
	// Name is the provider identifier
	Name = "dlsite"
	// Priority is the default priority for this provider
	Priority = 1000
)

// baseURL is the base URL for DLSite
var baseURL = "https://www.dlsite.com"

func init() {
	provider.Register(Name, factory)
}

// factory is the no-argument constructor required by provider.Register.
// It creates a DLSite provider with default options.
func factory() *DLSite {
	return New()
}

// DLSite is a provider for DLSite voice and video works.
// It implements the MovieProvider and MovieSearcher interfaces.
//
// This provider only supports RJ-prefixed work IDs (voice/video works).
// Other content types (VJ for games, BJ for comics, etc.) are not supported.
type DLSite struct {
	*scraper.Scraper
	options    *Options
	logger     Logger
	httpClient *HTTPClient
}

// New creates a new DLSite provider instance with optional configuration.
//
// The provider is configured with:
//   - Base URL: https://www.dlsite.com
//   - Language: Japanese
//   - Age verification cookie (adultchecked=1) for automatic adult content access
//   - Default timeout: 30 seconds (configurable via options)
//
// Parameters:
//   - opts: Optional configuration options (see WithTimeout, WithCache, etc.)
//
// Returns:
//   - A configured DLSite provider instance
//
// Example (basic usage):
//
//	provider := dlsite.New()
//	info, err := provider.GetMovieInfoByID("RJ123456")
//
// Example (with options):
//
//	provider := dlsite.New(
//	    dlsite.WithTimeout(60 * time.Second),
//	    dlsite.WithCache(1*time.Hour, 1000),
//	    dlsite.WithRateLimit(10),
//	)
func New(opts ...Option) *DLSite {
	// Create default options first
	options := defaultOptions()
	
	// Create a temporary DLSite to apply options and determine HTTP client config
	temp := &DLSite{options: options}
	for _, opt := range opts {
		opt(temp)
	}
	
	// Initialize HTTP client based on options
	var httpClient *HTTPClient
	if temp.options.HTTPClient != nil {
		// Wrap the custom HTTP client
		httpClient = &HTTPClient{
			client:    temp.options.HTTPClient,
			userAgent: temp.options.UserAgent,
		}
	} else {
		// Create optimized HTTP client with default settings
		timeout := temp.options.RequestTimeout
		if timeout == 0 {
			timeout = 30 * time.Second
		}
		httpClient = NewHTTPClient(timeout, temp.options.UserAgent)
	}
	
	// Create scraper with optimized transport
	scraperOpts := []scraper.Option{
		scraper.WithCookies(baseURL, []*http.Cookie{
			{Name: "adultchecked", Value: "1"},
		}),
	}
	
	// Add transport if available
	if transport := httpClient.Transport(); transport != nil {
		scraperOpts = append(scraperOpts, scraper.WithTransport(transport))
	}
	
	d := &DLSite{
		Scraper:    scraper.NewDefaultScraper(Name, baseURL, Priority, language.Japanese, scraperOpts...),
		options:    temp.options,
		httpClient: httpClient,
	}

	// Initialize logger if provided
	if d.options.Logger != nil {
		d.logger = d.options.Logger
	}

	// Apply timeout configuration if set
	if d.options.RequestTimeout > 0 {
		d.Scraper.SetRequestTimeout(d.options.RequestTimeout)
	}

	return d
}

// Name returns the provider name "dlsite".
//
// This method implements the Provider interface.
//
// Returns:
//   - The string "dlsite"
func (d *DLSite) Name() string {
	return Name
}

// URL returns the base URL of the DLSite website.
//
// This method implements the Provider interface.
//
// Returns:
//   - A *url.URL pointing to https://www.dlsite.com
func (d *DLSite) URL() *url.URL {
	u, _ := url.Parse(baseURL)
	return u
}

// Language returns the primary language tag for this provider.
//
// This method implements the Provider interface.
//
// Returns:
//   - language.Japanese (ja)
func (d *DLSite) Language() language.Tag {
	return language.Japanese
}
