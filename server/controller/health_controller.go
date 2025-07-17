package controller

import (
	"context"
	"develapar-server/config"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type HealthController struct {
	poolManager config.ConnectionPoolManager
}

type HealthResponse struct {
	Status    string                 `json:"status"`
	Timestamp time.Time              `json:"timestamp"`
	Database  DatabaseHealthStatus   `json:"database"`
	Stats     config.ConnectionStats `json:"connection_stats"`
}

type DatabaseHealthStatus struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}

// NewHealthController creates a new health controller
func NewHealthController(poolManager config.ConnectionPoolManager) *HealthController {
	return &HealthController{
		poolManager: poolManager,
	}
}

// HealthCheck godoc
// @Summary Health check endpoint
// @Description Get application and database health status
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} HealthResponse
// @Failure 503 {object} HealthResponse
// @Router /health [get]
func (hc *HealthController) HealthCheck(c *gin.Context) {
	ctx := c.Request.Context()
	
	// Create a timeout context for health check
	healthCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	response := HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now(),
		Database: DatabaseHealthStatus{
			Status: "healthy",
		},
	}

	// Check database health
	if err := hc.poolManager.HealthCheck(healthCtx); err != nil {
		response.Status = "unhealthy"
		response.Database.Status = "unhealthy"
		response.Database.Message = err.Error()
		
		c.JSON(http.StatusServiceUnavailable, response)
		return
	}

	// Get connection pool statistics
	response.Stats = hc.poolManager.GetStats(healthCtx)

	c.JSON(http.StatusOK, response)
}

// DatabaseStats godoc
// @Summary Database connection pool statistics
// @Description Get detailed database connection pool statistics
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} config.ConnectionStats
// @Router /health/database [get]
func (hc *HealthController) DatabaseStats(c *gin.Context) {
	ctx := c.Request.Context()
	
	// Create a timeout context for stats retrieval
	statsCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	stats := hc.poolManager.GetStats(statsCtx)
	c.JSON(http.StatusOK, stats)
}

// Route sets up the health check routes
func (hc *HealthController) Route(rg *gin.RouterGroup) {
	healthGroup := rg.Group("/health")
	{
		healthGroup.GET("", hc.HealthCheck)
		healthGroup.GET("/database", hc.DatabaseStats)
	}
}