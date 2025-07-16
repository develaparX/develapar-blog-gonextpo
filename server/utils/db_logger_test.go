package utils

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"strings"
	"testing"
	"time"

	_ "github.com/lib/pq" // PostgreSQL driver for testing
)

func TestNewDBLogger(t *testing.T) {
	var buf bytes.Buffer
	logger := NewJSONLogger(&buf, DebugLevel, "test")
	
	dbLogger := NewDBLogger(logger, 500*time.Millisecond)
	if dbLogger == nil {
		t.Fatal("NewDBLogger returned nil")
	}
}

func TestDBLogger_LogQuery(t *testing.T) {
	var buf bytes.Buffer
	logger := NewJSONLogger(&buf, DebugLevel, "test")
	dbLogger := NewDBLogger(logger, 500*time.Millisecond)

	ctx := context.WithValue(context.Background(), contextKey("request_id"), "req-123")
	query := "SELECT * FROM users WHERE id = $1"
	args := []interface{}{123}
	duration := 100 * time.Millisecond

	// Test successful query
	dbLogger.LogQuery(ctx, query, args, duration, nil)

	var entry LogEntry
	if err := json.Unmarshal(buf.Bytes(), &entry); err != nil {
		t.Fatalf("Failed to unmarshal log entry: %v", err)
	}

	if entry.Level != "DEBUG" {
		t.Errorf("Expected level DEBUG, got %s", entry.Level)
	}

	if entry.Message != "Database query executed" {
		t.Errorf("Expected message 'Database query executed', got '%s'", entry.Message)
	}

	if entry.Fields["query"] != query {
		t.Errorf("Expected query '%s', got '%v'", query, entry.Fields["query"])
	}

	if entry.Fields["duration"] != "100ms" {
		t.Errorf("Expected duration '100ms', got '%v'", entry.Fields["duration"])
	}

	if entry.Fields["arg_count"] != float64(1) {
		t.Errorf("Expected arg_count 1, got %v", entry.Fields["arg_count"])
	}
}

func TestDBLogger_LogQueryWithError(t *testing.T) {
	var buf bytes.Buffer
	logger := NewJSONLogger(&buf, DebugLevel, "test")
	dbLogger := NewDBLogger(logger, 500*time.Millisecond)

	ctx := context.Background()
	query := "SELECT * FROM users WHERE id = $1"
	args := []interface{}{123}
	duration := 100 * time.Millisecond
	testErr := errors.New("database connection failed")

	dbLogger.LogQuery(ctx, query, args, duration, testErr)

	var entry LogEntry
	if err := json.Unmarshal(buf.Bytes(), &entry); err != nil {
		t.Fatalf("Failed to unmarshal log entry: %v", err)
	}

	if entry.Level != "ERROR" {
		t.Errorf("Expected level ERROR, got %s", entry.Level)
	}

	if entry.Message != "Database query failed" {
		t.Errorf("Expected message 'Database query failed', got '%s'", entry.Message)
	}

	if entry.Error != "database connection failed" {
		t.Errorf("Expected error 'database connection failed', got '%s'", entry.Error)
	}
}

func TestDBLogger_LogSlowQuery(t *testing.T) {
	var buf bytes.Buffer
	logger := NewJSONLogger(&buf, DebugLevel, "test")
	dbLogger := NewDBLogger(logger, 500*time.Millisecond)

	ctx := context.Background()
	query := "SELECT * FROM users"
	args := []interface{}{}
	duration := 2 * time.Second // Slow query
	threshold := 500 * time.Millisecond

	dbLogger.LogSlowQuery(ctx, query, args, duration, threshold)

	var entry LogEntry
	if err := json.Unmarshal(buf.Bytes(), &entry); err != nil {
		t.Fatalf("Failed to unmarshal log entry: %v", err)
	}

	if entry.Level != "WARN" {
		t.Errorf("Expected level WARN, got %s", entry.Level)
	}

	if entry.Message != "Slow database query detected" {
		t.Errorf("Expected message 'Slow database query detected', got '%s'", entry.Message)
	}

	if entry.Fields["slow_query"] != true {
		t.Errorf("Expected slow_query true, got %v", entry.Fields["slow_query"])
	}

	if entry.Fields["duration"] != "2s" {
		t.Errorf("Expected duration '2s', got '%v'", entry.Fields["duration"])
	}

	if entry.Fields["threshold"] != "500ms" {
		t.Errorf("Expected threshold '500ms', got '%v'", entry.Fields["threshold"])
	}
}

func TestDBLogger_AutoSlowQueryDetection(t *testing.T) {
	var buf bytes.Buffer
	logger := NewJSONLogger(&buf, DebugLevel, "test")
	dbLogger := NewDBLogger(logger, 500*time.Millisecond)

	ctx := context.Background()
	query := "SELECT * FROM users"
	args := []interface{}{}
	duration := 2 * time.Second // This should trigger slow query warning

	dbLogger.LogQuery(ctx, query, args, duration, nil)

	// Should have 2 log entries: one for the query and one for slow query warning
	logOutput := buf.String()
	logLines := strings.Split(strings.TrimSpace(logOutput), "\n")
	
	if len(logLines) != 2 {
		t.Fatalf("Expected 2 log entries, got %d", len(logLines))
	}

	// Parse slow query warning (second entry)
	var slowQueryEntry LogEntry
	if err := json.Unmarshal([]byte(logLines[1]), &slowQueryEntry); err != nil {
		t.Fatalf("Failed to unmarshal slow query log entry: %v", err)
	}

	if slowQueryEntry.Level != "WARN" {
		t.Errorf("Expected slow query level WARN, got %s", slowQueryEntry.Level)
	}

	if slowQueryEntry.Message != "Slow database query detected" {
		t.Errorf("Expected slow query message, got '%s'", slowQueryEntry.Message)
	}
}

func TestSanitizeQuery(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "multiple spaces",
			input:    "SELECT   *   FROM   users   WHERE   id = $1",
			expected: "SELECT * FROM users WHERE id = $1",
		},
		{
			name:     "newlines and tabs",
			input:    "SELECT *\nFROM users\n\tWHERE id = $1",
			expected: "SELECT * FROM users WHERE id = $1",
		},
		{
			name:     "long query truncation",
			input:    strings.Repeat("SELECT * FROM very_long_table_name ", 20),
			expected: func() string {
				longQuery := strings.Repeat("SELECT * FROM very_long_table_name ", 20)
				if len(longQuery) > 500 {
					return longQuery[:497] + "..."
				}
				return longQuery
			}(),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := sanitizeQuery(test.input)
			if result != test.expected {
				t.Errorf("sanitizeQuery() = %v, want %v", result, test.expected)
			}
		})
	}
}

func TestSanitizeArgs(t *testing.T) {
	tests := []struct {
		name     string
		input    []interface{}
		expected []interface{}
	}{
		{
			name:     "normal args",
			input:    []interface{}{123, "john@example.com", true},
			expected: []interface{}{123, "john@example.com", true},
		},
		{
			name:     "password redaction",
			input:    []interface{}{"user123", "mypassword123"},
			expected: []interface{}{"user123", "[REDACTED]"},
		},
		{
			name:     "long string truncation",
			input:    []interface{}{strings.Repeat("a", 150)},
			expected: []interface{}{strings.Repeat("a", 97) + "..."},
		},
		{
			name:     "binary data",
			input:    []interface{}{[]byte{1, 2, 3, 4, 5}},
			expected: []interface{}{"[BINARY:5 bytes]"},
		},
		{
			name:     "token redaction",
			input:    []interface{}{"user123", "auth_token_abc123"},
			expected: []interface{}{"user123", "[REDACTED]"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := sanitizeArgs(test.input)
			if len(result) != len(test.expected) {
				t.Fatalf("Expected %d args, got %d", len(test.expected), len(result))
			}
			
			for i, expected := range test.expected {
				if result[i] != expected {
					t.Errorf("Arg %d: expected %v, got %v", i, expected, result[i])
				}
			}
		})
	}
}

func TestIsSensitiveString(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"password123", true},
		{"mytoken", true},
		{"secret_key", true},
		{"auth_header", true},
		{"user_credential", true},
		{"john@example.com", false},
		{"user123", false},
		{"normal_string", false},
		// Hash-like strings
		{"a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6", true},
		{"short", false},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			result := isSensitiveString(test.input)
			if result != test.expected {
				t.Errorf("isSensitiveString(%s) = %v, want %v", test.input, result, test.expected)
			}
		})
	}
}

func TestContainsMixedChars(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"abc123", true},
		{"ABC123", true},
		{"abcdef", false},
		{"123456", false},
		{"", false},
		{"a1", true},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			result := containsMixedChars(test.input)
			if result != test.expected {
				t.Errorf("containsMixedChars(%s) = %v, want %v", test.input, result, test.expected)
			}
		})
	}
}

// Integration test with actual database operations would require a test database
// For now, we'll test the wrapper structure and method signatures

func TestLoggingDB_Structure(t *testing.T) {
	var buf bytes.Buffer
	logger := NewJSONLogger(&buf, DebugLevel, "test")
	
	// We can't create a real database connection in unit tests without setup
	// So we'll test that the constructor works with nil (it should handle this gracefully)
	// In a real scenario, you'd pass a valid *sql.DB
	
	loggingDB := NewLoggingDB(nil, logger, 500*time.Millisecond)
	if loggingDB == nil {
		t.Fatal("NewLoggingDB returned nil")
	}
	
	if loggingDB.dbLogger == nil {
		t.Fatal("LoggingDB dbLogger is nil")
	}
}

// Mock test for database operations (would need actual DB for full integration test)
func TestLoggingDB_MethodSignatures(t *testing.T) {
	// This test ensures our method signatures are correct
	// In a real integration test, you'd use a test database
	
	var buf bytes.Buffer
	logger := NewJSONLogger(&buf, DebugLevel, "test")
	
	// Create a mock or use a test database
	// For this test, we're just checking that the methods exist and have correct signatures
	loggingDB := NewLoggingDB(nil, logger, 500*time.Millisecond)
	
	// Test that methods exist (they'll panic with nil DB, but signatures are correct)
	ctx := context.Background()
	
	// These would work with a real database connection
	_ = ctx
	_ = loggingDB
	
	// Just verify the struct has the expected methods by checking they compile
	var _ interface {
		QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
		QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
		ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
		PrepareContext(ctx context.Context, query string) (*LoggingStmt, error)
		BeginTx(ctx context.Context, opts *sql.TxOptions) (*LoggingTx, error)
		PingContext(ctx context.Context) error
		Close() error
		Stats() sql.DBStats
	} = loggingDB
}