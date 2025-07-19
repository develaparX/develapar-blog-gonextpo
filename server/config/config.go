package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type DbConfig struct {
	Host     string
	Port     string
	Name     string
	User     string
	Password string
	Driver   string
}

type AppConfig struct {
	AppPort string
}

type SecurityConfig struct {
	Key    string
	Durasi time.Duration
	Issues string
}

type ContextConfig struct {
	RequestTimeout    time.Duration `json:"request_timeout"`
	DatabaseTimeout   time.Duration `json:"database_timeout"`
	ValidationTimeout time.Duration `json:"validation_timeout"`
	LoggingTimeout    time.Duration `json:"logging_timeout"`
}

type LoggingConfig struct {
	Level           string        `json:"level"`
	Format          string        `json:"format"`
	OutputPath      string        `json:"output_path"`
	ErrorOutputPath string        `json:"error_output_path"`
	MaxSize         int           `json:"max_size"`
	MaxBackups      int           `json:"max_backups"`
	MaxAge          int           `json:"max_age"`
	Compress        bool          `json:"compress"`
	RequestTimeout  time.Duration `json:"request_timeout"`
}

type RateLimitConfig struct {
	Enabled           bool          `json:"enabled"`
	RequestsPerMinute int           `json:"requests_per_minute"`
	BurstSize         int           `json:"burst_size"`
	CleanupInterval   time.Duration `json:"cleanup_interval"`
	WindowSize        time.Duration `json:"window_size"`
	AuthenticatedRPM  int           `json:"authenticated_rpm"`
	AnonymousRPM      int           `json:"anonymous_rpm"`
	RequestTimeout    time.Duration `json:"request_timeout"`
}

type Config struct {
	DbConfig
	AppConfig
	SecurityConfig
	PoolConfig
	ContextConfig
	LoggingConfig
	RateLimitConfig
}

func (c *Config) readConfig() error {
	err := godotenv.Load()
	if err != nil {
		fmt.Println(err)
	}

	lifeTime, err := strconv.Atoi(os.Getenv("JWT_LIFE_TIME"))
	if err != nil {
		return err

	}

	c.SecurityConfig = SecurityConfig{
		Key:    os.Getenv("JWT_KEY"),
		Durasi: time.Duration(lifeTime) * time.Hour,
		Issues: os.Getenv("JWT_ISSUER_NAME"),
	}

	c.DbConfig = DbConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Name:     os.Getenv("DB_NAME"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Driver:   os.Getenv("DB_DRIVER"),
	}

	c.AppConfig = AppConfig{
		AppPort: os.Getenv("PORT_APP"),
	}

	// Load pool configuration with defaults
	c.PoolConfig = c.loadPoolConfig()

	// Load context configuration with defaults
	c.ContextConfig = c.loadContextConfig()

	// Load logging configuration with defaults
	c.LoggingConfig = c.loadLoggingConfig()

	// Load rate limiting configuration with defaults
	c.RateLimitConfig = c.loadRateLimitConfig()

	// Validate required configuration fields
	if err := c.validateConfig(); err != nil {
		return err
	}
	
	return nil

}

func (c *Config) loadPoolConfig() PoolConfig {
	// Start with default configuration
	poolConfig := DefaultPoolConfig()

	// Override with environment variables if present
	if maxOpenConns := os.Getenv("DB_MAX_OPEN_CONNS"); maxOpenConns != "" {
		if val, err := strconv.Atoi(maxOpenConns); err == nil && val > 0 {
			poolConfig.MaxOpenConns = val
		}
	}

	if maxIdleConns := os.Getenv("DB_MAX_IDLE_CONNS"); maxIdleConns != "" {
		if val, err := strconv.Atoi(maxIdleConns); err == nil && val > 0 {
			poolConfig.MaxIdleConns = val
		}
	}

	if connMaxLifetime := os.Getenv("DB_CONN_MAX_LIFETIME"); connMaxLifetime != "" {
		if val, err := time.ParseDuration(connMaxLifetime); err == nil && val > 0 {
			poolConfig.ConnMaxLifetime = val
		}
	}

	if connMaxIdleTime := os.Getenv("DB_CONN_MAX_IDLE_TIME"); connMaxIdleTime != "" {
		if val, err := time.ParseDuration(connMaxIdleTime); err == nil && val > 0 {
			poolConfig.ConnMaxIdleTime = val
		}
	}

	if connectTimeout := os.Getenv("DB_CONNECT_TIMEOUT"); connectTimeout != "" {
		if val, err := time.ParseDuration(connectTimeout); err == nil && val > 0 {
			poolConfig.ConnectTimeout = val
		}
	}

	if queryTimeout := os.Getenv("DB_QUERY_TIMEOUT"); queryTimeout != "" {
		if val, err := time.ParseDuration(queryTimeout); err == nil && val > 0 {
			poolConfig.QueryTimeout = val
		}
	}

	return poolConfig
}

func (c *Config) loadContextConfig() ContextConfig {
	// Start with default configuration
	contextConfig := DefaultContextConfig()

	// Override with environment variables if present
	if requestTimeout := os.Getenv("CONTEXT_REQUEST_TIMEOUT"); requestTimeout != "" {
		if val, err := time.ParseDuration(requestTimeout); err == nil && val > 0 {
			contextConfig.RequestTimeout = val
		}
	}

	if databaseTimeout := os.Getenv("CONTEXT_DATABASE_TIMEOUT"); databaseTimeout != "" {
		if val, err := time.ParseDuration(databaseTimeout); err == nil && val > 0 {
			contextConfig.DatabaseTimeout = val
		}
	}

	if validationTimeout := os.Getenv("CONTEXT_VALIDATION_TIMEOUT"); validationTimeout != "" {
		if val, err := time.ParseDuration(validationTimeout); err == nil && val > 0 {
			contextConfig.ValidationTimeout = val
		}
	}

	if loggingTimeout := os.Getenv("CONTEXT_LOGGING_TIMEOUT"); loggingTimeout != "" {
		if val, err := time.ParseDuration(loggingTimeout); err == nil && val > 0 {
			contextConfig.LoggingTimeout = val
		}
	}

	return contextConfig
}

func (c *Config) loadLoggingConfig() LoggingConfig {
	// Start with default configuration
	loggingConfig := DefaultLoggingConfig()

	// Override with environment variables if present
	if level := os.Getenv("LOG_LEVEL"); level != "" {
		loggingConfig.Level = level
	}

	if format := os.Getenv("LOG_FORMAT"); format != "" {
		loggingConfig.Format = format
	}

	if outputPath := os.Getenv("LOG_OUTPUT_PATH"); outputPath != "" {
		loggingConfig.OutputPath = outputPath
	}

	if errorOutputPath := os.Getenv("LOG_ERROR_OUTPUT_PATH"); errorOutputPath != "" {
		loggingConfig.ErrorOutputPath = errorOutputPath
	}

	if maxSize := os.Getenv("LOG_MAX_SIZE"); maxSize != "" {
		if val, err := strconv.Atoi(maxSize); err == nil && val > 0 {
			loggingConfig.MaxSize = val
		}
	}

	if maxBackups := os.Getenv("LOG_MAX_BACKUPS"); maxBackups != "" {
		if val, err := strconv.Atoi(maxBackups); err == nil && val >= 0 {
			loggingConfig.MaxBackups = val
		}
	}

	if maxAge := os.Getenv("LOG_MAX_AGE"); maxAge != "" {
		if val, err := strconv.Atoi(maxAge); err == nil && val >= 0 {
			loggingConfig.MaxAge = val
		}
	}

	if compress := os.Getenv("LOG_COMPRESS"); compress != "" {
		if val, err := strconv.ParseBool(compress); err == nil {
			loggingConfig.Compress = val
		}
	}

	if requestTimeout := os.Getenv("LOG_REQUEST_TIMEOUT"); requestTimeout != "" {
		if val, err := time.ParseDuration(requestTimeout); err == nil && val > 0 {
			loggingConfig.RequestTimeout = val
		}
	}

	return loggingConfig
}

func (c *Config) loadRateLimitConfig() RateLimitConfig {
	// Start with default configuration
	rateLimitConfig := DefaultRateLimitConfig()

	// Override with environment variables if present
	if enabled := os.Getenv("RATE_LIMIT_ENABLED"); enabled != "" {
		if val, err := strconv.ParseBool(enabled); err == nil {
			rateLimitConfig.Enabled = val
		}
	}

	if requestsPerMinute := os.Getenv("RATE_LIMIT_REQUESTS_PER_MINUTE"); requestsPerMinute != "" {
		if val, err := strconv.Atoi(requestsPerMinute); err == nil && val > 0 {
			rateLimitConfig.RequestsPerMinute = val
		}
	}

	if burstSize := os.Getenv("RATE_LIMIT_BURST_SIZE"); burstSize != "" {
		if val, err := strconv.Atoi(burstSize); err == nil && val > 0 {
			rateLimitConfig.BurstSize = val
		}
	}

	if cleanupInterval := os.Getenv("RATE_LIMIT_CLEANUP_INTERVAL"); cleanupInterval != "" {
		if val, err := time.ParseDuration(cleanupInterval); err == nil && val > 0 {
			rateLimitConfig.CleanupInterval = val
		}
	}

	if windowSize := os.Getenv("RATE_LIMIT_WINDOW_SIZE"); windowSize != "" {
		if val, err := time.ParseDuration(windowSize); err == nil && val > 0 {
			rateLimitConfig.WindowSize = val
		}
	}

	if authenticatedRPM := os.Getenv("RATE_LIMIT_AUTHENTICATED_RPM"); authenticatedRPM != "" {
		if val, err := strconv.Atoi(authenticatedRPM); err == nil && val > 0 {
			rateLimitConfig.AuthenticatedRPM = val
		}
	}

	if anonymousRPM := os.Getenv("RATE_LIMIT_ANONYMOUS_RPM"); anonymousRPM != "" {
		if val, err := strconv.Atoi(anonymousRPM); err == nil && val > 0 {
			rateLimitConfig.AnonymousRPM = val
		}
	}

	if requestTimeout := os.Getenv("RATE_LIMIT_REQUEST_TIMEOUT"); requestTimeout != "" {
		if val, err := time.ParseDuration(requestTimeout); err == nil && val > 0 {
			rateLimitConfig.RequestTimeout = val
		}
	}

	return rateLimitConfig
}

// DefaultContextConfig returns a default context configuration
func DefaultContextConfig() ContextConfig {
	return ContextConfig{
		RequestTimeout:    30 * time.Second,  // Maximum request processing time
		DatabaseTimeout:   15 * time.Second,  // Maximum database operation time
		ValidationTimeout: 5 * time.Second,   // Maximum validation processing time
		LoggingTimeout:    2 * time.Second,   // Maximum logging operation time
	}
}

// DefaultLoggingConfig returns a default logging configuration
func DefaultLoggingConfig() LoggingConfig {
	return LoggingConfig{
		Level:           "info",
		Format:          "json",
		OutputPath:      "stdout",
		ErrorOutputPath: "stderr",
		MaxSize:         100,                // 100 MB
		MaxBackups:      3,                  // Keep 3 backup files
		MaxAge:          28,                 // Keep logs for 28 days
		Compress:        true,               // Compress rotated logs
		RequestTimeout:  2 * time.Second,    // Logging operation timeout
	}
}

// DefaultRateLimitConfig returns a default rate limiting configuration
func DefaultRateLimitConfig() RateLimitConfig {
	return RateLimitConfig{
		Enabled:           true,
		RequestsPerMinute: 60,                // 60 requests per minute for general use
		BurstSize:         10,                // Allow burst of 10 requests
		CleanupInterval:   5 * time.Minute,   // Clean up old entries every 5 minutes
		WindowSize:        1 * time.Minute,   // 1 minute sliding window
		AuthenticatedRPM:  120,               // 120 requests per minute for authenticated users
		AnonymousRPM:      30,                // 30 requests per minute for anonymous users
		RequestTimeout:    1 * time.Second,   // Rate limit check timeout
	}
}

// LoadContextConfig loads context configuration from environment variables (public for testing)
func (c *Config) LoadContextConfig() ContextConfig {
	return c.loadContextConfig()
}

// LoadLoggingConfig loads logging configuration from environment variables (public for testing)
func (c *Config) LoadLoggingConfig() LoggingConfig {
	return c.loadLoggingConfig()
}

// LoadRateLimitConfig loads rate limiting configuration from environment variables (public for testing)
func (c *Config) LoadRateLimitConfig() RateLimitConfig {
	return c.loadRateLimitConfig()
}

// validateConfig validates all configuration settings including context settings
func (c *Config) validateConfig() error {
	// Validate database configuration
	if c.DbConfig.Host == "" {
		return errors.New("database host is required")
	}
	if c.DbConfig.Port == "" {
		return errors.New("database port is required")
	}
	if c.DbConfig.Name == "" {
		return errors.New("database name is required")
	}
	if c.DbConfig.User == "" {
		return errors.New("database user is required")
	}
	if c.DbConfig.Password == "" {
		return errors.New("database password is required")
	}
	if c.DbConfig.Driver == "" {
		return errors.New("database driver is required")
	}

	// Validate security configuration
	if c.SecurityConfig.Key == "" {
		return errors.New("JWT key is required")
	}
	if c.SecurityConfig.Durasi < 0 {
		return errors.New("JWT duration must be positive")
	}
	if c.SecurityConfig.Issues == "" {
		return errors.New("JWT issuer name is required")
	}

	// Validate app configuration
	if c.AppConfig.AppPort == "" {
		return errors.New("application port is required")
	}

	// Validate context configuration
	if c.ContextConfig.RequestTimeout <= 0 {
		return errors.New("context request timeout must be positive")
	}
	if c.ContextConfig.DatabaseTimeout <= 0 {
		return errors.New("context database timeout must be positive")
	}
	if c.ContextConfig.ValidationTimeout <= 0 {
		return errors.New("context validation timeout must be positive")
	}
	if c.ContextConfig.LoggingTimeout <= 0 {
		return errors.New("context logging timeout must be positive")
	}

	// Validate logging configuration
	if c.LoggingConfig.Level == "" {
		return errors.New("logging level is required")
	}
	if c.LoggingConfig.Format == "" {
		return errors.New("logging format is required")
	}
	if c.LoggingConfig.MaxSize <= 0 {
		return errors.New("logging max size must be positive")
	}
	if c.LoggingConfig.MaxBackups < 0 {
		return errors.New("logging max backups must be non-negative")
	}
	if c.LoggingConfig.MaxAge < 0 {
		return errors.New("logging max age must be non-negative")
	}
	if c.LoggingConfig.RequestTimeout <= 0 {
		return errors.New("logging request timeout must be positive")
	}

	// Validate rate limiting configuration
	if c.RateLimitConfig.RequestsPerMinute <= 0 {
		return errors.New("rate limit requests per minute must be positive")
	}
	if c.RateLimitConfig.BurstSize <= 0 {
		return errors.New("rate limit burst size must be positive")
	}
	if c.RateLimitConfig.CleanupInterval <= 0 {
		return errors.New("rate limit cleanup interval must be positive")
	}
	if c.RateLimitConfig.WindowSize <= 0 {
		return errors.New("rate limit window size must be positive")
	}
	if c.RateLimitConfig.AuthenticatedRPM <= 0 {
		return errors.New("rate limit authenticated RPM must be positive")
	}
	if c.RateLimitConfig.AnonymousRPM <= 0 {
		return errors.New("rate limit anonymous RPM must be positive")
	}
	if c.RateLimitConfig.RequestTimeout <= 0 {
		return errors.New("rate limit request timeout must be positive")
	}

	// Validate pool configuration
	if c.PoolConfig.MaxOpenConns <= 0 {
		return errors.New("database max open connections must be positive")
	}
	if c.PoolConfig.MaxIdleConns < 0 {
		return errors.New("database max idle connections must be non-negative")
	}
	if c.PoolConfig.ConnMaxLifetime <= 0 {
		return errors.New("database connection max lifetime must be positive")
	}
	if c.PoolConfig.ConnMaxIdleTime <= 0 {
		return errors.New("database connection max idle time must be positive")
	}
	if c.PoolConfig.ConnectTimeout <= 0 {
		return errors.New("database connect timeout must be positive")
	}
	if c.PoolConfig.QueryTimeout <= 0 {
		return errors.New("database query timeout must be positive")
	}

	return nil
}

func NewConfig() (*Config, error) {

	config := &Config{}
	if err := config.readConfig(); err != nil {
		return nil, err
	}

	return config, nil
}
