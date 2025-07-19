package controller

import (
	"context"
	"develapar-server/service"
	"develapar-server/utils"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockMetricsService is a mock implementation of MetricsService
type MockMetricsService struct {
	mock.Mock
}

func (m *MockMetricsService) RecordRequest(ctx context.Context, method, path string, statusCode int, duration time.Duration) {
	m.Called(ctx, method, path, statusCode, duration)
}

func (m *MockMetricsService) RecordError(ctx context.Context, errorType, operation string) {
	m.Called(ctx, errorType, operation)
}

func (m *MockMetricsService) RecordDatabaseQuery(ctx context.Context, operation string, duration time.Duration, success bool) {
	m.Called(ctx, operation, duration, success)
}

func (m *MockMetricsService) RecordConnectionPoolStats(ctx context.Context, open, inUse, idle int) {
	m.Called(ctx, open, inUse, idle)
}

func (m *MockMetricsService) RecordMemoryUsage(ctx context.Context, allocBytes, sysBytes uint64) {
	m.Called(ctx, allocBytes, sysBytes)
}

func (m *MockMetricsService) RecordGoroutineCount(ctx context.Context, count int) {
	m.Called(ctx, count)
}

func (m *MockMetricsService) GetRequestMetrics(ctx context.Context) service.RequestMetrics {
	args := m.Called(ctx)
	return args.Get(0).(service.RequestMetrics)
}

func (m *MockMetricsService) GetDatabaseMetrics(ctx context.Context) service.DatabaseMetrics {
	args := m.Called(ctx)
	return args.Get(0).(service.DatabaseMetrics)
}

func (m *MockMetricsService) GetApplicationMetrics(ctx context.Context) service.ApplicationMetrics {
	args := m.Called(ctx)
	return args.Get(0).(service.ApplicationMetrics)
}

func (m *MockMetricsService) GetErrorMetrics(ctx context.Context) service.ErrorMetrics {
	args := m.Called(ctx)
	return args.Get(0).(service.ErrorMetrics)
}

func (m *MockMetricsService) GetAllMetrics(ctx context.Context) service.AllMetrics {
	args := m.Called(ctx)
	return args.Get(0).(service.AllMetrics)
}

func (m *MockMetricsService) ResetMetrics(ctx context.Context) {
	m.Called(ctx)
}

func TestMetricsController_GetAllMetrics(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	mockService := new(MockMetricsService)
	controller := &MetricsController{
		metricsService: mockService,
		logger:         utils.NewDefaultLogger("test"),
	}

	// Mock data
	expectedMetrics := service.AllMetrics{
		Request: service.RequestMetrics{
			TotalRequests:     100,
			RequestsPerSecond: 10.5,
			AverageLatency:    50 * time.Millisecond,
			LastUpdated:       time.Now(),
		},
		Database: service.DatabaseMetrics{
			TotalQueries:      50,
			SuccessfulQueries: 48,
			FailedQueries:     2,
			AverageQueryTime:  25 * time.Millisecond,
			LastUpdated:       time.Now(),
		},
		Application: service.ApplicationMetrics{
			GoroutineCount: 25,
			Uptime:         time.Hour,
			StartTime:      time.Now().Add(-time.Hour),
			LastUpdated:    time.Now(),
		},
		Error: service.ErrorMetrics{
			TotalErrors: 5,
			ErrorRate:   5.0,
			LastUpdated: time.Now(),
		},
		Timestamp: time.Now(),
	}

	mockService.On("GetAllMetrics", mock.AnythingOfType("*context.timerCtx")).Return(expectedMetrics)

	// Create request
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req, _ := http.NewRequest("GET", "/metrics", nil)
	c.Request = req

	// Execute
	controller.GetAllMetrics(c)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	
	var response MetricsResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "success", response.Status)
	assert.NotNil(t, response.Data)

	mockService.AssertExpectations(t)
}

func TestMetricsController_GetRequestMetrics(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	mockService := new(MockMetricsService)
	controller := &MetricsController{
		metricsService: mockService,
		logger:         utils.NewDefaultLogger("test"),
	}

	// Mock data
	expectedMetrics := service.RequestMetrics{
		TotalRequests:     100,
		RequestsPerSecond: 10.5,
		AverageLatency:    50 * time.Millisecond,
		P95Latency:        100 * time.Millisecond,
		P99Latency:        200 * time.Millisecond,
		RequestsByMethod:  map[string]int64{"GET": 60, "POST": 40},
		RequestsByStatus:  map[int]int64{200: 90, 404: 5, 500: 5},
		LastUpdated:       time.Now(),
	}

	mockService.On("GetRequestMetrics", mock.AnythingOfType("*context.timerCtx")).Return(expectedMetrics)

	// Create request
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req, _ := http.NewRequest("GET", "/metrics/requests", nil)
	c.Request = req

	// Execute
	controller.GetRequestMetrics(c)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	
	var response MetricsResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "success", response.Status)
	assert.NotNil(t, response.Data)

	mockService.AssertExpectations(t)
}

func TestMetricsController_ResetMetrics(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	mockService := new(MockMetricsService)
	controller := &MetricsController{
		metricsService: mockService,
		logger:         utils.NewDefaultLogger("test"),
	}

	mockService.On("ResetMetrics", mock.AnythingOfType("*context.timerCtx")).Return()

	// Create request
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req, _ := http.NewRequest("POST", "/metrics/reset", nil)
	c.Request = req

	// Execute
	controller.ResetMetrics(c)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	
	var response MetricsResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "success", response.Status)
	assert.NotNil(t, response.Data)

	mockService.AssertExpectations(t)
}

func TestMetricsController_GetMetricsSummary(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	mockService := new(MockMetricsService)
	controller := &MetricsController{
		metricsService: mockService,
		logger:         utils.NewDefaultLogger("test"),
	}

	// Mock data with healthy metrics
	expectedMetrics := service.AllMetrics{
		Request: service.RequestMetrics{
			TotalRequests:     1000,
			RequestsPerSecond: 50.0,
			AverageLatency:    100 * time.Millisecond,
			P95Latency:        200 * time.Millisecond,
		},
		Database: service.DatabaseMetrics{
			TotalQueries:      500,
			SuccessfulQueries: 495,
			FailedQueries:     5,
			AverageQueryTime:  50 * time.Millisecond,
			SlowQueries:       2,
		},
		Application: service.ApplicationMetrics{
			GoroutineCount: 50,
			Uptime:         2 * time.Hour,
			MemoryUsage: service.MemoryMetrics{
				AllocBytes: 100 * 1024 * 1024, // 100MB
			},
		},
		Error: service.ErrorMetrics{
			TotalErrors: 10,
			ErrorRate:   1.0, // 1% error rate
		},
	}

	mockService.On("GetAllMetrics", mock.AnythingOfType("*context.timerCtx")).Return(expectedMetrics)

	// Create request
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req, _ := http.NewRequest("GET", "/metrics/summary", nil)
	c.Request = req

	// Execute
	controller.GetMetricsSummary(c)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	
	var response MetricsResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "success", response.Status)
	assert.NotNil(t, response.Data)

	// Verify summary structure
	data, ok := response.Data.(map[string]interface{})
	assert.True(t, ok)
	assert.Contains(t, data, "requests")
	assert.Contains(t, data, "database")
	assert.Contains(t, data, "application")
	assert.Contains(t, data, "errors")
	assert.Contains(t, data, "health")

	// Verify health status is healthy
	health, ok := data["health"].(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, "healthy", health["status"])

	mockService.AssertExpectations(t)
}