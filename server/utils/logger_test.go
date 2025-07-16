package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"
	"time"
)

func TestLogLevel_String(t *testing.T) {
	tests := []struct {
		level    LogLevel
		expected string
	}{
		{DebugLevel, "DEBUG"},
		{InfoLevel, "INFO"},
		{WarnLevel, "WARN"},
		{ErrorLevel, "ERROR"},
		{FatalLevel, "FATAL"},
		{LogLevel(99), "UNKNOWN"},
	}

	for _, test := range tests {
		if got := test.level.String(); got != test.expected {
			t.Errorf("LogLevel.String() = %v, want %v", got, test.expected)
		}
	}
}

func TestNewJSONLogger(t *testing.T) {
	var buf bytes.Buffer
	logger := NewJSONLogger(&buf, InfoLevel, "test-component")

	if logger == nil {
		t.Fatal("NewJSONLogger returned nil")
	}

	jsonLogger, ok := logger.(*JSONLogger)
	if !ok {
		t.Fatal("NewJSONLogger did not return *JSONLogger")
	}

	if jsonLogger.level != InfoLevel {
		t.Errorf("Expected level %v, got %v", InfoLevel, jsonLogger.level)
	}

	if jsonLogger.component != "test-component" {
		t.Errorf("Expected component 'test-component', got '%v'", jsonLogger.component)
	}
}

func TestNewDefaultLogger(t *testing.T) {
	logger := NewDefaultLogger("default-test")
	
	if logger == nil {
		t.Fatal("NewDefaultLogger returned nil")
	}

	if logger.GetLevel() != InfoLevel {
		t.Errorf("Expected default level %v, got %v", InfoLevel, logger.GetLevel())
	}
}

func TestJSONLogger_SetLevel(t *testing.T) {
	var buf bytes.Buffer
	logger := NewJSONLogger(&buf, InfoLevel, "test")

	logger.SetLevel(ErrorLevel)
	if logger.GetLevel() != ErrorLevel {
		t.Errorf("Expected level %v, got %v", ErrorLevel, logger.GetLevel())
	}
}

func TestJSONLogger_Info(t *testing.T) {
	var buf bytes.Buffer
	logger := NewJSONLogger(&buf, InfoLevel, "test-component")

	ctx := context.WithValue(context.Background(), contextKey("request_id"), "req-123")
	ctx = context.WithValue(ctx, contextKey("user_id"), "user-456")

	logger.Info(ctx, "test message", StringField("key1", "value1"))

	var entry LogEntry
	if err := json.Unmarshal(buf.Bytes(), &entry); err != nil {
		t.Fatalf("Failed to unmarshal log entry: %v", err)
	}

	if entry.Level != "INFO" {
		t.Errorf("Expected level INFO, got %v", entry.Level)
	}

	if entry.Message != "test message" {
		t.Errorf("Expected message 'test message', got '%v'", entry.Message)
	}

	if entry.RequestID != "req-123" {
		t.Errorf("Expected request ID 'req-123', got '%v'", entry.RequestID)
	}

	if entry.UserID != "user-456" {
		t.Errorf("Expected user ID 'user-456', got '%v'", entry.UserID)
	}

	if entry.Fields["component"] != "test-component" {
		t.Errorf("Expected component 'test-component', got '%v'", entry.Fields["component"])
	}

	if entry.Fields["key1"] != "value1" {
		t.Errorf("Expected field key1 'value1', got '%v'", entry.Fields["key1"])
	}
}

func TestJSONLogger_Error(t *testing.T) {
	var buf bytes.Buffer
	logger := NewJSONLogger(&buf, InfoLevel, "test-component")

	ctx := context.WithValue(context.Background(), contextKey("request_id"), "req-123")
	testErr := errors.New("test error")

	logger.Error(ctx, "error occurred", testErr, StringField("operation", "test-op"))

	var entry LogEntry
	if err := json.Unmarshal(buf.Bytes(), &entry); err != nil {
		t.Fatalf("Failed to unmarshal log entry: %v", err)
	}

	if entry.Level != "ERROR" {
		t.Errorf("Expected level ERROR, got %v", entry.Level)
	}

	if entry.Error != "test error" {
		t.Errorf("Expected error 'test error', got '%v'", entry.Error)
	}

	if entry.Caller == "" {
		t.Error("Expected caller information for error level")
	}

	if entry.Fields["operation"] != "test-op" {
		t.Errorf("Expected operation 'test-op', got '%v'", entry.Fields["operation"])
	}
}

func TestJSONLogger_LevelFiltering(t *testing.T) {
	var buf bytes.Buffer
	logger := NewJSONLogger(&buf, ErrorLevel, "test")

	ctx := context.Background()

	// These should not be logged due to level filtering
	logger.Debug(ctx, "debug message")
	logger.Info(ctx, "info message")
	logger.Warn(ctx, "warn message")

	if buf.Len() > 0 {
		t.Error("Expected no output for levels below ErrorLevel")
	}

	// This should be logged
	logger.Error(ctx, "error message", nil)

	if buf.Len() == 0 {
		t.Error("Expected output for ErrorLevel")
	}
}

func TestJSONLogger_WithoutContext(t *testing.T) {
	var buf bytes.Buffer
	logger := NewJSONLogger(&buf, InfoLevel, "test")

	logger.Info(nil, "message without context")

	var entry LogEntry
	if err := json.Unmarshal(buf.Bytes(), &entry); err != nil {
		t.Fatalf("Failed to unmarshal log entry: %v", err)
	}

	if entry.RequestID != "" {
		t.Errorf("Expected empty request ID, got '%v'", entry.RequestID)
	}

	if entry.UserID != "" {
		t.Errorf("Expected empty user ID, got '%v'", entry.UserID)
	}
}

func TestGetRequestIDFromContext(t *testing.T) {
	tests := []struct {
		name     string
		ctx      context.Context
		expected string
	}{
		{
			name:     "with request ID",
			ctx:      context.WithValue(context.Background(), contextKey("request_id"), "req-123"),
			expected: "req-123",
		},
		{
			name:     "without request ID",
			ctx:      context.Background(),
			expected: "",
		},
		{
			name:     "nil context",
			ctx:      nil,
			expected: "",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := GetRequestIDFromContext(test.ctx)
			if got != test.expected {
				t.Errorf("GetRequestIDFromContext() = %v, want %v", got, test.expected)
			}
		})
	}
}

func TestGetUserIDFromContext(t *testing.T) {
	tests := []struct {
		name     string
		ctx      context.Context
		expected string
	}{
		{
			name:     "with user ID",
			ctx:      context.WithValue(context.Background(), contextKey("user_id"), "user-456"),
			expected: "user-456",
		},
		{
			name:     "without user ID",
			ctx:      context.Background(),
			expected: "",
		},
		{
			name:     "nil context",
			ctx:      nil,
			expected: "",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := GetUserIDFromContext(test.ctx)
			if got != test.expected {
				t.Errorf("GetUserIDFromContext() = %v, want %v", got, test.expected)
			}
		})
	}
}

func TestFieldHelpers(t *testing.T) {
	tests := []struct {
		name     string
		field    Field
		expected interface{}
	}{
		{"StringField", StringField("key", "value"), "value"},
		{"IntField", IntField("key", 42), 42},
		{"Int64Field", Int64Field("key", int64(123)), int64(123)},
		{"BoolField", BoolField("key", true), true},
		{"DurationField", DurationField("key", time.Second), "1s"},
		{"ErrorField", ErrorField(errors.New("test error")), "test error"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.field.Value != test.expected {
				t.Errorf("Field value = %v, want %v", test.field.Value, test.expected)
			}
		})
	}
}

func TestLoggerFactory(t *testing.T) {
	factory := NewLoggerFactory(WarnLevel)

	logger1 := factory.GetLogger("component1")
	logger2 := factory.GetLogger("component2")

	if logger1.GetLevel() != WarnLevel {
		t.Errorf("Expected level %v, got %v", WarnLevel, logger1.GetLevel())
	}

	if logger2.GetLevel() != WarnLevel {
		t.Errorf("Expected level %v, got %v", WarnLevel, logger2.GetLevel())
	}

	// Test custom writer
	var buf bytes.Buffer
	factory.SetWriter("component1", &buf)
	
	logger3 := factory.GetLogger("component1")
	logger3.Warn(context.Background(), "test message")

	if buf.Len() == 0 {
		t.Error("Expected output to custom writer")
	}
}

func TestJSONLogger_MarshalError(t *testing.T) {
	var buf bytes.Buffer
	logger := NewJSONLogger(&buf, InfoLevel, "test")

	// Create a field that will cause JSON marshal error
	ctx := context.Background()
	
	// Use a channel which cannot be marshaled to JSON
	ch := make(chan int)
	logger.Info(ctx, "test message", Field{Key: "channel", Value: ch})

	// Should fallback to simple text logging
	output := buf.String()
	if !strings.Contains(output, "test message") {
		t.Error("Expected fallback message in output")
	}
	
	// Should not be valid JSON due to fallback
	var entry LogEntry
	if err := json.Unmarshal(buf.Bytes(), &entry); err == nil {
		t.Error("Expected JSON unmarshal to fail due to fallback logging")
	}
}