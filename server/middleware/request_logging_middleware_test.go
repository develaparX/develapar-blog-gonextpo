package middleware

import (
	"bytes"
	"develapar-server/utils"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

func TestNewRequestLogger(t *testing.T) {
	var buf bytes.Buffer
	logger := utils.NewJSONLogger(&buf, utils.InfoLevel, "test")
	
	requestLogger := NewRequestLogger(logger)
	if requestLogger == nil {
		t.Fatal("NewRequestLogger returned nil")
	}
}

func TestRequestLogger_LogRequests(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)
	
	var buf bytes.Buffer
	logger := utils.NewJSONLogger(&buf, utils.InfoLevel, "test")
	requestLogger := NewRequestLogger(logger)

	// Create test router
	router := gin.New()
	
	// Add context middleware first
	contextManager := NewContextManager()
	contextMiddleware := NewContextMiddleware(contextManager)
	router.Use(contextMiddleware.InjectContext())
	
	// Add request logging middleware
	router.Use(requestLogger.LogRequests())
	
	// Add test endpoint
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "test response"})
	})

	// Create test request
	req := httptest.NewRequest("GET", "/test?param=value", nil)
	req.Header.Set("User-Agent", "test-agent")
	req.Header.Set("X-Custom-Header", "custom-value")
	
	// Create response recorder
	w := httptest.NewRecorder()
	
	// Perform request
	router.ServeHTTP(w, req)

	// Check response
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	// Parse log entries
	logOutput := buf.String()
	logLines := strings.Split(strings.TrimSpace(logOutput), "\n")
	
	if len(logLines) < 2 {
		t.Fatalf("Expected at least 2 log entries (request and response), got %d", len(logLines))
	}

	// Parse first log entry (incoming request)
	var incomingLog utils.LogEntry
	if err := json.Unmarshal([]byte(logLines[0]), &incomingLog); err != nil {
		t.Fatalf("Failed to parse incoming request log: %v", err)
	}

	// Verify incoming request log
	if incomingLog.Level != "INFO" {
		t.Errorf("Expected level INFO, got %s", incomingLog.Level)
	}
	
	if incomingLog.Message != "Incoming request" {
		t.Errorf("Expected message 'Incoming request', got '%s'", incomingLog.Message)
	}

	if incomingLog.Fields["method"] != "GET" {
		t.Errorf("Expected method GET, got %v", incomingLog.Fields["method"])
	}

	if incomingLog.Fields["path"] != "/test" {
		t.Errorf("Expected path /test, got %v", incomingLog.Fields["path"])
	}

	if incomingLog.Fields["query"] != "param=value" {
		t.Errorf("Expected query 'param=value', got %v", incomingLog.Fields["query"])
	}

	// Parse second log entry (outgoing response)
	var outgoingLog utils.LogEntry
	if err := json.Unmarshal([]byte(logLines[1]), &outgoingLog); err != nil {
		t.Fatalf("Failed to parse outgoing response log: %v", err)
	}

	// Verify outgoing response log
	if outgoingLog.Level != "INFO" {
		t.Errorf("Expected level INFO, got %s", outgoingLog.Level)
	}
	
	if outgoingLog.Message != "Request completed successfully" {
		t.Errorf("Expected message 'Request completed successfully', got '%s'", outgoingLog.Message)
	}

	if outgoingLog.Fields["status_code"] != float64(200) {
		t.Errorf("Expected status_code 200, got %v", outgoingLog.Fields["status_code"])
	}

	// Check that processing_time is present
	if _, exists := outgoingLog.Fields["processing_time"]; !exists {
		t.Error("Expected processing_time field in response log")
	}
}

func TestRequestLogger_ErrorResponse(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	var buf bytes.Buffer
	logger := utils.NewJSONLogger(&buf, utils.InfoLevel, "test")
	requestLogger := NewRequestLogger(logger)

	router := gin.New()
	
	contextManager := NewContextManager()
	contextMiddleware := NewContextMiddleware(contextManager)
	router.Use(contextMiddleware.InjectContext())
	router.Use(requestLogger.LogRequests())
	
	// Add endpoint that returns error
	router.GET("/error", func(c *gin.Context) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
	})

	req := httptest.NewRequest("GET", "/error", nil)
	w := httptest.NewRecorder()
	
	router.ServeHTTP(w, req)

	// Parse log entries
	logOutput := buf.String()
	logLines := strings.Split(strings.TrimSpace(logOutput), "\n")
	
	if len(logLines) < 2 {
		t.Fatalf("Expected at least 2 log entries, got %d", len(logLines))
	}

	// Parse response log (should be ERROR level)
	var responseLog utils.LogEntry
	if err := json.Unmarshal([]byte(logLines[1]), &responseLog); err != nil {
		t.Fatalf("Failed to parse response log: %v", err)
	}

	if responseLog.Level != "ERROR" {
		t.Errorf("Expected level ERROR for 500 status, got %s", responseLog.Level)
	}

	if responseLog.Message != "Request completed with server error" {
		t.Errorf("Expected server error message, got '%s'", responseLog.Message)
	}
}

func TestRequestLogger_ClientError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	var buf bytes.Buffer
	logger := utils.NewJSONLogger(&buf, utils.InfoLevel, "test")
	requestLogger := NewRequestLogger(logger)

	router := gin.New()
	
	contextManager := NewContextManager()
	contextMiddleware := NewContextMiddleware(contextManager)
	router.Use(contextMiddleware.InjectContext())
	router.Use(requestLogger.LogRequests())
	
	// Add endpoint that returns client error
	router.GET("/bad-request", func(c *gin.Context) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
	})

	req := httptest.NewRequest("GET", "/bad-request", nil)
	w := httptest.NewRecorder()
	
	router.ServeHTTP(w, req)

	// Parse log entries
	logOutput := buf.String()
	logLines := strings.Split(strings.TrimSpace(logOutput), "\n")
	
	if len(logLines) < 2 {
		t.Fatalf("Expected at least 2 log entries, got %d", len(logLines))
	}

	// Parse response log (should be WARN level)
	var responseLog utils.LogEntry
	if err := json.Unmarshal([]byte(logLines[1]), &responseLog); err != nil {
		t.Fatalf("Failed to parse response log: %v", err)
	}

	if responseLog.Level != "WARN" {
		t.Errorf("Expected level WARN for 400 status, got %s", responseLog.Level)
	}

	if responseLog.Message != "Request completed with client error" {
		t.Errorf("Expected client error message, got '%s'", responseLog.Message)
	}
}

func TestRequestLogger_WithRequestBody(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	var buf bytes.Buffer
	logger := utils.NewJSONLogger(&buf, utils.InfoLevel, "test")
	requestLogger := NewRequestLogger(logger)

	router := gin.New()
	
	contextManager := NewContextManager()
	contextMiddleware := NewContextMiddleware(contextManager)
	router.Use(contextMiddleware.InjectContext())
	router.Use(requestLogger.LogRequests())
	
	router.POST("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"received": "ok"})
	})

	requestBody := `{"test": "data"}`
	req := httptest.NewRequest("POST", "/test", strings.NewReader(requestBody))
	req.Header.Set("Content-Type", "application/json")
	
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Parse log entries
	logOutput := buf.String()
	logLines := strings.Split(strings.TrimSpace(logOutput), "\n")
	
	if len(logLines) < 1 {
		t.Fatalf("Expected at least 1 log entry, got %d", len(logLines))
	}

	// Parse incoming request log
	var incomingLog utils.LogEntry
	if err := json.Unmarshal([]byte(logLines[0]), &incomingLog); err != nil {
		t.Fatalf("Failed to parse incoming request log: %v", err)
	}

	// Check that request body is logged
	if incomingLog.Fields["request_body"] != requestBody {
		t.Errorf("Expected request body '%s', got '%v'", requestBody, incomingLog.Fields["request_body"])
	}
}

func TestIsSensitiveHeader(t *testing.T) {
	tests := []struct {
		header   string
		expected bool
	}{
		{"Authorization", true},
		{"Cookie", true},
		{"Set-Cookie", true},
		{"X-Api-Key", true},
		{"X-Auth-Token", true},
		{"Proxy-Authorization", true},
		{"Content-Type", false},
		{"User-Agent", false},
		{"X-Custom-Header", false},
	}

	for _, test := range tests {
		t.Run(test.header, func(t *testing.T) {
			result := isSensitiveHeader(test.header)
			if result != test.expected {
				t.Errorf("isSensitiveHeader(%s) = %v, want %v", test.header, result, test.expected)
			}
		})
	}
}

func TestMetricsCollector(t *testing.T) {
	collector := NewMetricsCollector()

	// Record some requests
	collector.RecordRequest("GET", 200, 100*time.Millisecond)
	collector.RecordRequest("POST", 201, 200*time.Millisecond)
	collector.RecordRequest("GET", 404, 50*time.Millisecond)
	collector.RecordRequest("GET", 500, 1500*time.Millisecond) // Slow request

	metrics := collector.GetMetrics()

	// Check total requests
	if metrics.TotalRequests != 4 {
		t.Errorf("Expected 4 total requests, got %d", metrics.TotalRequests)
	}

	// Check requests by method
	if metrics.RequestsByMethod["GET"] != 3 {
		t.Errorf("Expected 3 GET requests, got %d", metrics.RequestsByMethod["GET"])
	}
	if metrics.RequestsByMethod["POST"] != 1 {
		t.Errorf("Expected 1 POST request, got %d", metrics.RequestsByMethod["POST"])
	}

	// Check requests by status
	if metrics.RequestsByStatus[200] != 1 {
		t.Errorf("Expected 1 request with status 200, got %d", metrics.RequestsByStatus[200])
	}
	if metrics.RequestsByStatus[404] != 1 {
		t.Errorf("Expected 1 request with status 404, got %d", metrics.RequestsByStatus[404])
	}

	// Check slow requests
	if metrics.SlowRequests != 1 {
		t.Errorf("Expected 1 slow request, got %d", metrics.SlowRequests)
	}

	// Test reset
	collector.Reset()
	resetMetrics := collector.GetMetrics()
	if resetMetrics.TotalRequests != 0 {
		t.Errorf("Expected 0 total requests after reset, got %d", resetMetrics.TotalRequests)
	}
}

func TestRequestLoggerWithMetrics(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	var buf bytes.Buffer
	logger := utils.NewJSONLogger(&buf, utils.InfoLevel, "test")
	requestLoggerWithMetrics := NewRequestLoggerWithMetrics(logger)

	router := gin.New()
	
	contextManager := NewContextManager()
	contextMiddleware := NewContextMiddleware(contextManager)
	router.Use(contextMiddleware.InjectContext())
	router.Use(requestLoggerWithMetrics.LogRequestsWithMetrics())
	
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "test"})
	})

	// Make a request
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Check metrics
	metrics := requestLoggerWithMetrics.GetMetrics()
	if metrics.TotalRequests != 1 {
		t.Errorf("Expected 1 total request, got %d", metrics.TotalRequests)
	}

	// Parse log entries to check metrics fields
	logOutput := buf.String()
	logLines := strings.Split(strings.TrimSpace(logOutput), "\n")
	
	if len(logLines) < 2 {
		t.Fatalf("Expected at least 2 log entries, got %d", len(logLines))
	}

	// Parse response log
	var responseLog utils.LogEntry
	if err := json.Unmarshal([]byte(logLines[1]), &responseLog); err != nil {
		t.Fatalf("Failed to parse response log: %v", err)
	}

	// Check that metrics fields are present
	if responseLog.Fields["total_requests"] != float64(1) {
		t.Errorf("Expected total_requests 1, got %v", responseLog.Fields["total_requests"])
	}

	if _, exists := responseLog.Fields["avg_response_time"]; !exists {
		t.Error("Expected avg_response_time field in response log")
	}

	if responseLog.Fields["slow_request"] != false {
		t.Errorf("Expected slow_request false, got %v", responseLog.Fields["slow_request"])
	}
}