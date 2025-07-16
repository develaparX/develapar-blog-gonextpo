package middleware

import (
	"bytes"
	"context"
	"develapar-server/utils"
	"io"
	"time"

	"github.com/gin-gonic/gin"
)

// RequestLogger interface for request logging middleware
type RequestLogger interface {
	LogRequests() gin.HandlerFunc
}

// requestLogger implements RequestLogger interface
type requestLogger struct {
	logger utils.Logger
}

// responseWriter wraps gin.ResponseWriter to capture response data
type responseWriter struct {
	gin.ResponseWriter
	body   *bytes.Buffer
	status int
}

// Write captures the response body
func (rw *responseWriter) Write(data []byte) (int, error) {
	rw.body.Write(data)
	return rw.ResponseWriter.Write(data)
}

// WriteHeader captures the status code
func (rw *responseWriter) WriteHeader(statusCode int) {
	rw.status = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}

// LogRequests middleware logs incoming requests and responses with context
func (rl *requestLogger) LogRequests() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start time for request processing
		startTime := time.Now()
		
		// Get context from request
		ctx := c.Request.Context()
		
		// Create response writer wrapper to capture response data
		responseBody := &bytes.Buffer{}
		writer := &responseWriter{
			ResponseWriter: c.Writer,
			body:          responseBody,
			status:        200, // Default status
		}
		c.Writer = writer

		// Log incoming request
		rl.logIncomingRequest(ctx, c, startTime)

		// Process request
		c.Next()

		// Calculate processing time
		processingTime := time.Since(startTime)

		// Log outgoing response
		rl.logOutgoingResponse(ctx, c, writer, processingTime)
	}
}

// logIncomingRequest logs details about the incoming request
func (rl *requestLogger) logIncomingRequest(ctx context.Context, c *gin.Context, startTime time.Time) {
	// Read and restore request body for logging
	var requestBody string
	if c.Request.Body != nil {
		bodyBytes, err := io.ReadAll(c.Request.Body)
		if err == nil {
			requestBody = string(bodyBytes)
			// Restore the request body for further processing
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}
	}

	// Get user agent and IP
	userAgent := c.GetHeader("User-Agent")
	clientIP := c.ClientIP()
	
	// Get content length
	contentLength := c.Request.ContentLength

	fields := []utils.Field{
		utils.StringField("method", c.Request.Method),
		utils.StringField("path", c.Request.URL.Path),
		utils.StringField("query", c.Request.URL.RawQuery),
		utils.StringField("user_agent", userAgent),
		utils.StringField("client_ip", clientIP),
		utils.Int64Field("content_length", contentLength),
		utils.StringField("protocol", c.Request.Proto),
	}

	// Add request body if present and not too large (limit to 1KB for logging)
	if requestBody != "" && len(requestBody) < 1024 {
		fields = append(fields, utils.StringField("request_body", requestBody))
	}

	// Add headers (excluding sensitive ones)
	headers := make(map[string]string)
	for key, values := range c.Request.Header {
		// Skip sensitive headers
		if !isSensitiveHeader(key) && len(values) > 0 {
			headers[key] = values[0]
		}
	}
	if len(headers) > 0 {
		fields = append(fields, utils.Field{Key: "headers", Value: headers})
	}

	rl.logger.Info(ctx, "Incoming request", fields...)
}

// logOutgoingResponse logs details about the outgoing response
func (rl *requestLogger) logOutgoingResponse(ctx context.Context, c *gin.Context, writer *responseWriter, processingTime time.Duration) {
	// Get response body (limit to 1KB for logging)
	responseBody := ""
	if writer.body.Len() > 0 && writer.body.Len() < 1024 {
		responseBody = writer.body.String()
	}

	// Get response headers
	responseHeaders := make(map[string]string)
	for key, values := range c.Writer.Header() {
		if len(values) > 0 {
			responseHeaders[key] = values[0]
		}
	}

	fields := []utils.Field{
		utils.StringField("method", c.Request.Method),
		utils.StringField("path", c.Request.URL.Path),
		utils.IntField("status_code", writer.status),
		utils.DurationField("processing_time", processingTime),
		utils.IntField("response_size", writer.body.Len()),
	}

	// Add response body if present
	if responseBody != "" {
		fields = append(fields, utils.StringField("response_body", responseBody))
	}

	// Add response headers
	if len(responseHeaders) > 0 {
		fields = append(fields, utils.Field{Key: "response_headers", Value: responseHeaders})
	}

	// Log with appropriate level based on status code
	if writer.status >= 500 {
		rl.logger.Error(ctx, "Request completed with server error", nil, fields...)
	} else if writer.status >= 400 {
		rl.logger.Warn(ctx, "Request completed with client error", fields...)
	} else {
		rl.logger.Info(ctx, "Request completed successfully", fields...)
	}
}

// isSensitiveHeader checks if a header contains sensitive information
func isSensitiveHeader(headerName string) bool {
	sensitiveHeaders := []string{
		"Authorization",
		"Cookie",
		"Set-Cookie",
		"X-Api-Key",
		"X-Auth-Token",
		"Proxy-Authorization",
	}

	for _, sensitive := range sensitiveHeaders {
		if headerName == sensitive {
			return true
		}
	}
	return false
}

// NewRequestLogger creates a new request logging middleware
func NewRequestLogger(logger utils.Logger) RequestLogger {
	return &requestLogger{
		logger: logger,
	}
}

// RequestMetrics holds metrics about request processing
type RequestMetrics struct {
	TotalRequests     int64
	RequestsByMethod  map[string]int64
	RequestsByStatus  map[int]int64
	AverageResponseTime time.Duration
	SlowRequests      int64 // Requests taking more than 1 second
}

// MetricsCollector interface for collecting request metrics
type MetricsCollector interface {
	RecordRequest(method string, status int, duration time.Duration)
	GetMetrics() RequestMetrics
	Reset()
}

// metricsCollector implements MetricsCollector interface
type metricsCollector struct {
	metrics RequestMetrics
}

// RecordRequest records metrics for a completed request
func (mc *metricsCollector) RecordRequest(method string, status int, duration time.Duration) {
	mc.metrics.TotalRequests++
	
	if mc.metrics.RequestsByMethod == nil {
		mc.metrics.RequestsByMethod = make(map[string]int64)
	}
	mc.metrics.RequestsByMethod[method]++
	
	if mc.metrics.RequestsByStatus == nil {
		mc.metrics.RequestsByStatus = make(map[int]int64)
	}
	mc.metrics.RequestsByStatus[status]++
	
	// Update average response time (simple moving average)
	if mc.metrics.TotalRequests == 1 {
		mc.metrics.AverageResponseTime = duration
	} else {
		mc.metrics.AverageResponseTime = time.Duration(
			(int64(mc.metrics.AverageResponseTime)*(mc.metrics.TotalRequests-1) + int64(duration)) / mc.metrics.TotalRequests,
		)
	}
	
	// Count slow requests (> 1 second)
	if duration > time.Second {
		mc.metrics.SlowRequests++
	}
}

// GetMetrics returns current metrics
func (mc *metricsCollector) GetMetrics() RequestMetrics {
	return mc.metrics
}

// Reset resets all metrics
func (mc *metricsCollector) Reset() {
	mc.metrics = RequestMetrics{
		RequestsByMethod: make(map[string]int64),
		RequestsByStatus: make(map[int]int64),
	}
}

// NewMetricsCollector creates a new metrics collector
func NewMetricsCollector() MetricsCollector {
	return &metricsCollector{
		metrics: RequestMetrics{
			RequestsByMethod: make(map[string]int64),
			RequestsByStatus: make(map[int]int64),
		},
	}
}

// RequestLoggerWithMetrics combines request logging with metrics collection
type RequestLoggerWithMetrics struct {
	logger    utils.Logger
	collector MetricsCollector
}

// LogRequestsWithMetrics middleware logs requests and collects metrics
func (rlm *RequestLoggerWithMetrics) LogRequestsWithMetrics() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()
		ctx := c.Request.Context()
		
		// Create response writer wrapper
		responseBody := &bytes.Buffer{}
		writer := &responseWriter{
			ResponseWriter: c.Writer,
			body:          responseBody,
			status:        200,
		}
		c.Writer = writer

		// Log incoming request
		rlm.logIncomingRequest(ctx, c, startTime)

		// Process request
		c.Next()

		// Calculate processing time
		processingTime := time.Since(startTime)

		// Record metrics
		rlm.collector.RecordRequest(c.Request.Method, writer.status, processingTime)

		// Log outgoing response with metrics
		rlm.logOutgoingResponseWithMetrics(ctx, c, writer, processingTime)
	}
}

// logIncomingRequest logs incoming request (same as regular request logger)
func (rlm *RequestLoggerWithMetrics) logIncomingRequest(ctx context.Context, c *gin.Context, startTime time.Time) {
	var requestBody string
	if c.Request.Body != nil {
		bodyBytes, err := io.ReadAll(c.Request.Body)
		if err == nil {
			requestBody = string(bodyBytes)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}
	}

	userAgent := c.GetHeader("User-Agent")
	clientIP := c.ClientIP()
	contentLength := c.Request.ContentLength

	fields := []utils.Field{
		utils.StringField("method", c.Request.Method),
		utils.StringField("path", c.Request.URL.Path),
		utils.StringField("query", c.Request.URL.RawQuery),
		utils.StringField("user_agent", userAgent),
		utils.StringField("client_ip", clientIP),
		utils.Int64Field("content_length", contentLength),
		utils.StringField("protocol", c.Request.Proto),
	}

	if requestBody != "" && len(requestBody) < 1024 {
		fields = append(fields, utils.StringField("request_body", requestBody))
	}

	headers := make(map[string]string)
	for key, values := range c.Request.Header {
		if !isSensitiveHeader(key) && len(values) > 0 {
			headers[key] = values[0]
		}
	}
	if len(headers) > 0 {
		fields = append(fields, utils.Field{Key: "headers", Value: headers})
	}

	rlm.logger.Info(ctx, "Incoming request", fields...)
}

// logOutgoingResponseWithMetrics logs response with additional metrics information
func (rlm *RequestLoggerWithMetrics) logOutgoingResponseWithMetrics(ctx context.Context, c *gin.Context, writer *responseWriter, processingTime time.Duration) {
	responseBody := ""
	if writer.body.Len() > 0 && writer.body.Len() < 1024 {
		responseBody = writer.body.String()
	}

	responseHeaders := make(map[string]string)
	for key, values := range c.Writer.Header() {
		if len(values) > 0 {
			responseHeaders[key] = values[0]
		}
	}

	// Get current metrics for additional context
	metrics := rlm.collector.GetMetrics()

	fields := []utils.Field{
		utils.StringField("method", c.Request.Method),
		utils.StringField("path", c.Request.URL.Path),
		utils.IntField("status_code", writer.status),
		utils.DurationField("processing_time", processingTime),
		utils.IntField("response_size", writer.body.Len()),
		utils.Int64Field("total_requests", metrics.TotalRequests),
		utils.DurationField("avg_response_time", metrics.AverageResponseTime),
		utils.BoolField("slow_request", processingTime > time.Second),
	}

	if responseBody != "" {
		fields = append(fields, utils.StringField("response_body", responseBody))
	}

	if len(responseHeaders) > 0 {
		fields = append(fields, utils.Field{Key: "response_headers", Value: responseHeaders})
	}

	// Log with appropriate level
	if writer.status >= 500 {
		rlm.logger.Error(ctx, "Request completed with server error", nil, fields...)
	} else if writer.status >= 400 {
		rlm.logger.Warn(ctx, "Request completed with client error", fields...)
	} else {
		rlm.logger.Info(ctx, "Request completed successfully", fields...)
	}
}

// GetMetrics returns current request metrics
func (rlm *RequestLoggerWithMetrics) GetMetrics() RequestMetrics {
	return rlm.collector.GetMetrics()
}

// NewRequestLoggerWithMetrics creates a new request logger with metrics collection
func NewRequestLoggerWithMetrics(logger utils.Logger) *RequestLoggerWithMetrics {
	return &RequestLoggerWithMetrics{
		logger:    logger,
		collector: NewMetricsCollector(),
	}
}