package dlsite

import (
	"net/http"
	"time"
)

// Options contains configuration options for the DLSite provider.
// These options allow customization of provider behavior without
// modifying the core implementation.
type Options struct {
	// Timeout is the overall timeout for HTTP requests
	Timeout time.Duration

	// RequestTimeout is the timeout for individual HTTP requests
	RequestTimeout time.Duration

	// RateLimit is the maximum number of requests per second (0 = no limit)
	RateLimit int

	// MaxConcurrentRequests is the maximum number of concurrent requests (0 = no limit)
	MaxConcurrentRequests int

	// CacheEnabled enables metadata caching
	CacheEnabled bool

	// CacheTTL is the time-to-live for cached entries
	CacheTTL time.Duration

	// CacheMaxSize is the maximum number of entries in the cache
	CacheMaxSize int

	// RetryEnabled enables automatic retry logic for failed requests
	RetryEnabled bool

	// MaxRetries is the maximum number of retry attempts
	MaxRetries int

	// BackoffMultiplier is the multiplier for exponential backoff
	BackoffMultiplier float64

	// MaxBackoff is the maximum backoff duration
	MaxBackoff time.Duration

	// HTTPClient is a custom HTTP client (optional)
	HTTPClient *http.Client

	// UserAgent is the User-Agent header for HTTP requests
	UserAgent string

	// Logger is a custom logger implementation (optional)
	Logger Logger

	// LogLevel is the minimum log level to output
	LogLevel LogLevel
}

// Option is a function that configures the DLSite provider.
// This follows the functional options pattern for flexible configuration.
type Option func(*DLSite)

// defaultOptions returns the default configuration options.
func defaultOptions() *Options {
	return &Options{
		Timeout:               30 * time.Second,
		RequestTimeout:        30 * time.Second,
		RateLimit:             0, // No rate limit by default
		MaxConcurrentRequests: 0, // No concurrency limit by default
		CacheEnabled:          false,
		CacheTTL:              1 * time.Hour,
		CacheMaxSize:          1000,
		RetryEnabled:          false,
		MaxRetries:            3,
		BackoffMultiplier:     2.0,
		MaxBackoff:            30 * time.Second,
		HTTPClient:            nil, // Use default HTTP client
		UserAgent:             "",  // Use default User-Agent
		Logger:                nil, // No logging by default
		LogLevel:              LogLevelInfo,
	}
}

// WithTimeout sets the overall timeout for HTTP requests.
//
// Example:
//
//	provider := dlsite.New(dlsite.WithTimeout(60 * time.Second))
func WithTimeout(timeout time.Duration) Option {
	return func(d *DLSite) {
		if d.options == nil {
			d.options = defaultOptions()
		}
		d.options.Timeout = timeout
	}
}

// WithRequestTimeout sets the timeout for individual HTTP requests.
//
// Example:
//
//	provider := dlsite.New(dlsite.WithRequestTimeout(30 * time.Second))
func WithRequestTimeout(timeout time.Duration) Option {
	return func(d *DLSite) {
		if d.options == nil {
			d.options = defaultOptions()
		}
		d.options.RequestTimeout = timeout
	}
}

// WithRateLimit sets the maximum number of requests per second.
// Set to 0 to disable rate limiting.
//
// Example:
//
//	provider := dlsite.New(dlsite.WithRateLimit(10)) // 10 requests per second
func WithRateLimit(requestsPerSecond int) Option {
	return func(d *DLSite) {
		if d.options == nil {
			d.options = defaultOptions()
		}
		d.options.RateLimit = requestsPerSecond
	}
}

// WithMaxConcurrentRequests sets the maximum number of concurrent requests.
// Set to 0 to disable concurrency limiting.
//
// Example:
//
//	provider := dlsite.New(dlsite.WithMaxConcurrentRequests(5))
func WithMaxConcurrentRequests(max int) Option {
	return func(d *DLSite) {
		if d.options == nil {
			d.options = defaultOptions()
		}
		d.options.MaxConcurrentRequests = max
	}
}

// WithCache enables metadata caching with the specified TTL and maximum size.
//
// Example:
//
//	provider := dlsite.New(dlsite.WithCache(1*time.Hour, 1000))
func WithCache(ttl time.Duration, maxSize int) Option {
	return func(d *DLSite) {
		if d.options == nil {
			d.options = defaultOptions()
		}
		d.options.CacheEnabled = true
		d.options.CacheTTL = ttl
		d.options.CacheMaxSize = maxSize
	}
}

// WithCacheDisabled explicitly disables metadata caching.
//
// Example:
//
//	provider := dlsite.New(dlsite.WithCacheDisabled())
func WithCacheDisabled() Option {
	return func(d *DLSite) {
		if d.options == nil {
			d.options = defaultOptions()
		}
		d.options.CacheEnabled = false
	}
}

// WithRetryPolicy enables automatic retry logic with the specified parameters.
//
// Example:
//
//	provider := dlsite.New(dlsite.WithRetryPolicy(3, 2.0))
func WithRetryPolicy(maxRetries int, backoffMultiplier float64) Option {
	return func(d *DLSite) {
		if d.options == nil {
			d.options = defaultOptions()
		}
		d.options.RetryEnabled = true
		d.options.MaxRetries = maxRetries
		d.options.BackoffMultiplier = backoffMultiplier
	}
}

// WithRetryDisabled explicitly disables automatic retry logic.
//
// Example:
//
//	provider := dlsite.New(dlsite.WithRetryDisabled())
func WithRetryDisabled() Option {
	return func(d *DLSite) {
		if d.options == nil {
			d.options = defaultOptions()
		}
		d.options.RetryEnabled = false
	}
}

// WithHTTPClient sets a custom HTTP client for the provider.
//
// Example:
//
//	client := &http.Client{Timeout: 60 * time.Second}
//	provider := dlsite.New(dlsite.WithHTTPClient(client))
func WithHTTPClient(client *http.Client) Option {
	return func(d *DLSite) {
		if d.options == nil {
			d.options = defaultOptions()
		}
		d.options.HTTPClient = client
	}
}

// WithUserAgent sets a custom User-Agent header for HTTP requests.
//
// Example:
//
//	provider := dlsite.New(dlsite.WithUserAgent("MyApp/1.0"))
func WithUserAgent(userAgent string) Option {
	return func(d *DLSite) {
		if d.options == nil {
			d.options = defaultOptions()
		}
		d.options.UserAgent = userAgent
	}
}

// WithLogger sets a custom logger for the provider.
// If no logger is set, logging is disabled by default.
//
// Example:
//
//	logger := dlsite.NewDefaultLogger(dlsite.LogLevelInfo)
//	provider := dlsite.New(dlsite.WithLogger(logger))
func WithLogger(logger Logger) Option {
	return func(d *DLSite) {
		if d.options == nil {
			d.options = defaultOptions()
		}
		d.options.Logger = logger
	}
}

// WithLogLevel sets the minimum log level for the default logger.
// This option only affects the DefaultLogger. Custom loggers should
// implement their own level filtering.
//
// Example:
//
//	provider := dlsite.New(
//	    dlsite.WithLogger(dlsite.NewDefaultLogger(dlsite.LogLevelDebug)),
//	    dlsite.WithLogLevel(dlsite.LogLevelDebug),
//	)
func WithLogLevel(level LogLevel) Option {
	return func(d *DLSite) {
		if d.options == nil {
			d.options = defaultOptions()
		}
		d.options.LogLevel = level
	}
}
