package utils

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"
)

// DBLogger interface for database query logging
type DBLogger interface {
	LogQuery(ctx context.Context, query string, args []interface{}, duration time.Duration, err error)
	LogSlowQuery(ctx context.Context, query string, args []interface{}, duration time.Duration, threshold time.Duration)
	LogError(ctx context.Context, query string, args []interface{}, err error)
}

// dbLogger implements DBLogger interface
type dbLogger struct {
	logger    Logger
	threshold time.Duration // Slow query threshold
}

// LogQuery logs a database query with context and execution time
func (dl *dbLogger) LogQuery(ctx context.Context, query string, args []interface{}, duration time.Duration, err error) {
	fields := []Field{
		StringField("query", sanitizeQuery(query)),
		DurationField("duration", duration),
		IntField("arg_count", len(args)),
	}

	// Add sanitized arguments (avoid logging sensitive data)
	if len(args) > 0 {
		sanitizedArgs := sanitizeArgs(args)
		fields = append(fields, Field{Key: "args", Value: sanitizedArgs})
	}

	if err != nil {
		fields = append(fields, ErrorField(err))
		dl.logger.Error(ctx, "Database query failed", err, fields...)
	} else {
		dl.logger.Debug(ctx, "Database query executed", fields...)
	}

	// Check for slow queries
	if duration > dl.threshold {
		dl.LogSlowQuery(ctx, query, args, duration, dl.threshold)
	}
}

// LogSlowQuery logs slow database queries
func (dl *dbLogger) LogSlowQuery(ctx context.Context, query string, args []interface{}, duration time.Duration, threshold time.Duration) {
	fields := []Field{
		StringField("query", sanitizeQuery(query)),
		DurationField("duration", duration),
		DurationField("threshold", threshold),
		IntField("arg_count", len(args)),
		BoolField("slow_query", true),
	}

	if len(args) > 0 {
		sanitizedArgs := sanitizeArgs(args)
		fields = append(fields, Field{Key: "args", Value: sanitizedArgs})
	}

	dl.logger.Warn(ctx, "Slow database query detected", fields...)
}

// LogError logs database errors with context
func (dl *dbLogger) LogError(ctx context.Context, query string, args []interface{}, err error) {
	fields := []Field{
		StringField("query", sanitizeQuery(query)),
		IntField("arg_count", len(args)),
		ErrorField(err),
	}

	if len(args) > 0 {
		sanitizedArgs := sanitizeArgs(args)
		fields = append(fields, Field{Key: "args", Value: sanitizedArgs})
	}

	dl.logger.Error(ctx, "Database error occurred", err, fields...)
}

// NewDBLogger creates a new database logger
func NewDBLogger(logger Logger, slowQueryThreshold time.Duration) DBLogger {
	if slowQueryThreshold == 0 {
		slowQueryThreshold = 1 * time.Second // Default threshold
	}
	return &dbLogger{
		logger:    logger,
		threshold: slowQueryThreshold,
	}
}

// LoggingDB wraps sql.DB to provide query logging with context
type LoggingDB struct {
	db       *sql.DB
	dbLogger DBLogger
}

// QueryContext executes a query with context and logging
func (ldb *LoggingDB) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	start := time.Now()
	rows, err := ldb.db.QueryContext(ctx, query, args...)
	duration := time.Since(start)
	
	ldb.dbLogger.LogQuery(ctx, query, args, duration, err)
	return rows, err
}

// QueryRowContext executes a query that returns a single row with context and logging
func (ldb *LoggingDB) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	start := time.Now()
	row := ldb.db.QueryRowContext(ctx, query, args...)
	duration := time.Since(start)
	
	// Note: QueryRow doesn't return an error directly, so we log without error
	ldb.dbLogger.LogQuery(ctx, query, args, duration, nil)
	return row
}

// ExecContext executes a query without returning rows with context and logging
func (ldb *LoggingDB) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	start := time.Now()
	result, err := ldb.db.ExecContext(ctx, query, args...)
	duration := time.Since(start)
	
	ldb.dbLogger.LogQuery(ctx, query, args, duration, err)
	return result, err
}

// PrepareContext prepares a statement with context and logging
func (ldb *LoggingDB) PrepareContext(ctx context.Context, query string) (*LoggingStmt, error) {
	start := time.Now()
	stmt, err := ldb.db.PrepareContext(ctx, query)
	duration := time.Since(start)
	
	ldb.dbLogger.LogQuery(ctx, fmt.Sprintf("PREPARE: %s", query), nil, duration, err)
	
	if err != nil {
		return nil, err
	}
	
	return &LoggingStmt{
		stmt:     stmt,
		query:    query,
		dbLogger: ldb.dbLogger,
	}, nil
}

// BeginTx starts a transaction with context and logging
func (ldb *LoggingDB) BeginTx(ctx context.Context, opts *sql.TxOptions) (*LoggingTx, error) {
	start := time.Now()
	tx, err := ldb.db.BeginTx(ctx, opts)
	duration := time.Since(start)
	
	ldb.dbLogger.LogQuery(ctx, "BEGIN TRANSACTION", nil, duration, err)
	
	if err != nil {
		return nil, err
	}
	
	return &LoggingTx{
		tx:       tx,
		dbLogger: ldb.dbLogger,
	}, nil
}

// PingContext pings the database with context and logging
func (ldb *LoggingDB) PingContext(ctx context.Context) error {
	start := time.Now()
	err := ldb.db.PingContext(ctx)
	duration := time.Since(start)
	
	ldb.dbLogger.LogQuery(ctx, "PING", nil, duration, err)
	return err
}

// Close closes the database connection
func (ldb *LoggingDB) Close() error {
	return ldb.db.Close()
}

// Stats returns database statistics
func (ldb *LoggingDB) Stats() sql.DBStats {
	return ldb.db.Stats()
}

// SetMaxOpenConns sets the maximum number of open connections
func (ldb *LoggingDB) SetMaxOpenConns(n int) {
	ldb.db.SetMaxOpenConns(n)
}

// SetMaxIdleConns sets the maximum number of idle connections
func (ldb *LoggingDB) SetMaxIdleConns(n int) {
	ldb.db.SetMaxIdleConns(n)
}

// SetConnMaxLifetime sets the maximum lifetime of connections
func (ldb *LoggingDB) SetConnMaxLifetime(d time.Duration) {
	ldb.db.SetConnMaxLifetime(d)
}

// SetConnMaxIdleTime sets the maximum idle time of connections
func (ldb *LoggingDB) SetConnMaxIdleTime(d time.Duration) {
	ldb.db.SetConnMaxIdleTime(d)
}

// GetUnderlyingDB returns the underlying sql.DB for cases where direct access is needed
func (ldb *LoggingDB) GetUnderlyingDB() *sql.DB {
	return ldb.db
}

// LoggingStmt wraps sql.Stmt to provide query logging
type LoggingStmt struct {
	stmt     *sql.Stmt
	query    string
	dbLogger DBLogger
}

// QueryContext executes a prepared statement query with context and logging
func (ls *LoggingStmt) QueryContext(ctx context.Context, args ...interface{}) (*sql.Rows, error) {
	start := time.Now()
	rows, err := ls.stmt.QueryContext(ctx, args...)
	duration := time.Since(start)
	
	ls.dbLogger.LogQuery(ctx, fmt.Sprintf("STMT: %s", ls.query), args, duration, err)
	return rows, err
}

// QueryRowContext executes a prepared statement query that returns a single row
func (ls *LoggingStmt) QueryRowContext(ctx context.Context, args ...interface{}) *sql.Row {
	start := time.Now()
	row := ls.stmt.QueryRowContext(ctx, args...)
	duration := time.Since(start)
	
	ls.dbLogger.LogQuery(ctx, fmt.Sprintf("STMT: %s", ls.query), args, duration, nil)
	return row
}

// ExecContext executes a prepared statement with context and logging
func (ls *LoggingStmt) ExecContext(ctx context.Context, args ...interface{}) (sql.Result, error) {
	start := time.Now()
	result, err := ls.stmt.ExecContext(ctx, args...)
	duration := time.Since(start)
	
	ls.dbLogger.LogQuery(ctx, fmt.Sprintf("STMT: %s", ls.query), args, duration, err)
	return result, err
}

// Close closes the prepared statement
func (ls *LoggingStmt) Close() error {
	return ls.stmt.Close()
}

// LoggingTx wraps sql.Tx to provide transaction logging
type LoggingTx struct {
	tx       *sql.Tx
	dbLogger DBLogger
}

// QueryContext executes a query within a transaction with context and logging
func (lt *LoggingTx) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	start := time.Now()
	rows, err := lt.tx.QueryContext(ctx, query, args...)
	duration := time.Since(start)
	
	lt.dbLogger.LogQuery(ctx, fmt.Sprintf("TX: %s", query), args, duration, err)
	return rows, err
}

// QueryRowContext executes a query that returns a single row within a transaction
func (lt *LoggingTx) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	start := time.Now()
	row := lt.tx.QueryRowContext(ctx, query, args...)
	duration := time.Since(start)
	
	lt.dbLogger.LogQuery(ctx, fmt.Sprintf("TX: %s", query), args, duration, nil)
	return row
}

// ExecContext executes a query within a transaction with context and logging
func (lt *LoggingTx) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	start := time.Now()
	result, err := lt.tx.ExecContext(ctx, query, args...)
	duration := time.Since(start)
	
	lt.dbLogger.LogQuery(ctx, fmt.Sprintf("TX: %s", query), args, duration, err)
	return result, err
}

// PrepareContext prepares a statement within a transaction
func (lt *LoggingTx) PrepareContext(ctx context.Context, query string) (*LoggingStmt, error) {
	start := time.Now()
	stmt, err := lt.tx.PrepareContext(ctx, query)
	duration := time.Since(start)
	
	lt.dbLogger.LogQuery(ctx, fmt.Sprintf("TX PREPARE: %s", query), nil, duration, err)
	
	if err != nil {
		return nil, err
	}
	
	return &LoggingStmt{
		stmt:     stmt,
		query:    query,
		dbLogger: lt.dbLogger,
	}, nil
}

// Commit commits the transaction with logging
func (lt *LoggingTx) Commit() error {
	start := time.Now()
	err := lt.tx.Commit()
	duration := time.Since(start)
	
	lt.dbLogger.LogQuery(context.Background(), "COMMIT", nil, duration, err)
	return err
}

// Rollback rolls back the transaction with logging
func (lt *LoggingTx) Rollback() error {
	start := time.Now()
	err := lt.tx.Rollback()
	duration := time.Since(start)
	
	lt.dbLogger.LogQuery(context.Background(), "ROLLBACK", nil, duration, err)
	return err
}

// NewLoggingDB creates a new logging database wrapper
func NewLoggingDB(db *sql.DB, logger Logger, slowQueryThreshold time.Duration) *LoggingDB {
	dbLogger := NewDBLogger(logger, slowQueryThreshold)
	return &LoggingDB{
		db:       db,
		dbLogger: dbLogger,
	}
}

// Helper functions

// sanitizeQuery removes extra whitespace and limits query length for logging
func sanitizeQuery(query string) string {
	// Replace multiple whitespaces with single space
	query = strings.Join(strings.Fields(query), " ")
	
	// Limit query length for logging (max 500 characters)
	if len(query) > 500 {
		query = query[:497] + "..."
	}
	
	return query
}

// sanitizeArgs sanitizes query arguments for logging (removes sensitive data)
func sanitizeArgs(args []interface{}) []interface{} {
	sanitized := make([]interface{}, len(args))
	
	for i, arg := range args {
		switch v := arg.(type) {
		case string:
			// Check if it looks like a password or sensitive data
			if isSensitiveString(v) {
				sanitized[i] = "[REDACTED]"
			} else if len(v) > 100 {
				// Truncate long strings
				sanitized[i] = v[:97] + "..."
			} else {
				sanitized[i] = v
			}
		case []byte:
			// Don't log binary data
			sanitized[i] = fmt.Sprintf("[BINARY:%d bytes]", len(v))
		default:
			sanitized[i] = arg
		}
	}
	
	return sanitized
}

// isSensitiveString checks if a string might contain sensitive information
func isSensitiveString(s string) bool {
	s = strings.ToLower(s)
	sensitivePatterns := []string{
		"password",
		"token",
		"secret",
		"key",
		"auth",
		"credential",
	}
	
	for _, pattern := range sensitivePatterns {
		if strings.Contains(s, pattern) {
			return true
		}
	}
	
	// Check if it looks like a hash (long string with mixed characters)
	if len(s) > 20 && containsMixedChars(s) {
		return true
	}
	
	return false
}

// containsMixedChars checks if string contains mixed alphanumeric characters (like a hash)
func containsMixedChars(s string) bool {
	hasLetter := false
	hasDigit := false
	
	for _, r := range s {
		if r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' {
			hasLetter = true
		} else if r >= '0' && r <= '9' {
			hasDigit = true
		}
		
		if hasLetter && hasDigit {
			return true
		}
	}
	
	return false
}