package service

import (
	"context"
	"develapar-server/utils"
	"sync"
	"time"
)

// MetricsService interface for collecting and reporting application metrics with context
type MetricsService interface {
	// Request metrics
	RecordRequest(ctx context.Context, method, path string, statusCode int, duration time.Duration)
	RecordError(ctx context.Context, errorType, operation string)
	
	// Database metrics
	RecordDatabaseQuery(ctx context.Context, operation string, duration time.Duration, success bool)
	RecordConnectionPoolStats(ctx context.Context, open, inUse, idle int)
	
	// Application metrics
	RecordMemoryUsage(ctx context.Context, allocBytes, sysBytes uint64)
	RecordGoroutineCount(ctx context.Context, count int)
	
	// Get metrics
	GetRequestMetrics(ctx context.Context) RequestMetrics
	GetDatabaseMetrics(ctx context.Context) DatabaseMetrics
	GetApplicationMetrics(ctx context.Context) ApplicationMetrics
	GetErrorMetrics(ctx context.Context) ErrorMetrics
	GetAllMetrics(ctx context.Context) AllMetrics
	
	// Reset metrics
	ResetMetrics(ctx context.Context)
}

// RequestMetrics represents HTTP request metrics
type RequestMetrics struct {
	TotalRequests     int64                    `json:"total_requests"`
	RequestsByMethod  map[string]int64         `json:"requests_by_method"`
	RequestsByPath    map[string]int64         `json:"requests_by_path"`
	RequestsByStatus  map[int]int64            `json:"requests_by_status"`
	AverageLatency    time.Duration            `json:"average_latency_ms"`
	P95Latency        time.Duration            `json:"p95_latency_ms"`
	P99Latency        time.Duration            `json:"p99_latency_ms"`
	RequestsPerSecond float64                  `json:"requests_per_second"`
	LastUpdated       time.Time                `json:"last_updated"`
}

// DatabaseMetrics represents database performance metrics
type DatabaseMetrics struct {
	TotalQueries        int64                    `json:"total_queries"`
	SuccessfulQueries   int64                    `json:"successful_queries"`
	FailedQueries       int64                    `json:"failed_queries"`
	AverageQueryTime    time.Duration            `json:"average_query_time_ms"`
	SlowQueries         int64                    `json:"slow_queries"`
	QueriesByOperation  map[string]int64         `json:"queries_by_operation"`
	ConnectionPoolStats ConnectionPoolMetrics    `json:"connection_pool_stats"`
	LastUpdated         time.Time                `json:"last_updated"`
}

// ConnectionPoolMetrics represents connection pool metrics
type ConnectionPoolMetrics struct {
	OpenConnections  int       `json:"open_connections"`
	InUseConnections int       `json:"in_use_connections"`
	IdleConnections  int       `json:"idle_connections"`
	LastUpdated      time.Time `json:"last_updated"`
}

// ApplicationMetrics represents application-level metrics
type ApplicationMetrics struct {
	MemoryUsage     MemoryMetrics `json:"memory_usage"`
	GoroutineCount  int           `json:"goroutine_count"`
	Uptime          time.Duration `json:"uptime"`
	StartTime       time.Time     `json:"start_time"`
	LastUpdated     time.Time     `json:"last_updated"`
}

// MemoryMetrics represents memory usage metrics
type MemoryMetrics struct {
	AllocBytes     uint64    `json:"alloc_bytes"`
	SysBytes       uint64    `json:"sys_bytes"`
	HeapAllocBytes uint64    `json:"heap_alloc_bytes"`
	HeapSysBytes   uint64    `json:"heap_sys_bytes"`
	LastUpdated    time.Time `json:"last_updated"`
}

// ErrorMetrics represents error tracking metrics
type ErrorMetrics struct {
	TotalErrors       int64            `json:"total_errors"`
	ErrorsByType      map[string]int64 `json:"errors_by_type"`
	ErrorsByOperation map[string]int64 `json:"errors_by_operation"`
	ErrorRate         float64          `json:"error_rate"`
	LastUpdated       time.Time        `json:"last_updated"`
}

// AllMetrics represents all metrics combined
type AllMetrics struct {
	Request     RequestMetrics     `json:"request"`
	Database    DatabaseMetrics    `json:"database"`
	Application ApplicationMetrics `json:"application"`
	Error       ErrorMetrics       `json:"error"`
	Timestamp   time.Time          `json:"timestamp"`
}

// metricsService implements MetricsService interface
type metricsService struct {
	mu                sync.RWMutex
	logger            utils.Logger
	startTime         time.Time
	
	// Request metrics
	totalRequests     int64
	requestsByMethod  map[string]int64
	requestsByPath    map[string]int64
	requestsByStatus  map[int]int64
	requestLatencies  []time.Duration
	requestStartTime  time.Time
	
	// Database metrics
	totalQueries      int64
	successfulQueries int64
	failedQueries     int64
	queryLatencies    []time.Duration
	queriesByOp       map[string]int64
	connectionStats   ConnectionPoolMetrics
	
	// Application metrics
	memoryStats       MemoryMetrics
	goroutineCount    int
	
	// Error metrics
	totalErrors       int64
	errorsByType      map[string]int64
	errorsByOperation map[string]int64
}

// NewMetricsService creates a new metrics service with context support
func NewMetricsService(logger utils.Logger) MetricsService {
	return &metricsService{
		logger:            logger,
		startTime:         time.Now(),
		requestsByMethod:  make(map[string]int64),
		requestsByPath:    make(map[string]int64),
		requestsByStatus:  make(map[int]int64),
		requestLatencies:  make([]time.Duration, 0, 1000), // Keep last 1000 requests
		requestStartTime:  time.Now(),
		queriesByOp:       make(map[string]int64),
		queryLatencies:    make([]time.Duration, 0, 1000), // Keep last 1000 queries
		errorsByType:      make(map[string]int64),
		errorsByOperation: make(map[string]int64),
	}
}

// RecordRequest records HTTP request metrics with context
func (ms *metricsService) RecordRequest(ctx context.Context, method, path string, statusCode int, duration time.Duration) {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	
	ms.totalRequests++
	ms.requestsByMethod[method]++
	ms.requestsByPath[path]++
	ms.requestsByStatus[statusCode]++
	
	// Keep only the last 1000 latencies for percentile calculations
	if len(ms.requestLatencies) >= 1000 {
		ms.requestLatencies = ms.requestLatencies[1:]
	}
	ms.requestLatencies = append(ms.requestLatencies, duration)
	
	// Log request metrics with context
	requestID := utils.GetRequestIDFromContext(ctx)
	userID := utils.GetUserIDFromContext(ctx)
	
	ms.logger.Debug(ctx, "Request recorded", 
		utils.Field{Key: "request_id", Value: requestID},
		utils.Field{Key: "user_id", Value: userID},
		utils.Field{Key: "method", Value: method},
		utils.Field{Key: "path", Value: path},
		utils.Field{Key: "status_code", Value: statusCode},
		utils.Field{Key: "duration_ms", Value: duration.Milliseconds()},
	)
}

// RecordError records error metrics with context
func (ms *metricsService) RecordError(ctx context.Context, errorType, operation string) {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	
	ms.totalErrors++
	ms.errorsByType[errorType]++
	ms.errorsByOperation[operation]++
	
	// Log error metrics with context
	requestID := utils.GetRequestIDFromContext(ctx)
	userID := utils.GetUserIDFromContext(ctx)
	
	ms.logger.Info(ctx, "Error recorded", 
		utils.Field{Key: "request_id", Value: requestID},
		utils.Field{Key: "user_id", Value: userID},
		utils.Field{Key: "error_type", Value: errorType},
		utils.Field{Key: "operation", Value: operation},
	)
}

// RecordDatabaseQuery records database query metrics with context
func (ms *metricsService) RecordDatabaseQuery(ctx context.Context, operation string, duration time.Duration, success bool) {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	
	ms.totalQueries++
	ms.queriesByOp[operation]++
	
	if success {
		ms.successfulQueries++
	} else {
		ms.failedQueries++
	}
	
	// Keep only the last 1000 query latencies
	if len(ms.queryLatencies) >= 1000 {
		ms.queryLatencies = ms.queryLatencies[1:]
	}
	ms.queryLatencies = append(ms.queryLatencies, duration)
	
	// Log database metrics with context
	requestID := utils.GetRequestIDFromContext(ctx)
	
	ms.logger.Debug(ctx, "Database query recorded", 
		utils.Field{Key: "request_id", Value: requestID},
		utils.Field{Key: "operation", Value: operation},
		utils.Field{Key: "duration_ms", Value: duration.Milliseconds()},
		utils.Field{Key: "success", Value: success},
	)
}

// RecordConnectionPoolStats records connection pool statistics with context
func (ms *metricsService) RecordConnectionPoolStats(ctx context.Context, open, inUse, idle int) {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	
	ms.connectionStats = ConnectionPoolMetrics{
		OpenConnections:  open,
		InUseConnections: inUse,
		IdleConnections:  idle,
		LastUpdated:      time.Now(),
	}
}

// RecordMemoryUsage records memory usage metrics with context
func (ms *metricsService) RecordMemoryUsage(ctx context.Context, allocBytes, sysBytes uint64) {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	
	ms.memoryStats = MemoryMetrics{
		AllocBytes:  allocBytes,
		SysBytes:    sysBytes,
		LastUpdated: time.Now(),
	}
}

// RecordGoroutineCount records goroutine count with context
func (ms *metricsService) RecordGoroutineCount(ctx context.Context, count int) {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	
	ms.goroutineCount = count
}

// GetRequestMetrics returns request metrics with context
func (ms *metricsService) GetRequestMetrics(ctx context.Context) RequestMetrics {
	ms.mu.RLock()
	defer ms.mu.RUnlock()
	
	avgLatency := ms.calculateAverageLatency(ms.requestLatencies)
	p95Latency := ms.calculatePercentile(ms.requestLatencies, 0.95)
	p99Latency := ms.calculatePercentile(ms.requestLatencies, 0.99)
	
	// Calculate requests per second
	elapsed := time.Since(ms.requestStartTime).Seconds()
	rps := float64(ms.totalRequests) / elapsed
	
	return RequestMetrics{
		TotalRequests:     ms.totalRequests,
		RequestsByMethod:  ms.copyStringInt64Map(ms.requestsByMethod),
		RequestsByPath:    ms.copyStringInt64Map(ms.requestsByPath),
		RequestsByStatus:  ms.copyIntInt64Map(ms.requestsByStatus),
		AverageLatency:    avgLatency,
		P95Latency:        p95Latency,
		P99Latency:        p99Latency,
		RequestsPerSecond: rps,
		LastUpdated:       time.Now(),
	}
}

// GetDatabaseMetrics returns database metrics with context
func (ms *metricsService) GetDatabaseMetrics(ctx context.Context) DatabaseMetrics {
	ms.mu.RLock()
	defer ms.mu.RUnlock()
	
	avgQueryTime := ms.calculateAverageLatency(ms.queryLatencies)
	slowQueries := ms.countSlowQueries(ms.queryLatencies, 1*time.Second)
	
	return DatabaseMetrics{
		TotalQueries:        ms.totalQueries,
		SuccessfulQueries:   ms.successfulQueries,
		FailedQueries:       ms.failedQueries,
		AverageQueryTime:    avgQueryTime,
		SlowQueries:         slowQueries,
		QueriesByOperation:  ms.copyStringInt64Map(ms.queriesByOp),
		ConnectionPoolStats: ms.connectionStats,
		LastUpdated:         time.Now(),
	}
}

// GetApplicationMetrics returns application metrics with context
func (ms *metricsService) GetApplicationMetrics(ctx context.Context) ApplicationMetrics {
	ms.mu.RLock()
	defer ms.mu.RUnlock()
	
	return ApplicationMetrics{
		MemoryUsage:    ms.memoryStats,
		GoroutineCount: ms.goroutineCount,
		Uptime:         time.Since(ms.startTime),
		StartTime:      ms.startTime,
		LastUpdated:    time.Now(),
	}
}

// GetErrorMetrics returns error metrics with context
func (ms *metricsService) GetErrorMetrics(ctx context.Context) ErrorMetrics {
	ms.mu.RLock()
	defer ms.mu.RUnlock()
	
	// Calculate error rate
	errorRate := float64(0)
	if ms.totalRequests > 0 {
		errorRate = float64(ms.totalErrors) / float64(ms.totalRequests) * 100
	}
	
	return ErrorMetrics{
		TotalErrors:       ms.totalErrors,
		ErrorsByType:      ms.copyStringInt64Map(ms.errorsByType),
		ErrorsByOperation: ms.copyStringInt64Map(ms.errorsByOperation),
		ErrorRate:         errorRate,
		LastUpdated:       time.Now(),
	}
}

// GetAllMetrics returns all metrics combined with context
func (ms *metricsService) GetAllMetrics(ctx context.Context) AllMetrics {
	return AllMetrics{
		Request:     ms.GetRequestMetrics(ctx),
		Database:    ms.GetDatabaseMetrics(ctx),
		Application: ms.GetApplicationMetrics(ctx),
		Error:       ms.GetErrorMetrics(ctx),
		Timestamp:   time.Now(),
	}
}

// ResetMetrics resets all metrics with context
func (ms *metricsService) ResetMetrics(ctx context.Context) {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	
	requestID := utils.GetRequestIDFromContext(ctx)
	ms.logger.Info(ctx, "Resetting all metrics", 
		utils.Field{Key: "request_id", Value: requestID})
	
	// Reset request metrics
	ms.totalRequests = 0
	ms.requestsByMethod = make(map[string]int64)
	ms.requestsByPath = make(map[string]int64)
	ms.requestsByStatus = make(map[int]int64)
	ms.requestLatencies = make([]time.Duration, 0, 1000)
	ms.requestStartTime = time.Now()
	
	// Reset database metrics
	ms.totalQueries = 0
	ms.successfulQueries = 0
	ms.failedQueries = 0
	ms.queryLatencies = make([]time.Duration, 0, 1000)
	ms.queriesByOp = make(map[string]int64)
	
	// Reset error metrics
	ms.totalErrors = 0
	ms.errorsByType = make(map[string]int64)
	ms.errorsByOperation = make(map[string]int64)
}

// Helper methods for calculations
func (ms *metricsService) calculateAverageLatency(latencies []time.Duration) time.Duration {
	if len(latencies) == 0 {
		return 0
	}
	
	var total time.Duration
	for _, latency := range latencies {
		total += latency
	}
	
	return total / time.Duration(len(latencies))
}

func (ms *metricsService) calculatePercentile(latencies []time.Duration, percentile float64) time.Duration {
	if len(latencies) == 0 {
		return 0
	}
	
	// Simple percentile calculation (not perfectly accurate but sufficient for monitoring)
	index := int(float64(len(latencies)) * percentile)
	if index >= len(latencies) {
		index = len(latencies) - 1
	}
	
	// Sort latencies for percentile calculation (simplified approach)
	sorted := make([]time.Duration, len(latencies))
	copy(sorted, latencies)
	
	// Simple bubble sort for small arrays
	for i := 0; i < len(sorted); i++ {
		for j := 0; j < len(sorted)-1-i; j++ {
			if sorted[j] > sorted[j+1] {
				sorted[j], sorted[j+1] = sorted[j+1], sorted[j]
			}
		}
	}
	
	return sorted[index]
}

func (ms *metricsService) countSlowQueries(latencies []time.Duration, threshold time.Duration) int64 {
	count := int64(0)
	for _, latency := range latencies {
		if latency > threshold {
			count++
		}
	}
	return count
}

func (ms *metricsService) copyStringInt64Map(original map[string]int64) map[string]int64 {
	copy := make(map[string]int64)
	for k, v := range original {
		copy[k] = v
	}
	return copy
}

func (ms *metricsService) copyIntInt64Map(original map[int]int64) map[int]int64 {
	copy := make(map[int]int64)
	for k, v := range original {
		copy[k] = v
	}
	return copy
}