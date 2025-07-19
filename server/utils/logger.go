package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
	"time"
)

// LogLevel represents the severity level of a log entry
type LogLevel int

const (
	DebugLevel LogLevel = iota
	InfoLevel
	WarnLevel
	ErrorLevel
	FatalLevel
)

// String returns the string representation of the log level
func (l LogLevel) String() string {
	switch l {
	case DebugLevel:
		return "DEBUG"
	case InfoLevel:
		return "INFO"
	case WarnLevel:
		return "WARN"
	case ErrorLevel:
		return "ERROR"
	case FatalLevel:
		return "FATAL"
	default:
		return "UNKNOWN"
	}
}

// Field represents a structured log field
type Field struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}

// LogEntry represents a structured log entry with context information
type LogEntry struct {
	Level     string                 `json:"level"`
	Message   string                 `json:"message"`
	Timestamp time.Time              `json:"timestamp"`
	RequestID string                 `json:"request_id,omitempty"`
	UserID    string                 `json:"user_id,omitempty"`
	Fields    map[string]interface{} `json:"fields,omitempty"`
	Error     string                 `json:"error,omitempty"`
	Caller    string                 `json:"caller,omitempty"`
}

// Logger interface defines the contract for context-aware logging
type Logger interface {
	Info(ctx context.Context, msg string, fields ...Field)
	Warn(ctx context.Context, msg string, fields ...Field)
	Error(ctx context.Context, msg string, err error, fields ...Field)
	Debug(ctx context.Context, msg string, fields ...Field)
	Fatal(ctx context.Context, msg string, err error, fields ...Field)
	WithContext(ctx context.Context) Logger
	SetLevel(level LogLevel)
	GetLevel() LogLevel
}

// JSONLogger implements the Logger interface with JSON structured logging
type JSONLogger struct {
	writer    io.Writer
	level     LogLevel
	component string
}

// NewJSONLogger creates a new JSON logger instance
func NewJSONLogger(writer io.Writer, level LogLevel, component string) Logger {
	if writer == nil {
		writer = os.Stdout
	}
	return &JSONLogger{
		writer:    writer,
		level:     level,
		component: component,
	}
}

// NewDefaultLogger creates a default JSON logger with INFO level
func NewDefaultLogger(component string) Logger {
	return NewJSONLogger(os.Stdout, InfoLevel, component)
}

// SetLevel sets the minimum log level
func (l *JSONLogger) SetLevel(level LogLevel) {
	l.level = level
}

// GetLevel returns the current log level
func (l *JSONLogger) GetLevel() LogLevel {
	return l.level
}

// WithContext returns a logger with context information
func (l *JSONLogger) WithContext(ctx context.Context) Logger {
	// Return the same logger as context is handled per log call
	return l
}

// Info logs an info level message with context
func (l *JSONLogger) Info(ctx context.Context, msg string, fields ...Field) {
	if l.level > InfoLevel {
		return
	}
	l.log(ctx, InfoLevel, msg, nil, fields...)
}

// Warn logs a warning level message with context
func (l *JSONLogger) Warn(ctx context.Context, msg string, fields ...Field) {
	if l.level > WarnLevel {
		return
	}
	l.log(ctx, WarnLevel, msg, nil, fields...)
}

// Error logs an error level message with context
func (l *JSONLogger) Error(ctx context.Context, msg string, err error, fields ...Field) {
	if l.level > ErrorLevel {
		return
	}
	l.log(ctx, ErrorLevel, msg, err, fields...)
}

// Debug logs a debug level message with context
func (l *JSONLogger) Debug(ctx context.Context, msg string, fields ...Field) {
	if l.level > DebugLevel {
		return
	}
	l.log(ctx, DebugLevel, msg, nil, fields...)
}

// Fatal logs a fatal level message with context and exits
func (l *JSONLogger) Fatal(ctx context.Context, msg string, err error, fields ...Field) {
	l.log(ctx, FatalLevel, msg, err, fields...)
	os.Exit(1)
}

// log is the internal logging method that handles context extraction and formatting
func (l *JSONLogger) log(ctx context.Context, level LogLevel, msg string, err error, fields ...Field) {
	entry := LogEntry{
		Level:     level.String(),
		Message:   msg,
		Timestamp: time.Now().UTC(),
		Fields:    make(map[string]interface{}),
	}

	// Extract context information if available
	if ctx != nil {
		if requestID := GetRequestIDFromContext(ctx); requestID != "" {
			entry.RequestID = requestID
		}
		if userID := GetUserIDFromContext(ctx); userID != "" {
			entry.UserID = userID
		}
	}

	// Add component information
	if l.component != "" {
		entry.Fields["component"] = l.component
	}

	// Add custom fields
	for _, field := range fields {
		entry.Fields[field.Key] = field.Value
	}

	// Add error information if present
	if err != nil {
		entry.Error = err.Error()
	}

	// Add caller information for error and fatal levels
	if level >= ErrorLevel {
		if caller := getCaller(3); caller != "" {
			entry.Caller = caller
		}
	}

	// Marshal and write the log entry
	jsonData, marshalErr := json.Marshal(entry)
	if marshalErr != nil {
		// Fallback to simple text logging if JSON marshaling fails
		fallbackMsg := fmt.Sprintf("[%s] %s %s - %s\n", 
			entry.Timestamp.Format(time.RFC3339), 
			entry.Level, 
			entry.RequestID, 
			msg)
		l.writer.Write([]byte(fallbackMsg))
		return
	}

	// Write JSON log entry
	l.writer.Write(jsonData)
	l.writer.Write([]byte("\n"))
}

// getCaller returns the caller information (file:line)
func getCaller(skip int) string {
	_, file, line, ok := runtime.Caller(skip)
	if !ok {
		return ""
	}
	
	// Get just the filename, not the full path
	parts := strings.Split(file, "/")
	if len(parts) > 0 {
		file = parts[len(parts)-1]
	}
	
	return fmt.Sprintf("%s:%d", file, line)
}

// Context helper functions
func GetRequestIDFromContext(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if requestID, ok := ctx.Value(contextKey("request_id")).(string); ok {
		return requestID
	}
	return ""
}

func GetUserIDFromContext(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if userID, ok := ctx.Value(contextKey("user_id")).(string); ok {
		return userID
	}
	return ""
}

// contextKey type for context keys (matching middleware package)
type contextKey string

// NewBackgroundContext creates a new background context for system operations
func NewBackgroundContext() context.Context {
	return context.Background()
}

// Helper functions to create common fields
func StringField(key, value string) Field {
	return Field{Key: key, Value: value}
}

func IntField(key string, value int) Field {
	return Field{Key: key, Value: value}
}

func Int64Field(key string, value int64) Field {
	return Field{Key: key, Value: value}
}

func DurationField(key string, value time.Duration) Field {
	return Field{Key: key, Value: value.String()}
}

func BoolField(key string, value bool) Field {
	return Field{Key: key, Value: value}
}

func ErrorField(err error) Field {
	return Field{Key: "error_details", Value: err.Error()}
}

// LoggerFactory manages logger instances for different components
type LoggerFactory struct {
	level   LogLevel
	writers map[string]io.Writer
}

// NewLoggerFactory creates a new logger factory
func NewLoggerFactory(level LogLevel) *LoggerFactory {
	return &LoggerFactory{
		level:   level,
		writers: make(map[string]io.Writer),
	}
}

// GetLogger returns a logger for the specified component
func (f *LoggerFactory) GetLogger(component string) Logger {
	writer, exists := f.writers[component]
	if !exists {
		writer = os.Stdout
	}
	return NewJSONLogger(writer, f.level, component)
}

// SetWriter sets a custom writer for a specific component
func (f *LoggerFactory) SetWriter(component string, writer io.Writer) {
	f.writers[component] = writer
}

// SetLevel sets the log level for all loggers created by this factory
func (f *LoggerFactory) SetLevel(level LogLevel) {
	f.level = level
}