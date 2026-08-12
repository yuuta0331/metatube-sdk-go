package dlsite

import (
	"context"
	"crypto/tls"
	"net"
	"net/http"
	"time"
)

// HTTPClient wraps an HTTP client with optimized settings for DLSite requests.
// It provides connection pooling, timeout configuration, and proper transport settings
// for efficient HTTP communication.
//
// The client is configured with:
//   - Connection pooling (MaxIdleConns: 100, MaxIdleConnsPerHost: 10)
//   - Optimized timeouts (Dial: 10s, TLS Handshake: 10s, Request: 30s)
//   - Keep-Alive enabled (30s)
//   - Automatic gzip decompression
//   - TLS 1.2+ with secure cipher suites
//
// This client is designed to be safe for concurrent use by multiple goroutines.
type HTTPClient struct {
	client    *http.Client
	userAgent string
}

// NewHTTPClient creates a new optimized HTTP client for DLSite requests.
//
// The client is configured with connection pooling and appropriate timeouts
// to maximize performance while maintaining reliability.
//
// Parameters:
//   - timeout: Overall request timeout (0 = no timeout)
//   - userAgent: Custom User-Agent header (empty = use default)
//
// Returns:
//   - A configured HTTPClient instance
//
// Example:
//
//	client := NewHTTPClient(30*time.Second, "MyApp/1.0")
//	resp, err := client.Do(req)
func NewHTTPClient(timeout time.Duration, userAgent string) *HTTPClient {
	// Create optimized transport with connection pooling
	transport := &http.Transport{
		// Connection pooling settings
		MaxIdleConns:        100,              // Maximum idle connections across all hosts
		MaxIdleConnsPerHost: 10,               // Maximum idle connections per host
		MaxConnsPerHost:     0,                // No limit on total connections per host
		IdleConnTimeout:     90 * time.Second, // How long idle connections are kept

		// Timeout settings
		DialContext: (&net.Dialer{
			Timeout:   10 * time.Second, // Connection timeout
			KeepAlive: 30 * time.Second, // Keep-alive probe interval
		}).DialContext,
		TLSHandshakeTimeout:   10 * time.Second, // TLS handshake timeout
		ResponseHeaderTimeout: 10 * time.Second, // Time to wait for response headers
		ExpectContinueTimeout: 1 * time.Second,  // Time to wait for 100-Continue response

		// Performance settings
		DisableCompression: false,                // Enable gzip compression
		DisableKeepAlives:  false,                // Enable HTTP keep-alive
		ForceAttemptHTTP2:  true,                 // Try HTTP/2 if available
		MaxResponseHeaderBytes: 10 * 1024 * 1024, // 10MB max response headers

		// Security settings
		TLSClientConfig: &tls.Config{
			MinVersion: tls.VersionTLS12, // Require TLS 1.2 or higher
			// Use secure cipher suites (Go's default is secure)
		},
	}

	// Create HTTP client with the optimized transport
	client := &http.Client{
		Transport: transport,
		Timeout:   timeout,
		// Don't follow redirects automatically - let colly handle them
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	return &HTTPClient{
		client:    client,
		userAgent: userAgent,
	}
}

// Do executes an HTTP request using the optimized client.
//
// This method adds the configured User-Agent header if set and executes
// the request with connection pooling and timeout settings.
//
// Parameters:
//   - req: The HTTP request to execute
//
// Returns:
//   - The HTTP response
//   - An error if the request fails
//
// Example:
//
//	req, _ := http.NewRequest("GET", "https://example.com", nil)
//	resp, err := client.Do(req)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer resp.Body.Close()
func (c *HTTPClient) Do(req *http.Request) (*http.Response, error) {
	// Add User-Agent header if configured
	if c.userAgent != "" && req.Header.Get("User-Agent") == "" {
		req.Header.Set("User-Agent", c.userAgent)
	}

	// Execute the request
	return c.client.Do(req)
}

// DoWithContext executes an HTTP request with a context for cancellation.
//
// This method is useful when you need to cancel requests or set custom timeouts
// beyond the client's default timeout.
//
// Parameters:
//   - ctx: Context for cancellation and timeout
//   - req: The HTTP request to execute
//
// Returns:
//   - The HTTP response
//   - An error if the request fails or context is cancelled
//
// Example:
//
//	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
//	defer cancel()
//	req, _ := http.NewRequest("GET", "https://example.com", nil)
//	resp, err := client.DoWithContext(ctx, req)
func (c *HTTPClient) DoWithContext(ctx context.Context, req *http.Request) (*http.Response, error) {
	// Clone request with context
	req = req.WithContext(ctx)

	// Add User-Agent header if configured
	if c.userAgent != "" && req.Header.Get("User-Agent") == "" {
		req.Header.Set("User-Agent", c.userAgent)
	}

	// Execute the request
	return c.client.Do(req)
}

// Client returns the underlying http.Client.
//
// This is useful when you need to access the raw client for advanced use cases
// or to integrate with libraries that expect an *http.Client.
//
// Returns:
//   - The underlying *http.Client
//
// Example:
//
//	httpClient := client.Client()
//	// Use with other libraries that need *http.Client
func (c *HTTPClient) Client() *http.Client {
	return c.client
}

// Transport returns the underlying http.Transport.
//
// This is useful when you need to access or modify transport settings,
// or to integrate with libraries like colly that need an http.RoundTripper.
//
// Returns:
//   - The underlying *http.Transport
//
// Example:
//
//	transport := client.Transport()
//	// Use with colly or other libraries
func (c *HTTPClient) Transport() *http.Transport {
	return c.client.Transport.(*http.Transport)
}

// SetUserAgent updates the User-Agent header for future requests.
//
// Parameters:
//   - userAgent: The new User-Agent string
//
// Example:
//
//	client.SetUserAgent("MyApp/2.0")
func (c *HTTPClient) SetUserAgent(userAgent string) {
	c.userAgent = userAgent
}

// Close closes all idle connections in the connection pool.
//
// This should be called when the client is no longer needed to free resources.
// After calling Close, the client should not be used for new requests.
//
// Example:
//
//	defer client.Close()
func (c *HTTPClient) Close() {
	c.client.CloseIdleConnections()
}
