package controller

import (
	"context"
	"develapar-server/service"
	"develapar-server/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// MetricsController handles metrics-related HTTP requests
type MetricsController struct {
	metricsService service.MetricsService
	logger         utils.Logger
}

// MetricsResponse represents the response structure for metrics endpoints
type MetricsResponse struct {
	Status    string                 `json:"status"`
	Timestamp time.Time              `json:"timestamp"`
	RequestID string                 `json:"request_id,omitempty"`
	Data      interface{}            `json:"data"`
}

// NewMetricsController creates a new metrics controller
func NewMetricsController(metricsService service.MetricsService) *MetricsController {
	logger := utils.NewDefaultLogger("metrics-controller")
	return &MetricsController{
		metricsService: metricsService,
		logger:         logger,
	}
}

// GetAllMetrics godoc
// @Summary Get all application metrics
// @Description Get comprehensive application metrics including request, database, application, and error metrics
// @Tags metrics
// @Accept json
// @Produce json
// @Success 200 {object} MetricsResponse
// @Failure 500 {object} MetricsResponse
// @Router /metrics [get]
func (mc *MetricsController) GetAllMetrics(c *gin.Context) {
	ctx := c.Request.Context()
	requestID := getRequestID(ctx)
	
	// Create timeout context for metrics collection
	metricsCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	mc.logger.Info(metricsCtx, "Retrieving all metrics", 
		utils.Field{Key: "request_id", Value: requestID})

	// Get all metrics with context
	metrics := mc.metricsService.GetAllMetrics(metricsCtx)

	response := MetricsResponse{
		Status:    "success",
		Timestamp: time.Now(),
		RequestID: requestID,
		Data:      metrics,
	}

	c.JSON(http.StatusOK, response)
}

// GetRequestMetrics godoc
// @Summary Get request metrics
// @Description Get HTTP request metrics including latency, throughput, and status code distribution
// @Tags metrics
// @Accept json
// @Produce json
// @Success 200 {object} MetricsResponse
// @Failure 500 {object} MetricsResponse
// @Router /metrics/requests [get]
func (mc *MetricsController) GetRequestMetrics(c *gin.Context) {
	ctx := c.Request.Context()
	requestID := getRequestID(ctx)
	
	// Create timeout context for metrics collection
	metricsCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	mc.logger.Info(metricsCtx, "Retrieving request metrics", 
		utils.Field{Key: "request_id", Value: requestID})

	// Get request metrics with context
	metrics := mc.metricsService.GetRequestMetrics(metricsCtx)

	response := MetricsResponse{
		Status:    "success",
		Timestamp: time.Now(),
		RequestID: requestID,
		Data:      metrics,
	}

	c.JSON(http.StatusOK, response)
}

// GetDatabaseMetrics godoc
// @Summary Get database metrics
// @Description Get database performance metrics including query times, connection pool stats, and error rates
// @Tags metrics
// @Accept json
// @Produce json
// @Success 200 {object} MetricsResponse
// @Failure 500 {object} MetricsResponse
// @Router /metrics/database [get]
func (mc *MetricsController) GetDatabaseMetrics(c *gin.Context) {
	ctx := c.Request.Context()
	requestID := getRequestID(ctx)
	
	// Create timeout context for metrics collection
	metricsCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	mc.logger.Info(metricsCtx, "Retrieving database metrics", 
		utils.Field{Key: "request_id", Value: requestID})

	// Get database metrics with context
	metrics := mc.metricsService.GetDatabaseMetrics(metricsCtx)

	response := MetricsResponse{
		Status:    "success",
		Timestamp: time.Now(),
		RequestID: requestID,
		Data:      metrics,
	}

	c.JSON(http.StatusOK, response)
}

// GetApplicationMetrics godoc
// @Summary Get application metrics
// @Description Get application-level metrics including memory usage, goroutine count, and uptime
// @Tags metrics
// @Accept json
// @Produce json
// @Success 200 {object} MetricsResponse
// @Failure 500 {object} MetricsResponse
// @Router /metrics/application [get]
func (mc *MetricsController) GetApplicationMetrics(c *gin.Context) {
	ctx := c.Request.Context()
	requestID := getRequestID(ctx)
	
	// Create timeout context for metrics collection
	metricsCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	mc.logger.Info(metricsCtx, "Retrieving application metrics", 
		utils.Field{Key: "request_id", Value: requestID})

	// Get application metrics with context
	metrics := mc.metricsService.GetApplicationMetrics(metricsCtx)

	response := MetricsResponse{
		Status:    "success",
		Timestamp: time.Now(),
		RequestID: requestID,
		Data:      metrics,
	}

	c.JSON(http.StatusOK, response)
}

// GetErrorMetrics godoc
// @Summary Get error metrics
// @Description Get error tracking metrics including error rates, error types, and error distribution
// @Tags metrics
// @Accept json
// @Produce json
// @Success 200 {object} MetricsResponse
// @Failure 500 {object} MetricsResponse
// @Router /metrics/errors [get]
func (mc *MetricsController) GetErrorMetrics(c *gin.Context) {
	ctx := c.Request.Context()
	requestID := getRequestID(ctx)
	
	// Create timeout context for metrics collection
	metricsCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	mc.logger.Info(metricsCtx, "Retrieving error metrics", 
		utils.Field{Key: "request_id", Value: requestID})

	// Get error metrics with context
	metrics := mc.metricsService.GetErrorMetrics(metricsCtx)

	response := MetricsResponse{
		Status:    "success",
		Timestamp: time.Now(),
		RequestID: requestID,
		Data:      metrics,
	}

	c.JSON(http.StatusOK, response)
}

// ResetMetrics godoc
// @Summary Reset all metrics
// @Description Reset all collected metrics to zero (useful for testing or periodic resets)
// @Tags metrics
// @Accept json
// @Produce json
// @Success 200 {object} MetricsResponse
// @Failure 500 {object} MetricsResponse
// @Router /metrics/reset [post]
func (mc *MetricsController) ResetMetrics(c *gin.Context) {
	ctx := c.Request.Context()
	requestID := getRequestID(ctx)
	
	// Create timeout context for metrics reset
	resetCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	mc.logger.Info(resetCtx, "Resetting all metrics", 
		utils.Field{Key: "request_id", Value: requestID})

	// Reset metrics with context
	mc.metricsService.ResetMetrics(resetCtx)

	response := MetricsResponse{
		Status:    "success",
		Timestamp: time.Now(),
		RequestID: requestID,
		Data:      map[string]string{"message": "All metrics have been reset"},
	}

	mc.logger.Info(resetCtx, "All metrics reset successfully", 
		utils.Field{Key: "request_id", Value: requestID})

	c.JSON(http.StatusOK, response)
}

// GetMetricsSummary godoc
// @Summary Get metrics summary
// @Description Get a high-level summary of key metrics for dashboard display
// @Tags metrics
// @Accept json
// @Produce json
// @Success 200 {object} MetricsResponse
// @Failure 500 {object} MetricsResponse
// @Router /metrics/summary [get]
func (mc *MetricsController) GetMetricsSummary(c *gin.Context) {
	ctx := c.Request.Context()
	requestID := getRequestID(ctx)
	
	// Create timeout context for metrics collection
	summaryCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	mc.logger.Info(summaryCtx, "Retrieving metrics summary", 
		utils.Field{Key: "request_id", Value: requestID})

	// Get all metrics and create summary
	allMetrics := mc.metricsService.GetAllMetrics(summaryCtx)
	
	summary := map[string]interface{}{
		"requests": map[string]interface{}{
			"total":              allMetrics.Request.TotalRequests,
			"requests_per_second": allMetrics.Request.RequestsPerSecond,
			"average_latency_ms": allMetrics.Request.AverageLatency.Milliseconds(),
			"p95_latency_ms":     allMetrics.Request.P95Latency.Milliseconds(),
		},
		"database": map[string]interface{}{
			"total_queries":        allMetrics.Database.TotalQueries,
			"successful_queries":   allMetrics.Database.SuccessfulQueries,
			"failed_queries":       allMetrics.Database.FailedQueries,
			"average_query_time_ms": allMetrics.Database.AverageQueryTime.Milliseconds(),
			"slow_queries":         allMetrics.Database.SlowQueries,
		},
		"application": map[string]interface{}{
			"uptime_seconds":   allMetrics.Application.Uptime.Seconds(),
			"memory_alloc_mb":  allMetrics.Application.MemoryUsage.AllocBytes / 1024 / 1024,
			"goroutine_count":  allMetrics.Application.GoroutineCount,
		},
		"errors": map[string]interface{}{
			"total_errors": allMetrics.Error.TotalErrors,
			"error_rate":   allMetrics.Error.ErrorRate,
		},
		"health": map[string]interface{}{
			"status": mc.getHealthStatus(allMetrics),
		},
	}

	response := MetricsResponse{
		Status:    "success",
		Timestamp: time.Now(),
		RequestID: requestID,
		Data:      summary,
	}

	c.JSON(http.StatusOK, response)
}

// getHealthStatus determines overall health status based on metrics
func (mc *MetricsController) getHealthStatus(metrics service.AllMetrics) string {
	// Check error rate
	if metrics.Error.ErrorRate > 10.0 { // More than 10% error rate
		return "unhealthy"
	}
	
	// Check memory usage
	memoryMB := metrics.Application.MemoryUsage.AllocBytes / 1024 / 1024
	if memoryMB > 1000 { // More than 1GB
		return "degraded"
	}
	
	// Check goroutine count
	if metrics.Application.GoroutineCount > 2000 {
		return "degraded"
	}
	
	// Check database performance
	if metrics.Database.AverageQueryTime > 2*time.Second {
		return "degraded"
	}
	
	// Check request latency
	if metrics.Request.P95Latency > 5*time.Second {
		return "degraded"
	}
	
	return "healthy"
}

// Route sets up the metrics routes
func (mc *MetricsController) Route(rg *gin.RouterGroup) {
	metricsGroup := rg.Group("/metrics")
	{
		metricsGroup.GET("", mc.GetAllMetrics)
		metricsGroup.GET("/summary", mc.GetMetricsSummary)
		metricsGroup.GET("/requests", mc.GetRequestMetrics)
		metricsGroup.GET("/database", mc.GetDatabaseMetrics)
		metricsGroup.GET("/application", mc.GetApplicationMetrics)
		metricsGroup.GET("/errors", mc.GetErrorMetrics)
		metricsGroup.POST("/reset", mc.ResetMetrics)
	}
}