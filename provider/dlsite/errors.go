package dlsite

import (
	"errors"
	"fmt"
	"time"
)

// Existing error types
var (
	// ErrInvalidWorkID indicates that the provided work ID format is invalid.
	ErrInvalidWorkID = errors.New("invalid work ID format")

	// ErrUnsupportedContentType indicates that the content type is not supported.
	// This error is returned for non-RJ prefixed IDs (VJ, BJ, etc.).
	ErrUnsupportedContentType = errors.New("unsupported content type")

	// ErrWorkNotFound indicates that the requested work does not exist (HTTP 404).
	ErrWorkNotFound = errors.New("work not found")

	// ErrParseError indicates that HTML parsing failed.
	ErrParseError = errors.New("failed to parse HTML")

	// ErrNetworkError indicates that a network request failed.
	ErrNetworkError = errors.New("network request failed")

	// ErrValidationError indicates that metadata validation failed.
	ErrValidationError = errors.New("metadata validation failed")
)

// New error types
var (
	// ErrRateLimitExceeded indicates that the rate limit has been exceeded.
	ErrRateLimitExceeded = errors.New("rate limit exceeded")

	// ErrTimeout indicates that a request timed out.
	ErrTimeout = errors.New("request timeout")

	// ErrRetryExhausted indicates that the maximum number of retries has been exceeded.
	ErrRetryExhausted = errors.New("maximum retries exceeded")

	// ErrCacheFull indicates that the cache is full and cannot accept new entries.
	ErrCacheFull = errors.New("cache is full")

	// ErrInvalidOption indicates that an invalid option was provided.
	ErrInvalidOption = errors.New("invalid option")
)

// DetailedError provides detailed context information for errors.
// It wraps an underlying error and adds contextual information such as
// URL, status code, and custom context fields.
type DetailedError struct {
	// Err is the underlying error
	Err error

	// Context contains additional contextual information about the error
	Context map[string]interface{}

	// StatusCode is the HTTP status code (if applicable)
	StatusCode int

	// URL is the URL that was being accessed when the error occurred
	URL string

	// Timestamp is when the error occurred
	Timestamp time.Time
}

// Error implements the error interface.
// It returns a formatted error message that includes the underlying error
// and contextual information.
func (e *DetailedError) Error() string {
	if len(e.Context) > 0 {
		return fmt.Sprintf("%v (context: %v)", e.Err, e.Context)
	}
	return e.Err.Error()
}

// Unwrap implements error unwrapping for error chaining.
// This allows errors.Is() and errors.As() to work correctly.
func (e *DetailedError) Unwrap() error {
	return e.Err
}

// NewDetailedError creates a new DetailedError with the given error and context.
func NewDetailedError(err error, context map[string]interface{}) *DetailedError {
	return &DetailedError{
		Err:       err,
		Context:   context,
		Timestamp: time.Now(),
	}
}

// NewDetailedErrorWithStatus creates a new DetailedError with error, context, and HTTP status.
func NewDetailedErrorWithStatus(err error, statusCode int, url string, context map[string]interface{}) *DetailedError {
	return &DetailedError{
		Err:        err,
		Context:    context,
		StatusCode: statusCode,
		URL:        url,
		Timestamp:  time.Now(),
	}
}
