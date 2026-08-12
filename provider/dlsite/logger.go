package dlsite

// Logger is an interface for structured logging.
// It allows custom logging implementations to be used with the DLSite provider.
//
// The logger supports four log levels (Debug, Info, Warn, Error) and accepts
// structured fields for additional context.
//
// Example implementation:
//
//	type MyLogger struct{}
//
//	func (l *MyLogger) Debug(msg string, fields ...Field) {
//	    log.Printf("[DEBUG] %s %v", msg, fields)
//	}
//
// Example usage:
//
//	provider := dlsite.New(
//	    dlsite.WithLogger(myLogger),
//	    dlsite.WithLogLevel(LogLevelInfo),
//	)
type Logger interface {
	// Debug logs a debug-level message with optional structured fields.
	// Debug messages are typically used for detailed diagnostic information.
	Debug(msg string, fields ...Field)

	// Info logs an info-level message with optional structured fields.
	// Info messages are used for general informational messages.
	Info(msg string, fields ...Field)

	// Warn logs a warning-level message with optional structured fields.
	// Warn messages indicate potential issues that don't prevent operation.
	Warn(msg string, fields ...Field)

	// Error logs an error-level message with optional structured fields.
	// Error messages indicate failures that prevent normal operation.
	Error(msg string, fields ...Field)
}

// Field represents a structured logging field with a key-value pair.
// Fields provide additional context for log messages.
//
// Example:
//
//	logger.Info("request completed",
//	    Field{Key: "url", Value: "https://example.com"},
//	    Field{Key: "duration", Value: time.Second},
//	)
type Field struct {
	// Key is the field name
	Key string

	// Value is the field value (can be any type)
	Value interface{}
}

// LogLevel represents the severity level for logging.
// Log levels are ordered from most verbose (Debug) to least verbose (Error).
type LogLevel int

const (
	// LogLevelDebug enables all log messages including detailed diagnostics
	LogLevelDebug LogLevel = iota

	// LogLevelInfo enables informational, warning, and error messages
	LogLevelInfo

	// LogLevelWarn enables warning and error messages only
	LogLevelWarn

	// LogLevelError enables error messages only
	LogLevelError
)

// String returns the string representation of the log level.
func (l LogLevel) String() string {
	switch l {
	case LogLevelDebug:
		return "DEBUG"
	case LogLevelInfo:
		return "INFO"
	case LogLevelWarn:
		return "WARN"
	case LogLevelError:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}

// DefaultLogger is a simple logger implementation using the standard log package.
// It supports log level filtering and automatic sanitization of sensitive information.
//
// Sensitive information (cookies, tokens, passwords) is automatically redacted
// from log output to prevent accidental exposure.
//
// Example:
//
//	logger := NewDefaultLogger(LogLevelInfo)
//	logger.Info("request completed", Field{Key: "url", Value: "https://example.com"})
type DefaultLogger struct {
	level LogLevel
}

// NewDefaultLogger creates a new DefaultLogger with the specified log level.
//
// Parameters:
//   - level: The minimum log level to output (messages below this level are filtered)
//
// Returns:
//   - A configured DefaultLogger instance
//
// Example:
//
//	logger := NewDefaultLogger(LogLevelInfo)
//	provider := dlsite.New(dlsite.WithLogger(logger))
func NewDefaultLogger(level LogLevel) *DefaultLogger {
	return &DefaultLogger{
		level: level,
	}
}

// Debug logs a debug-level message if the logger's level is Debug or lower.
func (l *DefaultLogger) Debug(msg string, fields ...Field) {
	if l.level <= LogLevelDebug {
		l.log("DEBUG", msg, fields)
	}
}

// Info logs an info-level message if the logger's level is Info or lower.
func (l *DefaultLogger) Info(msg string, fields ...Field) {
	if l.level <= LogLevelInfo {
		l.log("INFO", msg, fields)
	}
}

// Warn logs a warning-level message if the logger's level is Warn or lower.
func (l *DefaultLogger) Warn(msg string, fields ...Field) {
	if l.level <= LogLevelWarn {
		l.log("WARN", msg, fields)
	}
}

// Error logs an error-level message if the logger's level is Error or lower.
func (l *DefaultLogger) Error(msg string, fields ...Field) {
	if l.level <= LogLevelError {
		l.log("ERROR", msg, fields)
	}
}

// log is the internal logging method that formats and outputs log messages.
// It automatically sanitizes sensitive information from fields.
func (l *DefaultLogger) log(level, msg string, fields []Field) {
	// Sanitize sensitive information
	sanitized := l.sanitizeFields(fields)

	// Format fields as key=value pairs
	var fieldStr string
	if len(sanitized) > 0 {
		fieldStr = " "
		for i, field := range sanitized {
			if i > 0 {
				fieldStr += " "
			}
			fieldStr += field.Key + "=" + formatValue(field.Value)
		}
	}

	// Output log message
	println("[" + level + "] " + msg + fieldStr)
}

// sanitizeFields removes or redacts sensitive information from log fields.
// This prevents accidental exposure of cookies, tokens, passwords, etc.
func (l *DefaultLogger) sanitizeFields(fields []Field) []Field {
	if len(fields) == 0 {
		return fields
	}

	sanitized := make([]Field, 0, len(fields))
	for _, field := range fields {
		// Check if field contains sensitive information
		if isSensitiveKey(field.Key) {
			sanitized = append(sanitized, Field{
				Key:   field.Key,
				Value: "[REDACTED]",
			})
		} else {
			sanitized = append(sanitized, field)
		}
	}
	return sanitized
}

// isSensitiveKey checks if a field key contains sensitive information.
// Returns true if the key contains words like "cookie", "token", "password", etc.
func isSensitiveKey(key string) bool {
	// Convert to lowercase for case-insensitive matching
	lowerKey := ""
	for _, c := range key {
		if c >= 'A' && c <= 'Z' {
			lowerKey += string(c + 32)
		} else {
			lowerKey += string(c)
		}
	}

	// Check for sensitive keywords
	sensitiveKeywords := []string{
		"cookie", "token", "password", "secret", "key",
		"auth", "authorization", "credential", "api_key",
	}

	for _, keyword := range sensitiveKeywords {
		if contains(lowerKey, keyword) {
			return true
		}
	}

	return false
}

// contains checks if a string contains a substring (simple implementation).
func contains(s, substr string) bool {
	if len(substr) == 0 {
		return true
	}
	if len(s) < len(substr) {
		return false
	}
	for i := 0; i <= len(s)-len(substr); i++ {
		match := true
		for j := 0; j < len(substr); j++ {
			if s[i+j] != substr[j] {
				match = false
				break
			}
		}
		if match {
			return true
		}
	}
	return false
}

// formatValue converts a field value to a string representation.
func formatValue(v interface{}) string {
	if v == nil {
		return "<nil>"
	}

	switch val := v.(type) {
	case string:
		return val
	case int, int8, int16, int32, int64:
		return intToString(val)
	case uint, uint8, uint16, uint32, uint64:
		return uintToString(val)
	case float32, float64:
		return floatToString(val)
	case bool:
		if val {
			return "true"
		}
		return "false"
	default:
		// For other types, use a simple representation
		return "<value>"
	}
}

// intToString converts an integer to a string.
func intToString(v interface{}) string {
	var n int64
	switch val := v.(type) {
	case int:
		n = int64(val)
	case int8:
		n = int64(val)
	case int16:
		n = int64(val)
	case int32:
		n = int64(val)
	case int64:
		n = val
	default:
		return "<int>"
	}

	if n == 0 {
		return "0"
	}

	negative := n < 0
	if negative {
		n = -n
	}

	// Convert to string
	digits := make([]byte, 0, 20)
	for n > 0 {
		digits = append(digits, byte('0'+n%10))
		n /= 10
	}

	// Reverse digits
	for i, j := 0, len(digits)-1; i < j; i, j = i+1, j-1 {
		digits[i], digits[j] = digits[j], digits[i]
	}

	if negative {
		return "-" + string(digits)
	}
	return string(digits)
}

// uintToString converts an unsigned integer to a string.
func uintToString(v interface{}) string {
	var n uint64
	switch val := v.(type) {
	case uint:
		n = uint64(val)
	case uint8:
		n = uint64(val)
	case uint16:
		n = uint64(val)
	case uint32:
		n = uint64(val)
	case uint64:
		n = val
	default:
		return "<uint>"
	}

	if n == 0 {
		return "0"
	}

	// Convert to string
	digits := make([]byte, 0, 20)
	for n > 0 {
		digits = append(digits, byte('0'+n%10))
		n /= 10
	}

	// Reverse digits
	for i, j := 0, len(digits)-1; i < j; i, j = i+1, j-1 {
		digits[i], digits[j] = digits[j], digits[i]
	}

	return string(digits)
}

// floatToString converts a float to a string (simplified).
func floatToString(v interface{}) string {
	// Simplified float formatting - just return a placeholder
	// In a real implementation, you'd use strconv.FormatFloat
	return "<float>"
}
