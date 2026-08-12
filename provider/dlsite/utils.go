package dlsite

import (
	"bytes"
	"strings"
	"sync"
)

// bufferPool is a sync.Pool for reusing byte buffers.
// This reduces memory allocations for temporary buffer operations.
var bufferPool = sync.Pool{
	New: func() interface{} {
		return new(bytes.Buffer)
	},
}

// stringBuilderPool is a sync.Pool for reusing strings.Builder instances.
// This reduces memory allocations for string concatenation operations.
var stringBuilderPool = sync.Pool{
	New: func() interface{} {
		return new(strings.Builder)
	},
}

// getBuffer retrieves a buffer from the pool.
// The caller must call putBuffer when done to return it to the pool.
func getBuffer() *bytes.Buffer {
	buf := bufferPool.Get().(*bytes.Buffer)
	buf.Reset()
	return buf
}

// putBuffer returns a buffer to the pool for reuse.
// The buffer is reset before being returned to the pool.
func putBuffer(buf *bytes.Buffer) {
	buf.Reset()
	bufferPool.Put(buf)
}

// getStringBuilder retrieves a strings.Builder from the pool.
// The caller must call putStringBuilder when done to return it to the pool.
func getStringBuilder() *strings.Builder {
	sb := stringBuilderPool.Get().(*strings.Builder)
	sb.Reset()
	return sb
}

// putStringBuilder returns a strings.Builder to the pool for reuse.
// The builder is reset before being returned to the pool.
func putStringBuilder(sb *strings.Builder) {
	sb.Reset()
	stringBuilderPool.Put(sb)
}

// buildURL constructs a URL from parts using a pooled strings.Builder.
// This is more efficient than string concatenation or fmt.Sprintf for URL building.
//
// Example:
//   url := buildURL(baseURL, "/path/", id, ".html")
func buildURL(parts ...string) string {
	sb := getStringBuilder()
	defer putStringBuilder(sb)
	
	// Calculate total length for efficient allocation
	totalLen := 0
	for _, part := range parts {
		totalLen += len(part)
	}
	sb.Grow(totalLen)
	
	// Write all parts
	for _, part := range parts {
		sb.WriteString(part)
	}
	
	return sb.String()
}

// normalizeImageURL adds the https: prefix to URLs starting with //.
// Uses a pooled strings.Builder for efficient string construction.
func normalizeImageURL(url string) string {
	if strings.HasPrefix(url, "//") {
		sb := getStringBuilder()
		defer putStringBuilder(sb)
		
		sb.Grow(len("https:") + len(url))
		sb.WriteString("https:")
		sb.WriteString(url)
		return sb.String()
	}
	return url
}
