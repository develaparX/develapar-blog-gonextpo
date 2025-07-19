package controller

import (
	"context"
	"develapar-server/config"
	"develapar-server/middleware"
	"develapar-server/utils"
	"net/http"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
)

type HealthController struct {
	poolManager config.ConnectionPoolManager
	logger      utils.Logger
}

// HealthResponse represents the overall health status with context information
type HealthResponse struct {
	Status       string                    `json:"status"`
	Timestamp    time.Time                 `json:"timestamp"`
	RequestID    string                    `json:"request_id,omitempty"`
	Version      string                    `json:"version"`
	Uptime       time.Duration             `json:"uptime"`
	Database     DatabaseHealthStatus      `json:"database"`
	Application  ApplicationHealthStatus   `json:"application"`
	Dependencies []DependencyHealthStatus  `json:"dependencies"`
	Stats        config.ConnectionStats    `json:"connection_stats"`
}

// DatabaseHealthStatus represents database health information
type DatabaseHealthStatus struct {
	Status       string        `json:"status"`
	Message      string        `json:"message,omitempty"`
	ResponseTime time.Duration `json:"response_time_ms"`
	LastChecked  time.Time     `json:"last_checked"`
}

// ApplicationHealthStatus represents application-level health information
type ApplicationHealthStatus struct {
	Status         string            `json:"status"`
	Memory         MemoryStats       `json:"memory"`
	Goroutines     int               `json:"goroutines"`
	StartTime      time.Time         `json:"start_time"`
	ConfigLoaded   bool              `json:"config_loaded"`
	ServicesReady  bool              `json:"services_ready"`
}

// MemoryStats represents memory usage statistics
type MemoryStats struct {
	Alloc        uint64 `json:"alloc_bytes"`
	TotalAlloc   uint64 `json:"total_alloc_bytes"`
	Sys          uint64 `json:"sys_bytes"`
	NumGC        uint32 `json:"num_gc"`
	HeapAlloc    uint64 `json:"heap_alloc_bytes"`
	HeapSys      uint64 `json:"heap_sys_bytes"`
}

// DependencyHealthStatus represents external dependency health
type DependencyHealthStatus struct {
	Name         string        `json:"name"`
	Status       string        `json:"status"`
	Message      string        `json:"message,omitempty"`
	ResponseTime time.Duration `json:"response_time_ms"`
	LastChecked  time.Time     `json:"last_checked"`
}

// DetailedHealthResponse provides comprehensive health information
type DetailedHealthResponse struct {
	HealthResponse
	Checks []HealthCheck `json:"checks"`
}

// HealthCheck represents individual health check results
type HealthCheck struct {
	Name        string        `json:"name"`
	Status      string        `json:"status"`
	Message     string        `json:"message,omitempty"`
	Duration    time.Duration `json:"duration_ms"`
	Timestamp   time.Time     `json:"timestamp"`
	Critical    bool          `json:"critical"`
}

var (
	appStartTime = time.Now()
	appVersion   = "1.0.0" // This could be set via build flags
)

// Helper functions to extract context information
func getRequestID(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if requestID, ok := ctx.Value(middleware.RequestIDKey).(string); ok {
		return requestID
	}
	return ""
}

func getUserID(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if userID, ok := ctx.Value(middleware.UserIDKey).(string); ok {
		return userID
	}
	return ""
}

// NewHealthController creates a new health controller with logger support
func NewHealthController(poolManager config.ConnectionPoolManager) *HealthController {
	logger := utils.NewDefaultLogger("health-controller")
	return &HealthController{
		poolManager: poolManager,
		logger:      logger,
	}
}

// HealthCheck godoc
// @Summary Health check endpoint
// @Description Get comprehensive application and database health status with context
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} HealthResponse
// @Failure 503 {object} HealthResponse
// @Router /health [get]
func (hc *HealthController) HealthCheck(c *gin.Context) {
	ctx := c.Request.Context()
	startTime := time.Now()
	
	// Get request ID from context if available
	requestID := getRequestID(ctx)
	
	// Create a timeout context for health check
	healthCtx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	// Log health check start with context
	hc.logger.Info(healthCtx, "Starting health check", utils.Field{
		Key:   "request_id",
		Value: requestID,
	})

	response := HealthResponse{
		Status:       "healthy",
		Timestamp:    time.Now(),
		RequestID:    requestID,
		Version:      appVersion,
		Uptime:       time.Since(appStartTime),
		Database:     hc.checkDatabaseHealth(healthCtx),
		Application:  hc.checkApplicationHealth(healthCtx),
		Dependencies: hc.checkDependenciesHealth(healthCtx),
	}

	// Get connection pool statistics with context
	response.Stats = hc.poolManager.GetStats(healthCtx)

	// Determine overall health status
	overallHealthy := response.Database.Status == "healthy" && 
					 response.Application.Status == "healthy"
	
	// Check dependencies health
	for _, dep := range response.Dependencies {
		if dep.Status != "healthy" {
			overallHealthy = false
			break
		}
	}

	if !overallHealthy {
		response.Status = "unhealthy"
		hc.logger.Warn(healthCtx, "Health check failed", 
			utils.Field{Key: "request_id", Value: requestID},
			utils.Field{Key: "database_status", Value: response.Database.Status},
			utils.Field{Key: "application_status", Value: response.Application.Status},
		)
		c.JSON(http.StatusServiceUnavailable, response)
		return
	}

	// Log successful health check
	duration := time.Since(startTime)
	hc.logger.Info(healthCtx, "Health check completed successfully", 
		utils.Field{Key: "request_id", Value: requestID},
		utils.Field{Key: "duration_ms", Value: duration.Milliseconds()},
	)

	c.JSON(http.StatusOK, response)
}

// checkDatabaseHealth performs comprehensive database health check with context
func (hc *HealthController) checkDatabaseHealth(ctx context.Context) DatabaseHealthStatus {
	startTime := time.Now()
	
	status := DatabaseHealthStatus{
		Status:      "healthy",
		LastChecked: time.Now(),
	}

	// Perform database health check with context
	if err := hc.poolManager.HealthCheck(ctx); err != nil {
		status.Status = "unhealthy"
		status.Message = err.Error()
		hc.logger.Error(ctx, "Database health check failed", err, 
			utils.Field{Key: "error", Value: err.Error()})
	}

	status.ResponseTime = time.Since(startTime)
	return status
}

// checkApplicationHealth performs application-level health checks with context
func (hc *HealthController) checkApplicationHealth(ctx context.Context) ApplicationHealthStatus {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	status := ApplicationHealthStatus{
		Status:        "healthy",
		StartTime:     appStartTime,
		Goroutines:    runtime.NumGoroutine(),
		ConfigLoaded:  true, // Assume config is loaded if we got this far
		ServicesReady: true, // Assume services are ready if we got this far
		Memory: MemoryStats{
			Alloc:      memStats.Alloc,
			TotalAlloc: memStats.TotalAlloc,
			Sys:        memStats.Sys,
			NumGC:      memStats.NumGC,
			HeapAlloc:  memStats.HeapAlloc,
			HeapSys:    memStats.HeapSys,
		},
	}

	// Check for potential memory issues
	if memStats.Alloc > 500*1024*1024 { // 500MB threshold
		status.Status = "degraded"
		hc.logger.Warn(ctx, "High memory usage detected", 
			utils.Field{Key: "alloc_mb", Value: memStats.Alloc / 1024 / 1024})
	}

	// Check for too many goroutines
	if status.Goroutines > 1000 {
		status.Status = "degraded"
		hc.logger.Warn(ctx, "High goroutine count detected", 
			utils.Field{Key: "goroutines", Value: status.Goroutines})
	}

	return status
}

// checkDependenciesHealth checks external dependencies health with context
func (hc *HealthController) checkDependenciesHealth(ctx context.Context) []DependencyHealthStatus {
	var dependencies []DependencyHealthStatus

	// Check database connection pool as a dependency
	dbDep := hc.checkDatabaseDependency(ctx)
	dependencies = append(dependencies, dbDep)

	// Add more dependency checks here as needed
	// For example: Redis, external APIs, message queues, etc.

	return dependencies
}

// checkDatabaseDependency checks database as a dependency with context
func (hc *HealthController) checkDatabaseDependency(ctx context.Context) DependencyHealthStatus {
	startTime := time.Now()
	
	dep := DependencyHealthStatus{
		Name:        "PostgreSQL Database",
		Status:      "healthy",
		LastChecked: time.Now(),
	}

	// Check database connectivity with context
	if err := hc.poolManager.HealthCheck(ctx); err != nil {
		dep.Status = "unhealthy"
		dep.Message = err.Error()
	} else {
		// Check connection pool stats for additional health indicators
		stats := hc.poolManager.GetStats(ctx)
		if stats.OpenConnections == 0 {
			dep.Status = "unhealthy"
			dep.Message = "No open database connections"
		} else if float64(stats.InUseConnections)/float64(stats.OpenConnections) > 0.9 {
			dep.Status = "degraded"
			dep.Message = "High connection pool utilization"
		}
	}

	dep.ResponseTime = time.Since(startTime)
	return dep
}

// DetailedHealthCheck godoc
// @Summary Detailed health check endpoint
// @Description Get comprehensive health information with individual check results
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} DetailedHealthResponse
// @Failure 503 {object} DetailedHealthResponse
// @Router /health/detailed [get]
func (hc *HealthController) DetailedHealthCheck(c *gin.Context) {
	ctx := c.Request.Context()
	requestID := getRequestID(ctx)
	
	// Create a timeout context for detailed health check
	healthCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	// Perform basic health check
	basicResponse := hc.getBasicHealthResponse(healthCtx, requestID)
	
	// Perform detailed checks
	checks := hc.performDetailedChecks(healthCtx)
	
	response := DetailedHealthResponse{
		HealthResponse: basicResponse,
		Checks:         checks,
	}

	// Determine status based on critical checks
	for _, check := range checks {
		if check.Critical && check.Status != "healthy" {
			response.Status = "unhealthy"
			break
		}
	}

	statusCode := http.StatusOK
	if response.Status == "unhealthy" {
		statusCode = http.StatusServiceUnavailable
	}

	c.JSON(statusCode, response)
}

// getBasicHealthResponse creates a basic health response with context
func (hc *HealthController) getBasicHealthResponse(ctx context.Context, requestID string) HealthResponse {
	return HealthResponse{
		Status:       "healthy",
		Timestamp:    time.Now(),
		RequestID:    requestID,
		Version:      appVersion,
		Uptime:       time.Since(appStartTime),
		Database:     hc.checkDatabaseHealth(ctx),
		Application:  hc.checkApplicationHealth(ctx),
		Dependencies: hc.checkDependenciesHealth(ctx),
		Stats:        hc.poolManager.GetStats(ctx),
	}
}

// performDetailedChecks runs individual health checks with context
func (hc *HealthController) performDetailedChecks(ctx context.Context) []HealthCheck {
	var checks []HealthCheck

	// Database connectivity check
	checks = append(checks, hc.checkDatabaseConnectivity(ctx))
	
	// Database query performance check
	checks = append(checks, hc.checkDatabasePerformance(ctx))
	
	// Connection pool health check
	checks = append(checks, hc.checkConnectionPoolHealth(ctx))
	
	// Memory usage check
	checks = append(checks, hc.checkMemoryUsage(ctx))
	
	// Goroutine count check
	checks = append(checks, hc.checkGoroutineCount(ctx))

	return checks
}

// checkDatabaseConnectivity performs database connectivity check with context
func (hc *HealthController) checkDatabaseConnectivity(ctx context.Context) HealthCheck {
	startTime := time.Now()
	
	check := HealthCheck{
		Name:      "Database Connectivity",
		Timestamp: time.Now(),
		Critical:  true,
	}

	if err := hc.poolManager.HealthCheck(ctx); err != nil {
		check.Status = "unhealthy"
		check.Message = err.Error()
	} else {
		check.Status = "healthy"
		check.Message = "Database connection successful"
	}

	check.Duration = time.Since(startTime)
	return check
}

// checkDatabasePerformance performs database performance check with context
func (hc *HealthController) checkDatabasePerformance(ctx context.Context) HealthCheck {
	startTime := time.Now()
	
	check := HealthCheck{
		Name:      "Database Performance",
		Timestamp: time.Now(),
		Critical:  false,
	}

	// Create a timeout context for performance check
	perfCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := hc.poolManager.HealthCheck(perfCtx); err != nil {
		check.Status = "unhealthy"
		check.Message = err.Error()
	} else {
		duration := time.Since(startTime)
		if duration > 2*time.Second {
			check.Status = "degraded"
			check.Message = "Slow database response"
		} else {
			check.Status = "healthy"
			check.Message = "Database performance normal"
		}
	}

	check.Duration = time.Since(startTime)
	return check
}

// checkConnectionPoolHealth performs connection pool health check with context
func (hc *HealthController) checkConnectionPoolHealth(ctx context.Context) HealthCheck {
	startTime := time.Now()
	
	check := HealthCheck{
		Name:      "Connection Pool Health",
		Timestamp: time.Now(),
		Critical:  false,
	}

	stats := hc.poolManager.GetStats(ctx)
	
	if stats.OpenConnections == 0 {
		check.Status = "unhealthy"
		check.Message = "No open connections"
	} else if float64(stats.InUseConnections)/float64(stats.OpenConnections) > 0.9 {
		check.Status = "degraded"
		check.Message = "High connection pool utilization"
	} else {
		check.Status = "healthy"
		check.Message = "Connection pool healthy"
	}

	check.Duration = time.Since(startTime)
	return check
}

// checkMemoryUsage performs memory usage check with context
func (hc *HealthController) checkMemoryUsage(ctx context.Context) HealthCheck {
	startTime := time.Now()
	
	check := HealthCheck{
		Name:      "Memory Usage",
		Timestamp: time.Now(),
		Critical:  false,
	}

	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	allocMB := memStats.Alloc / 1024 / 1024
	
	if allocMB > 1000 { // 1GB threshold
		check.Status = "unhealthy"
		check.Message = "High memory usage"
	} else if allocMB > 500 { // 500MB threshold
		check.Status = "degraded"
		check.Message = "Elevated memory usage"
	} else {
		check.Status = "healthy"
		check.Message = "Memory usage normal"
	}

	check.Duration = time.Since(startTime)
	return check
}

// checkGoroutineCount performs goroutine count check with context
func (hc *HealthController) checkGoroutineCount(ctx context.Context) HealthCheck {
	startTime := time.Now()
	
	check := HealthCheck{
		Name:      "Goroutine Count",
		Timestamp: time.Now(),
		Critical:  false,
	}

	goroutines := runtime.NumGoroutine()
	
	if goroutines > 2000 {
		check.Status = "unhealthy"
		check.Message = "Too many goroutines"
	} else if goroutines > 1000 {
		check.Status = "degraded"
		check.Message = "High goroutine count"
	} else {
		check.Status = "healthy"
		check.Message = "Goroutine count normal"
	}

	check.Duration = time.Since(startTime)
	return check
}

// DatabaseStats godoc
// @Summary Database connection pool statistics
// @Description Get detailed database connection pool statistics with context
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} config.ConnectionStats
// @Router /health/database [get]
func (hc *HealthController) DatabaseStats(c *gin.Context) {
	ctx := c.Request.Context()
	requestID := getRequestID(ctx)
	
	// Create a timeout context for stats retrieval
	statsCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	hc.logger.Info(statsCtx, "Retrieving database statistics", 
		utils.Field{Key: "request_id", Value: requestID})

	stats := hc.poolManager.GetStats(statsCtx)
	c.JSON(http.StatusOK, stats)
}

// Route sets up the health check routes
func (hc *HealthController) Route(rg *gin.RouterGroup) {
	healthGroup := rg.Group("/health")
	{
		healthGroup.GET("", hc.HealthCheck)
		healthGroup.GET("/detailed", hc.DetailedHealthCheck)
		healthGroup.GET("/database", hc.DatabaseStats)
	}
}