package middleware

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInMemoryStore_Increment(t *testing.T) {
	logger := &defaultLogger{}
	store := NewInMemoryStore(logger)
	ctx := context.Background()
	
	tests := []struct {
		name     string
		key      string
		window   time.Duration
		expected int
		wantErr  bool
	}{
		{
			name:     "first increment",
			key:      "test-key-1",
			window:   time.Minute,
			expected: 1,
			wantErr:  false,
		},
		{
			name:     "second increment same key",
			key:      "test-key-1",
			window:   time.Minute,
			expected: 2,
			wantErr:  false,
		},
		{
			name:     "different key",
			key:      "test-key-2",
			window:   time.Minute,
			expected: 1,
			wantErr:  false,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			count, err := store.Increment(ctx, tt.key, tt.window)
			
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, count)
			}
		})
	}
}

func TestInMemoryStore_IncrementWithWindowExpiry(t *testing.T) {
	logger := &defaultLogger{}
	store := NewInMemoryStore(logger)
	ctx := context.Background()
	
	key := "test-key"
	window := 100 * time.Millisecond
	
	// First increment
	count1, err := store.Increment(ctx, key, window)
	require.NoError(t, err)
	assert.Equal(t, 1, count1)
	
	// Second increment within window
	count2, err := store.Increment(ctx, key, window)
	require.NoError(t, err)
	assert.Equal(t, 2, count2)
	
	// Wait for window to expire
	time.Sleep(150 * time.Millisecond)
	
	// Third increment after window expiry should reset
	count3, err := store.Increment(ctx, key, window)
	require.NoError(t, err)
	assert.Equal(t, 1, count3)
}

func TestInMemoryStore_Get(t *testing.T) {
	logger := &defaultLogger{}
	store := NewInMemoryStore(logger)
	ctx := context.Background()
	
	key := "test-key"
	window := time.Minute
	
	// Get non-existent key
	count, err := store.Get(ctx, key)
	require.NoError(t, err)
	assert.Equal(t, 0, count)
	
	// Increment and then get
	_, err = store.Increment(ctx, key, window)
	require.NoError(t, err)
	
	count, err = store.Get(ctx, key)
	require.NoError(t, err)
	assert.Equal(t, 1, count)
}

func TestInMemoryStore_Reset(t *testing.T) {
	logger := &defaultLogger{}
	store := NewInMemoryStore(logger)
	ctx := context.Background()
	
	key := "test-key"
	window := time.Minute
	
	// Increment first
	_, err := store.Increment(ctx, key, window)
	require.NoError(t, err)
	
	// Verify count
	count, err := store.Get(ctx, key)
	require.NoError(t, err)
	assert.Equal(t, 1, count)
	
	// Reset
	err = store.Reset(ctx, key)
	require.NoError(t, err)
	
	// Verify reset
	count, err = store.Get(ctx, key)
	require.NoError(t, err)
	assert.Equal(t, 0, count)
}

func TestInMemoryStore_Cleanup(t *testing.T) {
	logger := &defaultLogger{}
	store := NewInMemoryStore(logger)
	ctx := context.Background()
	
	key1 := "test-key-1"
	key2 := "test-key-2"
	shortWindow := 50 * time.Millisecond
	longWindow := time.Minute
	
	// Create entries with different windows
	_, err := store.Increment(ctx, key1, shortWindow)
	require.NoError(t, err)
	
	_, err = store.Increment(ctx, key2, longWindow)
	require.NoError(t, err)
	
	// Wait for short window to expire
	time.Sleep(100 * time.Millisecond)
	
	// Cleanup
	err = store.Cleanup(ctx)
	require.NoError(t, err)
	
	// Verify expired entry is removed
	count1, err := store.Get(ctx, key1)
	require.NoError(t, err)
	assert.Equal(t, 0, count1)
	
	// Verify non-expired entry remains
	count2, err := store.Get(ctx, key2)
	require.NoError(t, err)
	assert.Equal(t, 1, count2)
}

func TestInMemoryStore_ContextCancellation(t *testing.T) {
	logger := &defaultLogger{}
	store := NewInMemoryStore(logger)
	
	// Create cancelled context
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	
	key := "test-key"
	window := time.Minute
	
	// Test increment with cancelled context
	_, err := store.Increment(ctx, key, window)
	assert.Error(t, err)
	assert.Equal(t, context.Canceled, err)
	
	// Test get with cancelled context
	_, err = store.Get(ctx, key)
	assert.Error(t, err)
	assert.Equal(t, context.Canceled, err)
	
	// Test reset with cancelled context
	err = store.Reset(ctx, key)
	assert.Error(t, err)
	assert.Equal(t, context.Canceled, err)
	
	// Test cleanup with cancelled context
	err = store.Cleanup(ctx)
	assert.Error(t, err)
	assert.Equal(t, context.Canceled, err)
}

func TestSlidingWindowRateLimiter_Allow(t *testing.T) {
	logger := &defaultLogger{}
	store := NewInMemoryStore(logger)
	limiter := NewSlidingWindowRateLimiter(store, logger)
	ctx := context.Background()
	
	key := "test-key"
	limit := 3
	window := time.Minute
	
	// Test requests within limit
	for i := 1; i <= limit; i++ {
		allowed, err := limiter.Allow(ctx, key, limit, window)
		require.NoError(t, err)
		assert.True(t, allowed, "Request %d should be allowed", i)
	}
	
	// Test request exceeding limit
	allowed, err := limiter.Allow(ctx, key, limit, window)
	require.NoError(t, err)
	assert.False(t, allowed, "Request exceeding limit should be denied")
}

func TestSlidingWindowRateLimiter_AllowWithWindowReset(t *testing.T) {
	logger := &defaultLogger{}
	store := NewInMemoryStore(logger)
	limiter := NewSlidingWindowRateLimiter(store, logger)
	ctx := context.Background()
	
	key := "test-key"
	limit := 2
	window := 100 * time.Millisecond
	
	// Use up the limit
	allowed1, err := limiter.Allow(ctx, key, limit, window)
	require.NoError(t, err)
	assert.True(t, allowed1)
	
	allowed2, err := limiter.Allow(ctx, key, limit, window)
	require.NoError(t, err)
	assert.True(t, allowed2)
	
	// Should be denied
	allowed3, err := limiter.Allow(ctx, key, limit, window)
	require.NoError(t, err)
	assert.False(t, allowed3)
	
	// Wait for window to reset
	time.Sleep(150 * time.Millisecond)
	
	// Should be allowed again
	allowed4, err := limiter.Allow(ctx, key, limit, window)
	require.NoError(t, err)
	assert.True(t, allowed4)
}

func TestSlidingWindowRateLimiter_Reset(t *testing.T) {
	logger := &defaultLogger{}
	store := NewInMemoryStore(logger)
	limiter := NewSlidingWindowRateLimiter(store, logger)
	ctx := context.Background()
	
	key := "test-key"
	limit := 2
	window := time.Minute
	
	// Use up the limit
	_, err := limiter.Allow(ctx, key, limit, window)
	require.NoError(t, err)
	
	_, err = limiter.Allow(ctx, key, limit, window)
	require.NoError(t, err)
	
	// Should be denied
	allowed, err := limiter.Allow(ctx, key, limit, window)
	require.NoError(t, err)
	assert.False(t, allowed)
	
	// Reset
	err = limiter.Reset(ctx, key)
	require.NoError(t, err)
	
	// Should be allowed again
	allowed, err = limiter.Allow(ctx, key, limit, window)
	require.NoError(t, err)
	assert.True(t, allowed)
}

func TestSlidingWindowRateLimiter_GetStats(t *testing.T) {
	logger := &defaultLogger{}
	store := NewInMemoryStore(logger)
	limiter := NewSlidingWindowRateLimiter(store, logger)
	
	// Create context with request ID
	ctx := context.WithValue(context.Background(), "request_id", "test-request-123")
	
	key := "test-key"
	limit := 5
	window := time.Minute
	
	// Make some requests
	_, err := limiter.Allow(ctx, key, limit, window)
	require.NoError(t, err)
	
	_, err = limiter.Allow(ctx, key, limit, window)
	require.NoError(t, err)
	
	// Get stats
	stats, err := limiter.GetStats(ctx, key)
	require.NoError(t, err)
	require.NotNil(t, stats)
	
	assert.Equal(t, key, stats.Key)
	assert.Equal(t, 2, stats.Count)
	assert.Equal(t, "test-request-123", stats.RequestID)
	assert.True(t, stats.Remaining >= 0)
}

func TestSlidingWindowRateLimiter_ContextCancellation(t *testing.T) {
	logger := &defaultLogger{}
	store := NewInMemoryStore(logger)
	limiter := NewSlidingWindowRateLimiter(store, logger)
	
	// Create cancelled context
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	
	key := "test-key"
	limit := 5
	window := time.Minute
	
	// Test Allow with cancelled context
	_, err := limiter.Allow(ctx, key, limit, window)
	assert.Error(t, err)
	
	// Test GetRemaining with cancelled context
	_, err = limiter.GetRemaining(ctx, key)
	assert.Error(t, err)
	
	// Test Reset with cancelled context
	err = limiter.Reset(ctx, key)
	assert.Error(t, err)
	
	// Test GetStats with cancelled context
	_, err = limiter.GetStats(ctx, key)
	assert.Error(t, err)
	
	// Test Cleanup with cancelled context
	err = limiter.Cleanup(ctx)
	assert.Error(t, err)
}

func TestSlidingWindowRateLimiter_ContextTimeout(t *testing.T) {
	logger := &defaultLogger{}
	store := NewInMemoryStore(logger)
	limiter := NewSlidingWindowRateLimiter(store, logger)
	
	// Create context with very short timeout
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
	defer cancel()
	
	// Wait for timeout
	time.Sleep(1 * time.Millisecond)
	
	key := "test-key"
	limit := 5
	window := time.Minute
	
	// Test Allow with timed out context
	_, err := limiter.Allow(ctx, key, limit, window)
	assert.Error(t, err)
}

func TestSlidingWindowRateLimiter_Cleanup(t *testing.T) {
	logger := &defaultLogger{}
	store := NewInMemoryStore(logger)
	limiter := NewSlidingWindowRateLimiter(store, logger)
	ctx := context.Background()
	
	key1 := "test-key-1"
	key2 := "test-key-2"
	limit := 5
	shortWindow := 50 * time.Millisecond
	longWindow := time.Minute
	
	// Create entries with different windows
	_, err := limiter.Allow(ctx, key1, limit, shortWindow)
	require.NoError(t, err)
	
	_, err = limiter.Allow(ctx, key2, limit, longWindow)
	require.NoError(t, err)
	
	// Wait for short window to expire
	time.Sleep(100 * time.Millisecond)
	
	// Cleanup
	err = limiter.Cleanup(ctx)
	require.NoError(t, err)
	
	// Verify cleanup worked by checking if expired entries are gone
	// This is verified indirectly through the store's behavior
	count1, err := store.Get(ctx, key1)
	require.NoError(t, err)
	assert.Equal(t, 0, count1) // Should be 0 due to expiry
	
	count2, err := store.Get(ctx, key2)
	require.NoError(t, err)
	assert.Equal(t, 1, count2) // Should still be 1
}

// Tests for RateLimitMiddleware

func TestRateLimitConfig_Default(t *testing.T) {
	config := DefaultRateLimitConfig()
	
	assert.Equal(t, 100, config.DefaultLimit)
	assert.Equal(t, time.Hour, config.DefaultWindow)
	assert.Equal(t, 1000, config.AuthenticatedLimit)
	assert.Equal(t, time.Hour, config.AuthenticatedWindow)
	assert.Equal(t, 100, config.AnonymousLimit)
	assert.Equal(t, time.Hour, config.AnonymousWindow)
	assert.Contains(t, config.SkipPaths, "/health")
	assert.Contains(t, config.SkipPaths, "/metrics")
	assert.True(t, config.IncludeHeaders)
	assert.Equal(t, "ip", config.KeyStrategy)
}

func TestRateLimitMiddleware_GenerateKey(t *testing.T) {
	logger := &defaultLogger{}
	store := NewInMemoryStore(logger)
	limiter := NewSlidingWindowRateLimiter(store, logger)
	
	tests := []struct {
		name        string
		keyStrategy string
		clientIP    string
		userID      string
		expected    string
	}{
		{
			name:        "IP strategy with no user",
			keyStrategy: "ip",
			clientIP:    "192.168.1.1",
			userID:      "",
			expected:    "ip:192.168.1.1",
		},
		{
			name:        "IP strategy with user",
			keyStrategy: "ip",
			clientIP:    "192.168.1.1",
			userID:      "user123",
			expected:    "ip:192.168.1.1",
		},
		{
			name:        "User strategy with user",
			keyStrategy: "user",
			clientIP:    "192.168.1.1",
			userID:      "user123",
			expected:    "user:user123",
		},
		{
			name:        "User strategy without user",
			keyStrategy: "user",
			clientIP:    "192.168.1.1",
			userID:      "",
			expected:    "ip:192.168.1.1",
		},
		{
			name:        "IP+User strategy with user",
			keyStrategy: "ip_user",
			clientIP:    "192.168.1.1",
			userID:      "user123",
			expected:    "ip_user:192.168.1.1:user123",
		},
		{
			name:        "IP+User strategy without user",
			keyStrategy: "ip_user",
			clientIP:    "192.168.1.1",
			userID:      "",
			expected:    "ip:192.168.1.1",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := DefaultRateLimitConfig()
			config.KeyStrategy = tt.keyStrategy
			
			middleware := NewRateLimitMiddleware(limiter, config, logger)
			ctx := context.Background()
			
			key := middleware.generateKey(ctx, tt.clientIP, tt.userID)
			assert.Equal(t, tt.expected, key)
		})
	}
}

func TestRateLimitMiddleware_ShouldSkipPath(t *testing.T) {
	logger := &defaultLogger{}
	store := NewInMemoryStore(logger)
	limiter := NewSlidingWindowRateLimiter(store, logger)
	
	config := DefaultRateLimitConfig()
	config.SkipPaths = []string{"/health", "/metrics", "/api/v1/status"}
	
	middleware := NewRateLimitMiddleware(limiter, config, logger)
	
	tests := []struct {
		path     string
		expected bool
	}{
		{"/health", true},
		{"/metrics", true},
		{"/api/v1/status", true},
		{"/api/v1/users", false},
		{"/api/v1/articles", false},
		{"/healthcheck", false}, // Similar but not exact match
	}
	
	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			result := middleware.shouldSkipPath(tt.path)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestRateLimitMiddleware_GetLimitsForUser(t *testing.T) {
	logger := &defaultLogger{}
	store := NewInMemoryStore(logger)
	limiter := NewSlidingWindowRateLimiter(store, logger)
	
	config := DefaultRateLimitConfig()
	config.AuthenticatedLimit = 1000
	config.AuthenticatedWindow = time.Hour
	config.AnonymousLimit = 100
	config.AnonymousWindow = time.Minute
	
	middleware := NewRateLimitMiddleware(limiter, config, logger)
	
	// Test authenticated user
	limit, window := middleware.getLimitsForUser("user123")
	assert.Equal(t, 1000, limit)
	assert.Equal(t, time.Hour, window)
	
	// Test anonymous user
	limit, window = middleware.getLimitsForUser("")
	assert.Equal(t, 100, limit)
	assert.Equal(t, time.Minute, window)
}

// Tests for RateLimitMetrics

func TestRateLimitMetrics_IncrementOperations(t *testing.T) {
	metrics := NewRateLimitMetrics()
	
	// Test initial state
	snapshot := metrics.GetSnapshot()
	assert.Equal(t, int64(0), snapshot.TotalRequests)
	assert.Equal(t, int64(0), snapshot.AllowedRequests)
	assert.Equal(t, int64(0), snapshot.BlockedRequests)
	assert.Equal(t, int64(0), snapshot.ErrorCount)
	
	// Test increments
	metrics.IncrementTotal()
	metrics.IncrementAllowed()
	metrics.IncrementBlocked("test-key")
	metrics.IncrementError()
	metrics.IncrementKeyRequests("test-key")
	
	snapshot = metrics.GetSnapshot()
	assert.Equal(t, int64(1), snapshot.TotalRequests)
	assert.Equal(t, int64(1), snapshot.AllowedRequests)
	assert.Equal(t, int64(1), snapshot.BlockedRequests)
	assert.Equal(t, int64(1), snapshot.ErrorCount)
	assert.Equal(t, int64(1), snapshot.ViolationsByKey["test-key"])
	assert.Equal(t, int64(1), snapshot.RequestsByKey["test-key"])
}

func TestRateLimitMetrics_SetOperations(t *testing.T) {
	metrics := NewRateLimitMetrics()
	
	now := time.Now()
	metrics.SetActiveKeys(10)
	metrics.SetLastCleanup(now)
	
	snapshot := metrics.GetSnapshot()
	assert.Equal(t, 10, snapshot.ActiveKeys)
	assert.Equal(t, now, snapshot.LastCleanup)
}

func TestRateLimitMetrics_Reset(t *testing.T) {
	metrics := NewRateLimitMetrics()
	
	// Add some data
	metrics.IncrementTotal()
	metrics.IncrementAllowed()
	metrics.IncrementBlocked("test-key")
	metrics.IncrementError()
	metrics.SetActiveKeys(5)
	
	// Verify data exists
	snapshot := metrics.GetSnapshot()
	assert.Equal(t, int64(1), snapshot.TotalRequests)
	assert.Equal(t, int64(1), snapshot.AllowedRequests)
	assert.Equal(t, int64(1), snapshot.BlockedRequests)
	assert.Equal(t, int64(1), snapshot.ErrorCount)
	assert.Equal(t, 5, snapshot.ActiveKeys)
	
	// Reset
	metrics.Reset()
	
	// Verify reset
	snapshot = metrics.GetSnapshot()
	assert.Equal(t, int64(0), snapshot.TotalRequests)
	assert.Equal(t, int64(0), snapshot.AllowedRequests)
	assert.Equal(t, int64(0), snapshot.BlockedRequests)
	assert.Equal(t, int64(0), snapshot.ErrorCount)
	assert.Equal(t, 0, snapshot.ActiveKeys)
	assert.Empty(t, snapshot.ViolationsByKey)
	assert.Empty(t, snapshot.RequestsByKey)
}

// Tests for RateLimitMonitor

func TestRateLimitMonitor_LogViolation(t *testing.T) {
	logger := &defaultLogger{}
	store := NewInMemoryStore(logger)
	limiter := NewSlidingWindowRateLimiter(store, logger)
	monitor := NewRateLimitMonitor(limiter, logger)
	
	ctx := context.Background()
	key := "test-key"
	clientIP := "192.168.1.1"
	userID := "user123"
	limit := 10
	window := time.Minute
	count := 15
	
	// Log violation
	monitor.LogViolation(ctx, key, clientIP, userID, limit, window, count)
	
	// Check metrics
	snapshot := monitor.metrics.GetSnapshot()
	assert.Equal(t, int64(1), snapshot.BlockedRequests)
	assert.Equal(t, int64(1), snapshot.ViolationsByKey[key])
}

func TestRateLimitMonitor_LogAllowed(t *testing.T) {
	logger := &defaultLogger{}
	store := NewInMemoryStore(logger)
	limiter := NewSlidingWindowRateLimiter(store, logger)
	monitor := NewRateLimitMonitor(limiter, logger)
	
	ctx := context.Background()
	key := "test-key"
	clientIP := "192.168.1.1"
	userID := "user123"
	count := 5
	limit := 10
	remaining := 5
	
	// Log allowed request
	monitor.LogAllowed(ctx, key, clientIP, userID, count, limit, remaining)
	
	// Check metrics
	snapshot := monitor.metrics.GetSnapshot()
	assert.Equal(t, int64(1), snapshot.AllowedRequests)
	assert.Equal(t, int64(1), snapshot.RequestsByKey[key])
}

func TestRateLimitMonitor_LogError(t *testing.T) {
	logger := &defaultLogger{}
	store := NewInMemoryStore(logger)
	limiter := NewSlidingWindowRateLimiter(store, logger)
	monitor := NewRateLimitMonitor(limiter, logger)
	
	ctx := context.Background()
	key := "test-key"
	clientIP := "192.168.1.1"
	userID := "user123"
	err := fmt.Errorf("test error")
	operation := "test_operation"
	
	// Log error
	monitor.LogError(ctx, key, clientIP, userID, err, operation)
	
	// Check metrics
	snapshot := monitor.metrics.GetSnapshot()
	assert.Equal(t, int64(1), snapshot.ErrorCount)
}

func TestRateLimitMonitor_GetMetrics(t *testing.T) {
	logger := &defaultLogger{}
	store := NewInMemoryStore(logger)
	limiter := NewSlidingWindowRateLimiter(store, logger)
	monitor := NewRateLimitMonitor(limiter, logger)
	
	ctx := context.Background()
	
	// Add some data
	monitor.metrics.IncrementAllowed()
	monitor.metrics.IncrementBlocked("test-key")
	
	// Get metrics
	snapshot := monitor.GetMetrics(ctx)
	
	// Should include the metrics request itself
	assert.Equal(t, int64(1), snapshot.TotalRequests)
	assert.Equal(t, int64(1), snapshot.AllowedRequests)
	assert.Equal(t, int64(1), snapshot.BlockedRequests)
}

func TestRateLimitMonitor_ResetMetrics(t *testing.T) {
	logger := &defaultLogger{}
	store := NewInMemoryStore(logger)
	limiter := NewSlidingWindowRateLimiter(store, logger)
	monitor := NewRateLimitMonitor(limiter, logger)
	
	ctx := context.Background()
	
	// Add some data
	monitor.metrics.IncrementAllowed()
	monitor.metrics.IncrementBlocked("test-key")
	
	// Verify data exists
	snapshot := monitor.metrics.GetSnapshot()
	assert.Equal(t, int64(1), snapshot.AllowedRequests)
	assert.Equal(t, int64(1), snapshot.BlockedRequests)
	
	// Reset metrics
	monitor.ResetMetrics(ctx)
	
	// Verify reset
	snapshot = monitor.metrics.GetSnapshot()
	assert.Equal(t, int64(0), snapshot.AllowedRequests)
	assert.Equal(t, int64(0), snapshot.BlockedRequests)
}

func TestRateLimitMonitor_GetTopViolators(t *testing.T) {
	logger := &defaultLogger{}
	store := NewInMemoryStore(logger)
	limiter := NewSlidingWindowRateLimiter(store, logger)
	monitor := NewRateLimitMonitor(limiter, logger)
	
	violations := map[string]int64{
		"key1": 10,
		"key2": 5,
		"key3": 15,
		"key4": 2,
		"key5": 8,
	}
	
	topViolators := monitor.getTopViolators(violations, 3)
	
	// Should be sorted by violations descending
	require.Len(t, topViolators, 3)
	assert.Equal(t, "key3", topViolators[0]["key"])
	assert.Equal(t, int64(15), topViolators[0]["violations"])
	assert.Equal(t, "key1", topViolators[1]["key"])
	assert.Equal(t, int64(10), topViolators[1]["violations"])
	assert.Equal(t, "key5", topViolators[2]["key"])
	assert.Equal(t, int64(8), topViolators[2]["violations"])
}

// Tests for MonitoredRateLimitMiddleware

func TestNewMonitoredRateLimitMiddleware(t *testing.T) {
	logger := &defaultLogger{}
	store := NewInMemoryStore(logger)
	limiter := NewSlidingWindowRateLimiter(store, logger)
	config := DefaultRateLimitConfig()
	
	middleware := NewMonitoredRateLimitMiddleware(limiter, config, logger)
	
	assert.NotNil(t, middleware)
	assert.NotNil(t, middleware.RateLimitMiddleware)
	assert.NotNil(t, middleware.monitor)
	
	monitor := middleware.GetMonitor()
	assert.NotNil(t, monitor)
	assert.Equal(t, middleware.monitor, monitor)
}