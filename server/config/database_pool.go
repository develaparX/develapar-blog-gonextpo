package config

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

// ConnectionStats represents database connection pool statistics
type ConnectionStats struct {
	OpenConnections     int           `json:"open_connections"`
	InUseConnections    int           `json:"in_use_connections"`
	IdleConnections     int           `json:"idle_connections"`
	WaitCount           int64         `json:"wait_count"`
	WaitDuration        time.Duration `json:"wait_duration"`
	MaxIdleClosed       int64         `json:"max_idle_closed"`
	MaxIdleTimeClosed   int64         `json:"max_idle_time_closed"`
	MaxLifetimeClosed   int64         `json:"max_lifetime_closed"`
}

// PoolConfig represents database connection pool configuration
type PoolConfig struct {
	MaxOpenConns    int           `json:"max_open_conns"`
	MaxIdleConns    int           `json:"max_idle_conns"`
	ConnMaxLifetime time.Duration `json:"conn_max_lifetime"`
	ConnMaxIdleTime time.Duration `json:"conn_max_idle_time"`
	ConnectTimeout  time.Duration `json:"connect_timeout"`
	QueryTimeout    time.Duration `json:"query_timeout"`
}

// ConnectionPoolManager interface for managing database connections with context
type ConnectionPoolManager interface {
	GetConnection(ctx context.Context) (*sql.DB, error)
	GetStats(ctx context.Context) ConnectionStats
	Configure(ctx context.Context, config PoolConfig) error
	HealthCheck(ctx context.Context) error
	Close(ctx context.Context) error
}

// connectionPoolManager implements ConnectionPoolManager
type connectionPoolManager struct {
	db     *sql.DB
	config PoolConfig
}

// NewConnectionPoolManager creates a new connection pool manager with context support
func NewConnectionPoolManager(ctx context.Context, dbConfig DbConfig, poolConfig PoolConfig) (ConnectionPoolManager, error) {
	urlConnection := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", 
		dbConfig.Host, dbConfig.Port, dbConfig.User, dbConfig.Password, dbConfig.Name)

	// Create context with timeout for connection establishment
	connectCtx, cancel := context.WithTimeout(ctx, poolConfig.ConnectTimeout)
	defer cancel()

	db, err := sql.Open(dbConfig.Driver, urlConnection)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	// Test the connection with context
	if err := db.PingContext(connectCtx); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	manager := &connectionPoolManager{
		db:     db,
		config: poolConfig,
	}

	// Configure the connection pool
	if err := manager.Configure(ctx, poolConfig); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to configure connection pool: %w", err)
	}

	return manager, nil
}

// GetConnection returns the database connection with context validation
func (cpm *connectionPoolManager) GetConnection(ctx context.Context) (*sql.DB, error) {
	// Check if context is already cancelled
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	// Perform a quick health check with context
	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := cpm.db.PingContext(pingCtx); err != nil {
		return nil, fmt.Errorf("database connection health check failed: %w", err)
	}

	return cpm.db, nil
}

// GetStats returns connection pool statistics with context
func (cpm *connectionPoolManager) GetStats(ctx context.Context) ConnectionStats {
	// Check if context is cancelled
	if ctx.Err() != nil {
		return ConnectionStats{}
	}

	stats := cpm.db.Stats()
	return ConnectionStats{
		OpenConnections:     stats.OpenConnections,
		InUseConnections:    stats.InUse,
		IdleConnections:     stats.Idle,
		WaitCount:           stats.WaitCount,
		WaitDuration:        stats.WaitDuration,
		MaxIdleClosed:       stats.MaxIdleClosed,
		MaxIdleTimeClosed:   stats.MaxIdleTimeClosed,
		MaxLifetimeClosed:   stats.MaxLifetimeClosed,
	}
}

// Configure sets up the connection pool parameters with context
func (cpm *connectionPoolManager) Configure(ctx context.Context, config PoolConfig) error {
	// Check if context is cancelled
	if ctx.Err() != nil {
		return ctx.Err()
	}

	// Set maximum number of open connections
	cpm.db.SetMaxOpenConns(config.MaxOpenConns)

	// Set maximum number of idle connections
	cpm.db.SetMaxIdleConns(config.MaxIdleConns)

	// Set maximum lifetime of connections
	cpm.db.SetConnMaxLifetime(config.ConnMaxLifetime)

	// Set maximum idle time of connections
	cpm.db.SetConnMaxIdleTime(config.ConnMaxIdleTime)

	// Update internal config
	cpm.config = config

	return nil
}

// HealthCheck performs a comprehensive health check with context
func (cpm *connectionPoolManager) HealthCheck(ctx context.Context) error {
	// Check if context is cancelled
	if ctx.Err() != nil {
		return ctx.Err()
	}

	// Create context with timeout for health check
	healthCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// Test basic connectivity
	if err := cpm.db.PingContext(healthCtx); err != nil {
		return fmt.Errorf("ping failed: %w", err)
	}

	// Test query execution
	var result int
	query := "SELECT 1"
	if err := cpm.db.QueryRowContext(healthCtx, query).Scan(&result); err != nil {
		return fmt.Errorf("test query failed: %w", err)
	}

	if result != 1 {
		return fmt.Errorf("test query returned unexpected result: %d", result)
	}

	// Check connection pool stats
	stats := cpm.GetStats(ctx)
	if stats.OpenConnections == 0 {
		return fmt.Errorf("no open connections available")
	}

	return nil
}

// Close closes the database connection pool with context
func (cpm *connectionPoolManager) Close(ctx context.Context) error {
	// Check if context is cancelled
	if ctx.Err() != nil {
		return ctx.Err()
	}

	if cpm.db != nil {
		return cpm.db.Close()
	}
	return nil
}

// DefaultPoolConfig returns a default pool configuration
func DefaultPoolConfig() PoolConfig {
	return PoolConfig{
		MaxOpenConns:    25,                // Maximum number of open connections
		MaxIdleConns:    10,                // Maximum number of idle connections
		ConnMaxLifetime: 30 * time.Minute,  // Maximum lifetime of connections
		ConnMaxIdleTime: 15 * time.Minute,  // Maximum idle time of connections
		ConnectTimeout:  10 * time.Second,  // Connection establishment timeout
		QueryTimeout:    30 * time.Second,  // Query execution timeout
	}
}