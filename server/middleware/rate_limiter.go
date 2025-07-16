package middleware

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"develapar-server/utils"
)

// RateLimiter interface defines rate limiting operations with context support
type RateLimiter interface {
	Allow(ctx context.Context, key string, limit int, window time.Duration) (bool, error)
	GetRemaining(ctx context.Context, key string) (int, error)
	Reset(ctx context.Context, key string) error
	GetStats(ctx context.Context, key string) (*RateLimitStats, error)
	Cleanup(ctx context.Context) error
}

// RateLimitStore interface defines storage operations for rate limiting
type RateLimitStore interface {
	Increment(ctx context.Context, key string, window time.Duration) (int, error)
	Get(ctx context.Context, key string) (int, error)
	Reset(ctx context.Context, key string) error
	GetExpiry(ctx context.Context, key string) (time.Time, error)
	SetExpiry(ctx context.Context, key string, expiry time.Time) error
	Cleanup(ctx context.Context) error
}

// RateLimitStats contains statistics about rate limiting
type RateLimitStats struct {
	Key           string        `json:"key"`
	Count         int           `json:"count"`
	Limit         int           `json:"limit"`
	Remaining     int           `json:"remaining"`
	ResetTime     time.Time     `json:"reset_time"`
	Window        time.Duration `json:"window"`
	RequestID     string        `json:"request_id,omitempty"`
}

// RateLimitEntry represents a rate limit entry in the store
type RateLimitEntry struct {
	Count     int       `json:"count"`
	Window    time.Duration `json:"window"`
	StartTime time.Time `json:"start_time"`
	LastReset time.Time `json:"last_reset"`
}

// inMemoryStore implements RateLimitStore using in-memory storage
type inMemoryStore struct {
	mu      sync.RWMutex
	entries map[string]*RateLimitEntry
	logger  Logger
}

// NewInMemoryStore creates a new in-memory rate limit store
func NewInMemoryStore(logger Logger) RateLimitStore {
	if logger == nil {
		logger = &defaultLogger{}
	}
	
	return &inMemoryStore{
		entries: make(map[string]*RateLimitEntry),
		logger:  logger,
	}
}

// Increment increments the counter for a key within the specified window
func (s *inMemoryStore) Increment(ctx context.Context, key string, window time.Duration) (int, error) {
	// Check for context cancellation
	select {
	case <-ctx.Done():
		return 0, ctx.Err()
	default:
	}
	
	s.mu.Lock()
	defer s.mu.Unlock()
	
	now := time.Now()
	entry, exists := s.entries[key]
	
	if !exists {
		// Create new entry
		entry = &RateLimitEntry{
			Count:     1,
			Window:    window,
			StartTime: now,
			LastReset: now,
		}
		s.entries[key] = entry
		
		s.logger.Info(ctx, "Rate limit entry created", map[string]interface{}{
			"key":    key,
			"window": window,
		})
		
		return 1, nil
	}
	
	// Check if window has expired
	if now.Sub(entry.StartTime) >= entry.Window {
		// Reset the window
		entry.Count = 1
		entry.StartTime = now
		entry.LastReset = now
		entry.Window = window
		
		s.logger.Info(ctx, "Rate limit window reset", map[string]interface{}{
			"key":    key,
			"window": window,
		})
		
		return 1, nil
	}
	
	// Increment within current window
	entry.Count++
	
	return entry.Count, nil
}

// Get retrieves the current count for a key
func (s *inMemoryStore) Get(ctx context.Context, key string) (int, error) {
	// Check for context cancellation
	select {
	case <-ctx.Done():
		return 0, ctx.Err()
	default:
	}
	
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	entry, exists := s.entries[key]
	if !exists {
		return 0, nil
	}
	
	// Check if window has expired
	now := time.Now()
	if now.Sub(entry.StartTime) >= entry.Window {
		return 0, nil
	}
	
	return entry.Count, nil
}

// Reset resets the counter for a key
func (s *inMemoryStore) Reset(ctx context.Context, key string) error {
	// Check for context cancellation
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	
	s.mu.Lock()
	defer s.mu.Unlock()
	
	if entry, exists := s.entries[key]; exists {
		entry.Count = 0
		entry.StartTime = time.Now()
		entry.LastReset = time.Now()
		
		s.logger.Info(ctx, "Rate limit reset", map[string]interface{}{
			"key": key,
		})
	}
	
	return nil
}

// GetExpiry returns the expiry time for a key
func (s *inMemoryStore) GetExpiry(ctx context.Context, key string) (time.Time, error) {
	// Check for context cancellation
	select {
	case <-ctx.Done():
		return time.Time{}, ctx.Err()
	default:
	}
	
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	entry, exists := s.entries[key]
	if !exists {
		return time.Time{}, nil
	}
	
	return entry.StartTime.Add(entry.Window), nil
}

// SetExpiry sets the expiry time for a key
func (s *inMemoryStore) SetExpiry(ctx context.Context, key string, expiry time.Time) error {
	// Check for context cancellation
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	
	s.mu.Lock()
	defer s.mu.Unlock()
	
	if entry, exists := s.entries[key]; exists {
		entry.StartTime = expiry.Add(-entry.Window)
	}
	
	return nil
}

// Cleanup removes expired entries
func (s *inMemoryStore) Cleanup(ctx context.Context) error {
	// Check for context cancellation
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	
	s.mu.Lock()
	defer s.mu.Unlock()
	
	now := time.Now()
	expiredKeys := make([]string, 0)
	
	for key, entry := range s.entries {
		if now.Sub(entry.StartTime) >= entry.Window {
			expiredKeys = append(expiredKeys, key)
		}
	}
	
	for _, key := range expiredKeys {
		delete(s.entries, key)
	}
	
	if len(expiredKeys) > 0 {
		s.logger.Info(ctx, "Rate limit cleanup completed", map[string]interface{}{
			"expired_entries": len(expiredKeys),
		})
	}
	
	return nil
}

// slidingWindowRateLimiter implements RateLimiter using sliding window algorithm
type slidingWindowRateLimiter struct {
	store  RateLimitStore
	logger Logger
	wrapper utils.ErrorWrapper
}

// NewSlidingWindowRateLimiter creates a new sliding window rate limiter
func NewSlidingWindowRateLimiter(store RateLimitStore, logger Logger) RateLimiter {
	if logger == nil {
		logger = &defaultLogger{}
	}
	
	return &slidingWindowRateLimiter{
		store:   store,
		logger:  logger,
		wrapper: utils.NewErrorWrapper(),
	}
}

// Allow checks if a request is allowed within the rate limit
func (rl *slidingWindowRateLimiter) Allow(ctx context.Context, key string, limit int, window time.Duration) (bool, error) {
	// Check for context cancellation
	select {
	case <-ctx.Done():
		return false, rl.wrapper.CancellationError(ctx, "rate limit check")
	default:
	}
	
	// Increment the counter
	count, err := rl.store.Increment(ctx, key, window)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return false, rl.wrapper.TimeoutError(ctx, "rate limit check")
		}
		if ctx.Err() == context.Canceled {
			return false, rl.wrapper.CancellationError(ctx, "rate limit check")
		}
		return false, rl.wrapper.InternalError(ctx, err, "Failed to check rate limit")
	}
	
	allowed := count <= limit
	
	// Log rate limit check
	rl.logger.Info(ctx, "Rate limit check", map[string]interface{}{
		"key":     key,
		"count":   count,
		"limit":   limit,
		"allowed": allowed,
		"window":  window,
	})
	
	// Log rate limit violation
	if !allowed {
		rl.logger.Warn(ctx, "Rate limit exceeded", map[string]interface{}{
			"key":    key,
			"count":  count,
			"limit":  limit,
			"window": window,
		})
	}
	
	return allowed, nil
}

// GetRemaining returns the number of remaining requests
func (rl *slidingWindowRateLimiter) GetRemaining(ctx context.Context, key string) (int, error) {
	// Check for context cancellation
	select {
	case <-ctx.Done():
		return 0, rl.wrapper.CancellationError(ctx, "get remaining requests")
	default:
	}
	
	count, err := rl.store.Get(ctx, key)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return 0, rl.wrapper.TimeoutError(ctx, "get remaining requests")
		}
		if ctx.Err() == context.Canceled {
			return 0, rl.wrapper.CancellationError(ctx, "get remaining requests")
		}
		return 0, rl.wrapper.InternalError(ctx, err, "Failed to get remaining requests")
	}
	
	// Note: This is a simplified implementation
	// In a real sliding window, we'd need to track the limit per key
	// For now, we'll return 0 if we have any count
	if count > 0 {
		return 0, nil
	}
	
	return 100, nil // Default limit assumption
}

// Reset resets the rate limit for a key
func (rl *slidingWindowRateLimiter) Reset(ctx context.Context, key string) error {
	// Check for context cancellation
	select {
	case <-ctx.Done():
		return rl.wrapper.CancellationError(ctx, "rate limit reset")
	default:
	}
	
	err := rl.store.Reset(ctx, key)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return rl.wrapper.TimeoutError(ctx, "rate limit reset")
		}
		if ctx.Err() == context.Canceled {
			return rl.wrapper.CancellationError(ctx, "rate limit reset")
		}
		return rl.wrapper.InternalError(ctx, err, "Failed to reset rate limit")
	}
	
	rl.logger.Info(ctx, "Rate limit reset", map[string]interface{}{
		"key": key,
	})
	
	return nil
}

// GetStats returns statistics for a rate limit key
func (rl *slidingWindowRateLimiter) GetStats(ctx context.Context, key string) (*RateLimitStats, error) {
	// Check for context cancellation
	select {
	case <-ctx.Done():
		return nil, rl.wrapper.CancellationError(ctx, "get rate limit stats")
	default:
	}
	
	count, err := rl.store.Get(ctx, key)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return nil, rl.wrapper.TimeoutError(ctx, "get rate limit stats")
		}
		if ctx.Err() == context.Canceled {
			return nil, rl.wrapper.CancellationError(ctx, "get rate limit stats")
		}
		return nil, rl.wrapper.InternalError(ctx, err, "Failed to get rate limit stats")
	}
	
	expiry, err := rl.store.GetExpiry(ctx, key)
	if err != nil {
		return nil, rl.wrapper.InternalError(ctx, err, "Failed to get rate limit expiry")
	}
	
	// Extract request ID from context
	requestID := ""
	if ctx != nil {
		if rid, ok := ctx.Value("request_id").(string); ok {
			requestID = rid
		}
	}
	
	// Default values for demonstration
	limit := 100
	window := time.Minute
	
	stats := &RateLimitStats{
		Key:       key,
		Count:     count,
		Limit:     limit,
		Remaining: limit - count,
		ResetTime: expiry,
		Window:    window,
		RequestID: requestID,
	}
	
	if stats.Remaining < 0 {
		stats.Remaining = 0
	}
	
	return stats, nil
}

// Cleanup removes expired rate limit entries
func (rl *slidingWindowRateLimiter) Cleanup(ctx context.Context) error {
	// Check for context cancellation
	select {
	case <-ctx.Done():
		return rl.wrapper.CancellationError(ctx, "rate limit cleanup")
	default:
	}
	
	err := rl.store.Cleanup(ctx)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return rl.wrapper.TimeoutError(ctx, "rate limit cleanup")
		}
		if ctx.Err() == context.Canceled {
			return rl.wrapper.CancellationError(ctx, "rate limit cleanup")
		}
		return rl.wrapper.InternalError(ctx, err, "Failed to cleanup rate limits")
	}
	
	return nil
}

// RateLimitConfig holds configuration for rate limiting
type RateLimitConfig struct {
	// Default limits
	DefaultLimit       int           `json:"default_limit"`
	DefaultWindow      time.Duration `json:"default_window"`
	
	// Authenticated user limits
	AuthenticatedLimit int           `json:"authenticated_limit"`
	AuthenticatedWindow time.Duration `json:"authenticated_window"`
	
	// Anonymous user limits
	AnonymousLimit     int           `json:"anonymous_limit"`
	AnonymousWindow    time.Duration `json:"anonymous_window"`
	
	// Skip rate limiting for certain paths
	SkipPaths          []string      `json:"skip_paths"`
	
	// Headers to include in response
	IncludeHeaders     bool          `json:"include_headers"`
	
	// Key generation strategy
	KeyStrategy        string        `json:"key_strategy"` // "ip", "user", "ip_user"
}

// DefaultRateLimitConfig returns default rate limiting configuration
func DefaultRateLimitConfig() *RateLimitConfig {
	return &RateLimitConfig{
		DefaultLimit:        100,
		DefaultWindow:       time.Hour,
		AuthenticatedLimit:  1000,
		AuthenticatedWindow: time.Hour,
		AnonymousLimit:      100,
		AnonymousWindow:     time.Hour,
		SkipPaths:          []string{"/health", "/metrics"},
		IncludeHeaders:     true,
		KeyStrategy:        "ip",
	}
}

// RateLimitMiddleware provides rate limiting functionality for HTTP requests
type RateLimitMiddleware struct {
	limiter RateLimiter
	config  *RateLimitConfig
	logger  Logger
	wrapper utils.ErrorWrapper
}

// NewRateLimitMiddleware creates a new rate limiting middleware
func NewRateLimitMiddleware(limiter RateLimiter, config *RateLimitConfig, logger Logger) *RateLimitMiddleware {
	if config == nil {
		config = DefaultRateLimitConfig()
	}
	if logger == nil {
		logger = &defaultLogger{}
	}
	
	return &RateLimitMiddleware{
		limiter: limiter,
		config:  config,
		logger:  logger,
		wrapper: utils.NewErrorWrapper(),
	}
}

// generateKey generates a rate limiting key based on the configured strategy
func (rlm *RateLimitMiddleware) generateKey(ctx context.Context, clientIP string, userID string) string {
	switch rlm.config.KeyStrategy {
	case "user":
		if userID != "" {
			return "user:" + userID
		}
		return "ip:" + clientIP
	case "ip_user":
		if userID != "" {
			return "ip_user:" + clientIP + ":" + userID
		}
		return "ip:" + clientIP
	default: // "ip"
		return "ip:" + clientIP
	}
}

// shouldSkipPath checks if the path should skip rate limiting
func (rlm *RateLimitMiddleware) shouldSkipPath(path string) bool {
	for _, skipPath := range rlm.config.SkipPaths {
		if path == skipPath {
			return true
		}
	}
	return false
}

// getLimitsForUser returns appropriate limits based on user authentication status
func (rlm *RateLimitMiddleware) getLimitsForUser(userID string) (int, time.Duration) {
	if userID != "" {
		// Authenticated user
		return rlm.config.AuthenticatedLimit, rlm.config.AuthenticatedWindow
	}
	// Anonymous user
	return rlm.config.AnonymousLimit, rlm.config.AnonymousWindow
}

// setRateLimitHeaders sets rate limiting headers in the response
func (rlm *RateLimitMiddleware) setRateLimitHeaders(ctx context.Context, c *gin.Context, stats *RateLimitStats, limit int, window time.Duration) {
	if !rlm.config.IncludeHeaders {
		return
	}
	
	c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", limit))
	c.Header("X-RateLimit-Remaining", fmt.Sprintf("%d", stats.Remaining))
	c.Header("X-RateLimit-Reset", fmt.Sprintf("%d", stats.ResetTime.Unix()))
	c.Header("X-RateLimit-Window", window.String())
	
	// Add request ID if available
	if requestID := stats.RequestID; requestID != "" {
		c.Header("X-Request-ID", requestID)
	}
}

// Middleware returns a Gin middleware function for rate limiting
func (rlm *RateLimitMiddleware) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		
		// Check if path should skip rate limiting
		if rlm.shouldSkipPath(c.Request.URL.Path) {
			rlm.logger.Info(ctx, "Skipping rate limit for path", map[string]interface{}{
				"path": c.Request.URL.Path,
			})
			c.Next()
			return
		}
		
		// Extract client information
		clientIP := c.ClientIP()
		userID := ""
		
		// Try to get user ID from context or Gin context
		if uid, exists := c.Get("user_id"); exists {
			if uidStr, ok := uid.(string); ok {
				userID = uidStr
			}
		}
		if userID == "" {
			if ctx != nil {
				if uid, ok := ctx.Value("user_id").(string); ok {
					userID = uid
				}
			}
		}
		
		// Generate rate limiting key
		key := rlm.generateKey(ctx, clientIP, userID)
		
		// Get appropriate limits
		limit, window := rlm.getLimitsForUser(userID)
		
		// Check rate limit
		allowed, err := rlm.limiter.Allow(ctx, key, limit, window)
		if err != nil {
			rlm.logger.Error(ctx, "Rate limit check failed", err, map[string]interface{}{
				"key":      key,
				"client_ip": clientIP,
				"user_id":   userID,
			})
			
			// On error, allow the request but log the issue
			c.Next()
			return
		}
		
		// Get stats for headers
		stats, err := rlm.limiter.GetStats(ctx, key)
		if err != nil {
			rlm.logger.Warn(ctx, "Failed to get rate limit stats", map[string]interface{}{
				"key": key,
				"error": err.Error(),
			})
			// Create basic stats
			stats = &RateLimitStats{
				Key:       key,
				Count:     0,
				Limit:     limit,
				Remaining: limit,
				ResetTime: time.Now().Add(window),
				Window:    window,
			}
		}
		
		// Set rate limit headers
		rlm.setRateLimitHeaders(ctx, c, stats, limit, window)
		
		if !allowed {
			// Rate limit exceeded
			rlm.logger.Warn(ctx, "Rate limit exceeded", map[string]interface{}{
				"key":       key,
				"client_ip": clientIP,
				"user_id":   userID,
				"limit":     limit,
				"window":    window,
				"count":     stats.Count,
			})
			
			// Create rate limit error
			rateLimitErr := rlm.wrapper.RateLimitError(ctx, limit, window)
			
			// Set additional headers for rate limited requests
			c.Header("Retry-After", fmt.Sprintf("%.0f", window.Seconds()))
			
			// Create error response
			errorResponse := ErrorResponse{
				Success: false,
				Error: &ErrorDetail{
					Code:      rateLimitErr.Code,
					Message:   rateLimitErr.Message,
					RequestID: rateLimitErr.RequestID,
					Timestamp: rateLimitErr.Timestamp,
					Details: map[string]string{
						"limit":     fmt.Sprintf("%d", limit),
						"window":    window.String(),
						"reset_at":  stats.ResetTime.Format(time.RFC3339),
					},
				},
				Meta: &ResponseMetadata{
					RequestID:      rateLimitErr.RequestID,
					ProcessingTime: rlm.calculateProcessingTime(ctx),
					Version:        "1.0.0",
					Timestamp:      time.Now(),
				},
			}
			
			c.JSON(429, errorResponse)
			c.Abort()
			return
		}
		
		// Log successful rate limit check
		rlm.logger.Info(ctx, "Rate limit check passed", map[string]interface{}{
			"key":       key,
			"client_ip": clientIP,
			"user_id":   userID,
			"count":     stats.Count,
			"limit":     limit,
			"remaining": stats.Remaining,
		})
		
		c.Next()
	}
}

// calculateProcessingTime calculates request processing time from context
func (rlm *RateLimitMiddleware) calculateProcessingTime(ctx context.Context) time.Duration {
	if ctx != nil {
		if startTime, ok := ctx.Value("start_time").(time.Time); ok {
			return time.Since(startTime)
		}
	}
	return 0
}

// CleanupMiddleware returns a middleware that periodically cleans up expired rate limit entries
func (rlm *RateLimitMiddleware) CleanupMiddleware(interval time.Duration) gin.HandlerFunc {
	// Start cleanup goroutine
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		
		for {
			select {
			case <-ticker.C:
				ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
				err := rlm.limiter.Cleanup(ctx)
				if err != nil {
					rlm.logger.Error(ctx, "Rate limit cleanup failed", err, map[string]interface{}{
						"interval": interval,
					})
				} else {
					rlm.logger.Info(ctx, "Rate limit cleanup completed", map[string]interface{}{
						"interval": interval,
					})
				}
				cancel()
			}
		}
	}()
	
	// Return a no-op middleware since cleanup runs in background
	return func(c *gin.Context) {
		c.Next()
	}
}

// RateLimitMetrics holds metrics for rate limiting
type RateLimitMetrics struct {
	TotalRequests     int64             `json:"total_requests"`
	AllowedRequests   int64             `json:"allowed_requests"`
	BlockedRequests   int64             `json:"blocked_requests"`
	ErrorCount        int64             `json:"error_count"`
	ActiveKeys        int               `json:"active_keys"`
	LastCleanup       time.Time         `json:"last_cleanup"`
	ViolationsByKey   map[string]int64  `json:"violations_by_key"`
	RequestsByKey     map[string]int64  `json:"requests_by_key"`
	mu                sync.RWMutex      `json:"-"`
}

// NewRateLimitMetrics creates a new metrics collector
func NewRateLimitMetrics() *RateLimitMetrics {
	return &RateLimitMetrics{
		ViolationsByKey: make(map[string]int64),
		RequestsByKey:   make(map[string]int64),
	}
}

// IncrementTotal increments total request count
func (m *RateLimitMetrics) IncrementTotal() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.TotalRequests++
}

// IncrementAllowed increments allowed request count
func (m *RateLimitMetrics) IncrementAllowed() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.AllowedRequests++
}

// IncrementBlocked increments blocked request count and violations for key
func (m *RateLimitMetrics) IncrementBlocked(key string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.BlockedRequests++
	m.ViolationsByKey[key]++
}

// IncrementError increments error count
func (m *RateLimitMetrics) IncrementError() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.ErrorCount++
}

// IncrementKeyRequests increments request count for a specific key
func (m *RateLimitMetrics) IncrementKeyRequests(key string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.RequestsByKey[key]++
}

// SetActiveKeys sets the number of active keys
func (m *RateLimitMetrics) SetActiveKeys(count int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.ActiveKeys = count
}

// SetLastCleanup sets the last cleanup time
func (m *RateLimitMetrics) SetLastCleanup(t time.Time) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.LastCleanup = t
}

// GetSnapshot returns a snapshot of current metrics
func (m *RateLimitMetrics) GetSnapshot() RateLimitMetrics {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	// Create deep copy of maps
	violationsCopy := make(map[string]int64)
	requestsCopy := make(map[string]int64)
	
	for k, v := range m.ViolationsByKey {
		violationsCopy[k] = v
	}
	
	for k, v := range m.RequestsByKey {
		requestsCopy[k] = v
	}
	
	return RateLimitMetrics{
		TotalRequests:   m.TotalRequests,
		AllowedRequests: m.AllowedRequests,
		BlockedRequests: m.BlockedRequests,
		ErrorCount:      m.ErrorCount,
		ActiveKeys:      m.ActiveKeys,
		LastCleanup:     m.LastCleanup,
		ViolationsByKey: violationsCopy,
		RequestsByKey:   requestsCopy,
	}
}

// Reset resets all metrics
func (m *RateLimitMetrics) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	m.TotalRequests = 0
	m.AllowedRequests = 0
	m.BlockedRequests = 0
	m.ErrorCount = 0
	m.ActiveKeys = 0
	m.LastCleanup = time.Time{}
	m.ViolationsByKey = make(map[string]int64)
	m.RequestsByKey = make(map[string]int64)
}

// RateLimitMonitor provides monitoring and logging for rate limiting
type RateLimitMonitor struct {
	metrics *RateLimitMetrics
	logger  Logger
	limiter RateLimiter
}

// NewRateLimitMonitor creates a new rate limit monitor
func NewRateLimitMonitor(limiter RateLimiter, logger Logger) *RateLimitMonitor {
	if logger == nil {
		logger = &defaultLogger{}
	}
	
	return &RateLimitMonitor{
		metrics: NewRateLimitMetrics(),
		logger:  logger,
		limiter: limiter,
	}
}

// LogViolation logs a rate limit violation with context
func (m *RateLimitMonitor) LogViolation(ctx context.Context, key string, clientIP string, userID string, limit int, window time.Duration, count int) {
	m.metrics.IncrementBlocked(key)
	
	m.logger.Warn(ctx, "Rate limit violation", map[string]interface{}{
		"key":       key,
		"client_ip": clientIP,
		"user_id":   userID,
		"limit":     limit,
		"window":    window.String(),
		"count":     count,
		"violation_type": "rate_limit_exceeded",
	})
}

// LogAllowed logs an allowed request with context
func (m *RateLimitMonitor) LogAllowed(ctx context.Context, key string, clientIP string, userID string, count int, limit int, remaining int) {
	m.metrics.IncrementAllowed()
	m.metrics.IncrementKeyRequests(key)
	
	// Only log detailed info for debugging or if close to limit
	if remaining <= limit/10 { // Log when 90% of limit is used
		m.logger.Info(ctx, "Rate limit check - approaching limit", map[string]interface{}{
			"key":       key,
			"client_ip": clientIP,
			"user_id":   userID,
			"count":     count,
			"limit":     limit,
			"remaining": remaining,
			"usage_percent": float64(count) / float64(limit) * 100,
		})
	}
}

// LogError logs a rate limiting error with context
func (m *RateLimitMonitor) LogError(ctx context.Context, key string, clientIP string, userID string, err error, operation string) {
	m.metrics.IncrementError()
	
	m.logger.Error(ctx, "Rate limiting error", err, map[string]interface{}{
		"key":       key,
		"client_ip": clientIP,
		"user_id":   userID,
		"operation": operation,
	})
}

// LogCleanup logs cleanup operations with context
func (m *RateLimitMonitor) LogCleanup(ctx context.Context, expiredCount int, totalKeys int, duration time.Duration) {
	m.metrics.SetLastCleanup(time.Now())
	m.metrics.SetActiveKeys(totalKeys - expiredCount)
	
	m.logger.Info(ctx, "Rate limit cleanup completed", map[string]interface{}{
		"expired_entries": expiredCount,
		"active_keys":     totalKeys - expiredCount,
		"cleanup_duration": duration.String(),
	})
}

// GetMetrics returns current metrics
func (m *RateLimitMonitor) GetMetrics(ctx context.Context) RateLimitMetrics {
	m.metrics.IncrementTotal() // Count this metrics request
	
	snapshot := m.metrics.GetSnapshot()
	
	m.logger.Info(ctx, "Rate limit metrics requested", map[string]interface{}{
		"total_requests":   snapshot.TotalRequests,
		"allowed_requests": snapshot.AllowedRequests,
		"blocked_requests": snapshot.BlockedRequests,
		"error_count":      snapshot.ErrorCount,
		"active_keys":      snapshot.ActiveKeys,
	})
	
	return snapshot
}

// ResetMetrics resets all metrics
func (m *RateLimitMonitor) ResetMetrics(ctx context.Context) {
	m.metrics.Reset()
	
	m.logger.Info(ctx, "Rate limit metrics reset", map[string]interface{}{
		"reset_time": time.Now(),
	})
}

// StartPeriodicLogging starts periodic logging of rate limit statistics
func (m *RateLimitMonitor) StartPeriodicLogging(ctx context.Context, interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		
		for {
			select {
			case <-ctx.Done():
				m.logger.Info(ctx, "Stopping periodic rate limit logging", map[string]interface{}{
					"reason": "context cancelled",
				})
				return
			case <-ticker.C:
				snapshot := m.metrics.GetSnapshot()
				
				// Calculate rates
				var blockRate float64
				if snapshot.TotalRequests > 0 {
					blockRate = float64(snapshot.BlockedRequests) / float64(snapshot.TotalRequests) * 100
				}
				
				m.logger.Info(ctx, "Rate limit periodic stats", map[string]interface{}{
					"total_requests":   snapshot.TotalRequests,
					"allowed_requests": snapshot.AllowedRequests,
					"blocked_requests": snapshot.BlockedRequests,
					"block_rate_percent": blockRate,
					"error_count":      snapshot.ErrorCount,
					"active_keys":      snapshot.ActiveKeys,
					"top_violators":    m.getTopViolators(snapshot.ViolationsByKey, 5),
				})
			}
		}
	}()
}

// getTopViolators returns the top N keys with most violations
func (m *RateLimitMonitor) getTopViolators(violations map[string]int64, n int) []map[string]interface{} {
	type keyViolation struct {
		key        string
		violations int64
	}
	
	var sorted []keyViolation
	for key, count := range violations {
		sorted = append(sorted, keyViolation{key: key, violations: count})
	}
	
	// Simple bubble sort for small n
	for i := 0; i < len(sorted)-1; i++ {
		for j := 0; j < len(sorted)-i-1; j++ {
			if sorted[j].violations < sorted[j+1].violations {
				sorted[j], sorted[j+1] = sorted[j+1], sorted[j]
			}
		}
	}
	
	// Take top n
	if len(sorted) > n {
		sorted = sorted[:n]
	}
	
	result := make([]map[string]interface{}, len(sorted))
	for i, kv := range sorted {
		result[i] = map[string]interface{}{
			"key":        kv.key,
			"violations": kv.violations,
		}
	}
	
	return result
}

// Enhanced RateLimitMiddleware with monitoring
type MonitoredRateLimitMiddleware struct {
	*RateLimitMiddleware
	monitor *RateLimitMonitor
}

// NewMonitoredRateLimitMiddleware creates a new monitored rate limiting middleware
func NewMonitoredRateLimitMiddleware(limiter RateLimiter, config *RateLimitConfig, logger Logger) *MonitoredRateLimitMiddleware {
	baseMiddleware := NewRateLimitMiddleware(limiter, config, logger)
	monitor := NewRateLimitMonitor(limiter, logger)
	
	return &MonitoredRateLimitMiddleware{
		RateLimitMiddleware: baseMiddleware,
		monitor:            monitor,
	}
}

// GetMonitor returns the rate limit monitor
func (mrlm *MonitoredRateLimitMiddleware) GetMonitor() *RateLimitMonitor {
	return mrlm.monitor
}

// MonitoredMiddleware returns a Gin middleware function with monitoring
func (mrlm *MonitoredRateLimitMiddleware) MonitoredMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		
		// Increment total requests
		mrlm.monitor.metrics.IncrementTotal()
		
		// Check if path should skip rate limiting
		if mrlm.shouldSkipPath(c.Request.URL.Path) {
			mrlm.logger.Info(ctx, "Skipping rate limit for path", map[string]interface{}{
				"path": c.Request.URL.Path,
			})
			c.Next()
			return
		}
		
		// Extract client information
		clientIP := c.ClientIP()
		userID := ""
		
		// Try to get user ID from context or Gin context
		if uid, exists := c.Get("user_id"); exists {
			if uidStr, ok := uid.(string); ok {
				userID = uidStr
			}
		}
		if userID == "" {
			if ctx != nil {
				if uid, ok := ctx.Value("user_id").(string); ok {
					userID = uid
				}
			}
		}
		
		// Generate rate limiting key
		key := mrlm.generateKey(ctx, clientIP, userID)
		
		// Get appropriate limits
		limit, window := mrlm.getLimitsForUser(userID)
		
		// Check rate limit
		allowed, err := mrlm.limiter.Allow(ctx, key, limit, window)
		if err != nil {
			mrlm.monitor.LogError(ctx, key, clientIP, userID, err, "rate_limit_check")
			
			// On error, allow the request but log the issue
			c.Next()
			return
		}
		
		// Get stats for headers and monitoring
		stats, err := mrlm.limiter.GetStats(ctx, key)
		if err != nil {
			mrlm.logger.Warn(ctx, "Failed to get rate limit stats", map[string]interface{}{
				"key": key,
				"error": err.Error(),
			})
			// Create basic stats
			stats = &RateLimitStats{
				Key:       key,
				Count:     0,
				Limit:     limit,
				Remaining: limit,
				ResetTime: time.Now().Add(window),
				Window:    window,
			}
		}
		
		// Set rate limit headers
		mrlm.setRateLimitHeaders(ctx, c, stats, limit, window)
		
		if !allowed {
			// Rate limit exceeded - log violation
			mrlm.monitor.LogViolation(ctx, key, clientIP, userID, limit, window, stats.Count)
			
			// Create rate limit error
			rateLimitErr := mrlm.wrapper.RateLimitError(ctx, limit, window)
			
			// Set additional headers for rate limited requests
			c.Header("Retry-After", fmt.Sprintf("%.0f", window.Seconds()))
			
			// Create error response
			errorResponse := ErrorResponse{
				Success: false,
				Error: &ErrorDetail{
					Code:      rateLimitErr.Code,
					Message:   rateLimitErr.Message,
					RequestID: rateLimitErr.RequestID,
					Timestamp: rateLimitErr.Timestamp,
					Details: map[string]string{
						"limit":     fmt.Sprintf("%d", limit),
						"window":    window.String(),
						"reset_at":  stats.ResetTime.Format(time.RFC3339),
					},
				},
				Meta: &ResponseMetadata{
					RequestID:      rateLimitErr.RequestID,
					ProcessingTime: mrlm.calculateProcessingTime(ctx),
					Version:        "1.0.0",
					Timestamp:      time.Now(),
				},
			}
			
			c.JSON(429, errorResponse)
			c.Abort()
			return
		}
		
		// Log allowed request
		mrlm.monitor.LogAllowed(ctx, key, clientIP, userID, stats.Count, limit, stats.Remaining)
		
		c.Next()
	}
}