package middleware

import (
	"develapar-server/service"
	"develapar-server/utils"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
)

// MetricsMiddleware interface for collecting request metrics
type MetricsMiddleware interface {
	CollectMetrics() gin.HandlerFunc
	CollectSystemMetrics() gin.HandlerFunc
}

// metricsMiddleware implements MetricsMiddleware interface
type metricsMiddleware struct {
	metricsService service.MetricsService
	logger         utils.Logger
}

// NewMetricsMiddleware creates a new metrics middleware
func NewMetricsMiddleware(metricsService service.MetricsService, logger utils.Logger) MetricsMiddleware {
	return &metricsMiddleware{
		metricsService: metricsService,
		logger:         logger,
	}
}

// CollectMetrics middleware collects HTTP request metrics with context
func (mm *metricsMiddleware) CollectMetrics() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()
		ctx := c.Request.Context()
		
		// Get request information
		method := c.Request.Method
		path := c.FullPath()
		if path == "" {
			path = c.Request.URL.Path
		}
		
		// Process request
		c.Next()
		
		// Calculate duration
		duration := time.Since(startTime)
		statusCode := c.Writer.Status()
		
		// Record request metrics with context
		mm.metricsService.RecordRequest(ctx, method, path, statusCode, duration)
		
		// Record error if status code indicates an error
		if statusCode >= 400 {
			errorType := mm.getErrorType(statusCode)
			operation := method + " " + path
			mm.metricsService.RecordError(ctx, errorType, operation)
		}
		
		// Log request completion with context
		requestID := getRequestIDFromGinContext(c)
		userID := getUserIDFromGinContext(c)
		
		mm.logger.Info(ctx, "Request completed",
			utils.Field{Key: "request_id", Value: requestID},
			utils.Field{Key: "user_id", Value: userID},
			utils.Field{Key: "method", Value: method},
			utils.Field{Key: "path", Value: path},
			utils.Field{Key: "status_code", Value: statusCode},
			utils.Field{Key: "duration_ms", Value: duration.Milliseconds()},
		)
	}
}

// CollectSystemMetrics middleware periodically collects system metrics
func (mm *metricsMiddleware) CollectSystemMetrics() gin.HandlerFunc {
	// Start a goroutine to collect system metrics periodically
	go mm.collectSystemMetricsRoutine()
	
	// Return a no-op middleware since system metrics are collected in background
	return func(c *gin.Context) {
		c.Next()
	}
}

// collectSystemMetricsRoutine runs in background to collect system metrics
func (mm *metricsMiddleware) collectSystemMetricsRoutine() {
	ticker := time.NewTicker(30 * time.Second) // Collect every 30 seconds
	defer ticker.Stop()
	
	for range ticker.C {
		ctx := utils.NewBackgroundContext()
		
		// Collect memory statistics
		var memStats runtime.MemStats
		runtime.ReadMemStats(&memStats)
		
		mm.metricsService.RecordMemoryUsage(ctx, memStats.Alloc, memStats.Sys)
		
		// Collect goroutine count
		goroutineCount := runtime.NumGoroutine()
		mm.metricsService.RecordGoroutineCount(ctx, goroutineCount)
		
		// Log system metrics collection
		mm.logger.Debug(ctx, "System metrics collected",
			utils.Field{Key: "alloc_mb", Value: memStats.Alloc / 1024 / 1024},
			utils.Field{Key: "sys_mb", Value: memStats.Sys / 1024 / 1024},
			utils.Field{Key: "goroutines", Value: goroutineCount},
		)
	}
}

// getErrorType categorizes HTTP status codes into error types
func (mm *metricsMiddleware) getErrorType(statusCode int) string {
	switch {
	case statusCode >= 400 && statusCode < 500:
		switch statusCode {
		case 400:
			return "bad_request"
		case 401:
			return "unauthorized"
		case 403:
			return "forbidden"
		case 404:
			return "not_found"
		case 409:
			return "conflict"
		case 422:
			return "validation_error"
		case 429:
			return "rate_limit"
		default:
			return "client_error"
		}
	case statusCode >= 500:
		switch statusCode {
		case 500:
			return "internal_error"
		case 502:
			return "bad_gateway"
		case 503:
			return "service_unavailable"
		case 504:
			return "gateway_timeout"
		default:
			return "server_error"
		}
	default:
		return "unknown_error"
	}
}

// Helper functions to get context information from Gin context
func getRequestIDFromGinContext(c *gin.Context) string {
	if requestID, exists := c.Get("request_id"); exists {
		if id, ok := requestID.(string); ok {
			return id
		}
	}
	return ""
}

func getUserIDFromGinContext(c *gin.Context) string {
	if userID, exists := c.Get("user_id"); exists {
		if id, ok := userID.(string); ok {
			return id
		}
	}
	return ""
}

// DatabaseMetricsCollector interface for collecting database metrics
type DatabaseMetricsCollector interface {
	RecordQuery(ctx gin.Context, operation string, duration time.Duration, success bool)
	RecordConnectionStats(ctx gin.Context, open, inUse, idle int)
}

// databaseMetricsCollector implements DatabaseMetricsCollector
type databaseMetricsCollector struct {
	metricsService service.MetricsService
	logger         utils.Logger
}

// NewDatabaseMetricsCollector creates a new database metrics collector
func NewDatabaseMetricsCollector(metricsService service.MetricsService, logger utils.Logger) DatabaseMetricsCollector {
	return &databaseMetricsCollector{
		metricsService: metricsService,
		logger:         logger,
	}
}

// RecordQuery records database query metrics
func (dmc *databaseMetricsCollector) RecordQuery(c gin.Context, operation string, duration time.Duration, success bool) {
	ctx := c.Request.Context()
	dmc.metricsService.RecordDatabaseQuery(ctx, operation, duration, success)
}

// RecordConnectionStats records database connection pool statistics
func (dmc *databaseMetricsCollector) RecordConnectionStats(c gin.Context, open, inUse, idle int) {
	ctx := c.Request.Context()
	dmc.metricsService.RecordConnectionPoolStats(ctx, open, inUse, idle)
}