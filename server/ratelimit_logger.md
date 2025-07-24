# Rate Limiter dan Logger Implementation

## Overview

Dokumen ini menjelaskan implementasi Rate Limiter dan Logger di server Develapar Blog, termasuk konfigurasi, cara kerja, dan optimasi untuk performa.

## Table of Contents

1. [Rate Limiter](#rate-limiter)
2. [Logger System](#logger-system)
3. [Integrasi dan Konfigurasi](#integrasi-dan-konfigurasi)
4. [Optimasi untuk Stress Testing](#optimasi-untuk-stress-testing)
5. [Monitoring dan Metrics](#monitoring-dan-metrics)
6. [Troubleshooting](#troubleshooting)

---

## Rate Limiter

### Arsitektur Rate Limiter

Server menggunakan **Sliding Window Rate Limiter** dengan komponen-komponen berikut:

```
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│   Middleware    │───▶│   Rate Limiter   │───▶│  Memory Store   │
│                 │    │                  │    │                 │
│ - Key Generation│    │ - Allow/Deny     │    │ - Entry Storage │
│ - Header Setting│    │ - Stats Tracking │    │ - Cleanup       │
│ - Error Response│    │ - Monitoring     │    │ - Expiry Check  │
└─────────────────┘    └──────────────────┘    └─────────────────┘
```

### Komponen Utama

#### 1. RateLimitStore Interface

```go
type RateLimitStore interface {
    Increment(ctx context.Context, key string, window time.Duration) (int, error)
    Get(ctx context.Context, key string) (int, error)
    Reset(ctx context.Context, key string) error
    GetExpiry(ctx context.Context, key string) (time.Time, error)
    SetExpiry(ctx context.Context, key string, expiry time.Time) error
    Cleanup(ctx context.Context) error
}
```

**Implementasi:** `inMemoryStore`

- Menggunakan `map[string]*RateLimitEntry` untuk storage
- Thread-safe dengan `sync.RWMutex`
- Automatic cleanup untuk expired entries

#### 2. RateLimiter Interface

```go
type RateLimiter interface {
    Allow(ctx context.Context, key string, limit int, window time.Duration) (bool, error)
    GetRemaining(ctx context.Context, key string) (int, error)
    Reset(ctx context.Context, key string) error
    GetStats(ctx context.Context, key string) (*RateLimitStats, error)
    Cleanup(ctx context.Context) error
}
```

**Implementasi:** `slidingWindowRateLimiter`

- Sliding window algorithm untuk rate limiting
- Context-aware operations dengan timeout support
- Comprehensive error handling

#### 3. RateLimitMiddleware

```go
type MonitoredRateLimitMiddleware struct {
    *RateLimitMiddleware
    monitor *RateLimitMonitor
}
```

**Fitur:**

- Automatic key generation berdasarkan IP/User ID
- Differentiated limits untuk authenticated vs anonymous users
- HTTP headers untuk rate limit information
- Skip paths untuk health checks dan metrics

### Konfigurasi Rate Limiter

#### Environment Variables

```env
# Rate Limiting Configuration
RATE_LIMIT_ENABLED=true                    # Enable/disable rate limiting
RATE_LIMIT_REQUESTS_PER_MINUTE=100         # Default requests per minute
RATE_LIMIT_BURST_SIZE=5                    # Burst allowance
RATE_LIMIT_CLEANUP_INTERVAL=10s            # Cleanup frequency
RATE_LIMIT_WINDOW_SIZE=10s                 # Sliding window size
RATE_LIMIT_AUTHENTICATED_RPM=200           # Limit for authenticated users
RATE_LIMIT_ANONYMOUS_RPM=50                # Limit for anonymous users
RATE_LIMIT_REQUEST_TIMEOUT=1s              # Operation timeout
```

#### Konfigurasi Berdasarkan Environment

**Development:**

```env
RATE_LIMIT_AUTHENTICATED_RPM=1000
RATE_LIMIT_ANONYMOUS_RPM=500
RATE_LIMIT_WINDOW_SIZE=1m
```

**Production:**

```env
RATE_LIMIT_AUTHENTICATED_RPM=200
RATE_LIMIT_ANONYMOUS_RPM=50
RATE_LIMIT_WINDOW_SIZE=10s
```

**Stress Testing:**

```env
RATE_LIMIT_AUTHENTICATED_RPM=100
RATE_LIMIT_ANONYMOUS_RPM=25
RATE_LIMIT_WINDOW_SIZE=5s
```

### Key Generation Strategy

Rate limiter menggunakan strategi key generation yang fleksibel:

```go
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
```

**Strategi yang Tersedia:**

- `ip`: Rate limit berdasarkan IP address (default)
- `user`: Rate limit berdasarkan User ID, fallback ke IP
- `ip_user`: Kombinasi IP dan User ID

### Response Headers

Rate limiter menambahkan headers informatif:

```http
X-RateLimit-Limit: 50
X-RateLimit-Remaining: 42
X-RateLimit-Reset: 1753363252
X-RateLimit-Window: 10s
X-Request-ID: bd217edc-b466-4508-b892-b88381793adf
```

### Error Response

Ketika rate limit terlampaui:

```json
{
  "status": 429,
  "message": "Rate limit exceeded",
  "retry_after": 10,
  "request_id": "bd217edc-b466-4508-b892-b88381793adf",
  "timestamp": "2025-07-24T13:20:44Z"
}
```

---

## Logger System

### Arsitektur Logger

```
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│ Logger Factory  │───▶│   JSON Logger    │───▶│   Output        │
│                 │    │                  │    │                 │
│ - Level Config  │    │ - Context Aware  │    │ - stdout/stderr │
│ - Component Mgmt│    │ - Structured Log │    │ - File Output   │
│ - Writer Config │    │ - Field Support  │    │ - Log Rotation  │
└─────────────────┘    └──────────────────┘    └─────────────────┘
```

### Log Level Hierarchy

```go
const (
    DebugLevel LogLevel = iota  // 0 - Detailed debugging info
    InfoLevel                   // 1 - General information
    WarnLevel                   // 2 - Warning conditions
    ErrorLevel                  // 3 - Error conditions
    FatalLevel                  // 4 - Fatal errors (exits app)
)
```

**Filtering Logic:**

- Jika `LOG_LEVEL=ERROR`, hanya ERROR dan FATAL yang akan di-log
- Jika `LOG_LEVEL=INFO`, INFO, WARN, ERROR, dan FATAL yang akan di-log

### Logger Interface

```go
type Logger interface {
    Info(ctx context.Context, msg string, fields ...Field)
    Warn(ctx context.Context, msg string, fields ...Field)
    Error(ctx context.Context, msg string, err error, fields ...Field)
    Debug(ctx context.Context, msg string, fields ...Field)
    Fatal(ctx context.Context, msg string, err error, fields ...Field)
    WithContext(ctx context.Context) Logger
    SetLevel(level LogLevel)
    GetLevel() LogLevel
}
```

### Structured Logging

Logger menggunakan structured JSON format:

```json
{
  "level": "INFO",
  "message": "Request completed successfully",
  "timestamp": "2025-07-24T13:20:44.915720631Z",
  "request_id": "bd217edc-b466-4508-b892-b88381793adf",
  "user_id": "user-123",
  "fields": {
    "component": "request",
    "method": "GET",
    "path": "/api/v1/health",
    "status_code": 200,
    "processing_time": "1.912314ms",
    "response_size": 899
  },
  "caller": "middleware.go:245"
}
```

### Context-Aware Logging

Logger otomatis mengekstrak informasi dari context:

```go
// Extract context information if available
if ctx != nil {
    if requestID := GetRequestIDFromContext(ctx); requestID != "" {
        entry.RequestID = requestID
    }
    if userID := GetUserIDFromContext(ctx); userID != "" {
        entry.UserID = userID
    }
}
```

### Logger Factory

LoggerFactory mengelola multiple logger instances dengan konfigurasi terpusat:

```go
type LoggerFactory struct {
    level   LogLevel
    writers map[string]io.Writer
}

// Usage
loggerFactory := utils.NewLoggerFactory(logLevel)
metricsLogger := loggerFactory.GetLogger("metrics")
rateLimitLogger := loggerFactory.GetLogger("rate_limiter")
requestLogger := loggerFactory.GetLogger("request")
```

### Request Logging Middleware

Middleware untuk logging request/response dengan fitur:

- **Request Information**: Method, path, headers, body (limited)
- **Response Information**: Status code, headers, body (limited), processing time
- **Performance Metrics**: Response time, request size, response size
- **Security**: Automatic filtering of sensitive headers
- **Context Correlation**: Request ID dan User ID tracking

**Optimasi untuk Performance:**

- Body logging dibatasi 1KB
- Sensitive headers di-filter otomatis
- Optional body capture untuk stress testing

---

## Integrasi dan Konfigurasi

### Middleware Order

Urutan middleware yang optimal untuk context propagation:

```go
1. Recovery middleware (panic handling)
2. CORS middleware
3. Context middleware (request_id, user_id injection)
4. Request logging middleware
5. Rate limiting middleware
6. Metrics collection middleware
7. Error handling middleware
```

### Konfigurasi Environment

#### Development Environment

```env
LOG_LEVEL=DEBUG
RATE_LIMIT_ENABLED=false
RATE_LIMIT_AUTHENTICATED_RPM=1000
RATE_LIMIT_ANONYMOUS_RPM=500
```

#### Production Environment

```env
LOG_LEVEL=WARN
RATE_LIMIT_ENABLED=true
RATE_LIMIT_AUTHENTICATED_RPM=200
RATE_LIMIT_ANONYMOUS_RPM=50
```

#### Stress Testing Environment

```env
LOG_LEVEL=ERROR
RATE_LIMIT_ENABLED=true
RATE_LIMIT_AUTHENTICATED_RPM=100
RATE_LIMIT_ANONYMOUS_RPM=25
```

### Logger Adapter

Untuk kompatibilitas dengan middleware yang membutuhkan interface berbeda:

```go
type loggerAdapter struct {
    logger utils.Logger
}

func (la *loggerAdapter) Error(ctx context.Context, msg string, err error, fields map[string]interface{}) {
    var utilsFields []utils.Field
    for key, value := range fields {
        utilsFields = append(utilsFields, utils.Field{Key: key, Value: value})
    }
    la.logger.Error(ctx, msg, err, utilsFields...)
}
```

---

## Optimasi untuk Stress Testing

### Memory Optimization

#### 1. Rate Limiter Memory Management

```go
// Cleanup lebih sering untuk mencegah memory leak
RATE_LIMIT_CLEANUP_INTERVAL=10s  // Dari 5m ke 10s

// Window size lebih kecil untuk mengurangi entries
RATE_LIMIT_WINDOW_SIZE=10s       // Dari 1m ke 10s
```

#### 2. Logging Optimization

```go
// Disable request/response body logging
// Skip body reading to reduce memory usage
// if c.Request.Body != nil {
//     bodyBytes, err := io.ReadAll(c.Request.Body)
//     ...
// }
```

#### 3. Reduced Logging Frequency

```go
// Log setiap 100th reset untuk mengurangi I/O
if entry.Count%100 == 0 {
    s.logger.Info(ctx, "Rate limit window reset", ...)
}
```

### Database Connection Optimization

```env
# Reduced untuk VPS dengan resource terbatas
DB_MAX_OPEN_CONNS=3      # Dari 25 ke 3
DB_MAX_IDLE_CONNS=1      # Dari 10 ke 1
DB_CONN_MAX_LIFETIME=2m  # Dari 30m ke 2m
```

### Context Timeout Optimization

```env
# Faster failure untuk resource terbatas
CONTEXT_REQUEST_TIMEOUT=10s      # Dari 30s ke 10s
CONTEXT_DATABASE_TIMEOUT=5s      # Dari 15s ke 5s
CONTEXT_VALIDATION_TIMEOUT=2s    # Dari 5s ke 2s
CONTEXT_LOGGING_TIMEOUT=1s       # Dari 2s ke 1s
```

---

## Monitoring dan Metrics

### Rate Limit Metrics

```go
type RateLimitMetrics struct {
    TotalRequests     int64             // Total requests processed
    AllowedRequests   int64             // Requests allowed
    BlockedRequests   int64             // Requests blocked
    ErrorCount        int64             // Errors during rate limiting
    ActiveKeys        int               // Active rate limit keys
    LastCleanup       time.Time         // Last cleanup time
    ViolationsByKey   map[string]int64  // Violations per key
    RequestsByKey     map[string]int64  // Requests per key
}
```

### Periodic Monitoring

Rate limiter monitor melakukan periodic logging:

```go
// Start periodic logging setiap 10 menit
rateLimitMonitor.StartPeriodicLogging(ctx, 10*time.Minute)
```

**Output Example:**

```json
{
  "level": "INFO",
  "message": "Rate limit periodic stats",
  "fields": {
    "total_requests": 1500,
    "allowed_requests": 1450,
    "blocked_requests": 50,
    "block_rate_percent": 3.33,
    "error_count": 0,
    "active_keys": 25,
    "top_violators": [
      { "key": "ip:192.168.1.100", "violations": 15 },
      { "key": "ip:10.0.0.50", "violations": 8 }
    ]
  }
}
```

### Request Metrics

```go
type RequestMetrics struct {
    TotalRequests       int64
    RequestsByMethod    map[string]int64
    RequestsByStatus    map[int]int64
    AverageResponseTime time.Duration
    SlowRequests        int64  // Requests > 1 second
}
```

---

## Troubleshooting

### Common Issues

#### 1. Log Level Tidak Berpengaruh

**Gejala:** Masih melihat INFO logs padahal `LOG_LEVEL=ERROR`

**Penyebab:** Logger menggunakan hardcoded InfoLevel

**Solusi:**

```go
// Sebelum (salah)
logger := utils.NewDefaultLogger("component")

// Sesudah (benar)
logLevel := parseLogLevel(config.LoggingConfig.Level)
loggerFactory := utils.NewLoggerFactory(logLevel)
logger := loggerFactory.GetLogger("component")
```

#### 2. Rate Limiter Memory Leak

**Gejala:** Memory usage terus naik saat stress testing

**Penyebab:** Cleanup interval terlalu lama, entries menumpuk

**Solusi:**

```env
# Cleanup lebih sering
RATE_LIMIT_CLEANUP_INTERVAL=10s

# Window size lebih kecil
RATE_LIMIT_WINDOW_SIZE=10s
```

#### 3. Database Connection Pool Exhausted

**Gejala:** "too many connections" error saat stress test

**Penyebab:** Terlalu banyak concurrent connections untuk VPS kecil

**Solusi:**

```env
# Kurangi connection pool
DB_MAX_OPEN_CONNS=3
DB_MAX_IDLE_CONNS=1
DB_QUERY_TIMEOUT=5s
```

#### 4. Context Timeout Errors

**Gejala:** Banyak context deadline exceeded errors

**Penyebab:** Timeout terlalu pendek untuk operasi yang kompleks

**Solusi:**

```env
# Sesuaikan timeout berdasarkan kebutuhan
CONTEXT_REQUEST_TIMEOUT=15s
CONTEXT_DATABASE_TIMEOUT=10s
```

### Debugging Tips

#### 1. Enable Debug Logging Sementara

```env
LOG_LEVEL=DEBUG  # Untuk debugging, kembalikan ke ERROR setelah selesai
```

#### 2. Monitor Resource Usage

```bash
# Monitor memory dan CPU
./monitor_resources.sh

# Monitor database connections
ss -tuln | grep :5432 | wc -l
```

#### 3. Rate Limit Testing

```bash
# Test rate limit dengan curl
for i in {1..60}; do
  curl -w "%{http_code}\n" -o /dev/null -s http://localhost:4300/api/v1/health
  sleep 0.1
done
```

#### 4. Log Analysis

```bash
# Filter hanya ERROR logs
./test-server 2>&1 | grep '"level":"ERROR"'

# Count log levels
./test-server 2>&1 | jq -r '.level' | sort | uniq -c
```

---

## Best Practices

### 1. Rate Limiting

- Gunakan strategi key yang sesuai dengan use case
- Set cleanup interval yang reasonable (tidak terlalu sering, tidak terlalu jarang)
- Monitor metrics secara berkala
- Adjust limits berdasarkan capacity server

### 2. Logging

- Gunakan log level yang sesuai dengan environment
- Hindari logging sensitive information
- Limit body size untuk request/response logging
- Gunakan structured logging untuk easier parsing

### 3. Performance

- Disable unnecessary logging saat stress testing
- Monitor memory usage dan cleanup secara berkala
- Adjust timeout values berdasarkan server capacity
- Use connection pooling yang optimal

### 4. Monitoring

- Set up periodic metrics collection
- Monitor top violators untuk rate limiting
- Track response times dan error rates
- Alert pada threshold tertentu

---

## Kesimpulan

Implementasi Rate Limiter dan Logger di server Develapar Blog dirancang untuk:

1. **Scalability**: Mendukung high-traffic dengan efficient memory usage
2. **Observability**: Comprehensive logging dan metrics untuk monitoring
3. **Flexibility**: Configurable limits dan log levels untuk berbagai environment
4. **Performance**: Optimized untuk resource-constrained environments
5. **Reliability**: Graceful degradation dan error handling

Dengan konfigurasi yang tepat, sistem ini dapat menangani stress testing dan production load dengan stabil sambil memberikan visibility yang baik untuk monitoring dan debugging.
