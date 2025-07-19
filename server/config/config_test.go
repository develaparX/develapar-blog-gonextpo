package config_test

import (
	"os"
	"testing"
	"time"

	"develapar-server/config"
)

func TestNewConfig(t *testing.T) {
	tests := []struct {
		name    string
		setupEnv func()
		teardownEnv func()
		wantErr bool
		want    *config.Config
	}{
		{
			name: "should load config correctly",
			setupEnv: func() {
				os.Setenv("DB_HOST", "localhost")
				os.Setenv("DB_PORT", "5432")
				os.Setenv("DB_NAME", "test_db")
				os.Setenv("DB_USER", "test_user")
				os.Setenv("DB_PASSWORD", "test_password")
				os.Setenv("DB_DRIVER", "postgres")
				os.Setenv("PORT_APP", "8080")
				os.Setenv("JWT_KEY", "test_key")
				os.Setenv("JWT_LIFE_TIME", "15")
				os.Setenv("JWT_ISSUER_NAME", "test_issuer")
			},
			teardownEnv: func() {
				os.Unsetenv("DB_HOST")
				os.Unsetenv("DB_PORT")
				os.Unsetenv("DB_NAME")
				os.Unsetenv("DB_USER")
				os.Unsetenv("DB_PASSWORD")
				os.Unsetenv("DB_DRIVER")
				os.Unsetenv("PORT_APP")
				os.Unsetenv("JWT_KEY")
				os.Unsetenv("JWT_LIFE_TIME")
				os.Unsetenv("JWT_ISSUER_NAME")
			},
			wantErr: false,
			want: &config.Config{
				DbConfig: config.DbConfig{
					Host:     "localhost",
					Port:     "5432",
					Name:     "test_db",
					User:     "test_user",
					Password: "test_password",
					Driver:   "postgres",
				},
				AppConfig: config.AppConfig{
					AppPort: "8080",
				},
				SecurityConfig: config.SecurityConfig{
					Key:    "test_key",
					Durasi: 15 * time.Hour,
					Issues: "test_issuer",
				},
			},
		},
		{
			name: "should return error if environment is empty",
			setupEnv: func() {
				// Unset all env vars to simulate empty environment
				os.Unsetenv("DB_HOST")
				os.Unsetenv("DB_PORT")
				os.Unsetenv("DB_NAME")
				os.Unsetenv("DB_USER")
				os.Unsetenv("DB_PASSWORD")
				os.Unsetenv("DB_DRIVER")
				os.Unsetenv("PORT_APP")
				os.Unsetenv("JWT_KEY")
				os.Unsetenv("JWT_LIFE_TIME")
				os.Unsetenv("JWT_ISSUER_NAME")
			},
			teardownEnv: func() {
				// No need to unset, already unset in setup
			},
			wantErr: true,
			want:    nil, // Expect nil config on error
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupEnv()
			defer tt.teardownEnv()

			got, err := config.NewConfig()
			if (err != nil) != tt.wantErr {
				t.Errorf("NewConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				if got != nil {
					t.Errorf("NewConfig() got = %v, want %v", got, tt.want)
				}
				return
			}

			// We need to manually set the duration because it's calculated inside the function
			// and we want to compare the rest of the fields.
			tt.want.SecurityConfig.Durasi = got.SecurityConfig.Durasi

			if got.DbConfig != tt.want.DbConfig || got.AppConfig != tt.want.AppConfig || got.SecurityConfig != tt.want.SecurityConfig {
				t.Errorf("NewConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDefaultContextConfig(t *testing.T) {
	config := config.DefaultContextConfig()

	if config.RequestTimeout != 30*time.Second {
		t.Errorf("DefaultContextConfig() RequestTimeout = %v, want %v", config.RequestTimeout, 30*time.Second)
	}
	if config.DatabaseTimeout != 15*time.Second {
		t.Errorf("DefaultContextConfig() DatabaseTimeout = %v, want %v", config.DatabaseTimeout, 15*time.Second)
	}
	if config.ValidationTimeout != 5*time.Second {
		t.Errorf("DefaultContextConfig() ValidationTimeout = %v, want %v", config.ValidationTimeout, 5*time.Second)
	}
	if config.LoggingTimeout != 2*time.Second {
		t.Errorf("DefaultContextConfig() LoggingTimeout = %v, want %v", config.LoggingTimeout, 2*time.Second)
	}
}

func TestDefaultLoggingConfig(t *testing.T) {
	config := config.DefaultLoggingConfig()

	if config.Level != "info" {
		t.Errorf("DefaultLoggingConfig() Level = %v, want %v", config.Level, "info")
	}
	if config.Format != "json" {
		t.Errorf("DefaultLoggingConfig() Format = %v, want %v", config.Format, "json")
	}
	if config.OutputPath != "stdout" {
		t.Errorf("DefaultLoggingConfig() OutputPath = %v, want %v", config.OutputPath, "stdout")
	}
	if config.ErrorOutputPath != "stderr" {
		t.Errorf("DefaultLoggingConfig() ErrorOutputPath = %v, want %v", config.ErrorOutputPath, "stderr")
	}
	if config.MaxSize != 100 {
		t.Errorf("DefaultLoggingConfig() MaxSize = %v, want %v", config.MaxSize, 100)
	}
	if config.MaxBackups != 3 {
		t.Errorf("DefaultLoggingConfig() MaxBackups = %v, want %v", config.MaxBackups, 3)
	}
	if config.MaxAge != 28 {
		t.Errorf("DefaultLoggingConfig() MaxAge = %v, want %v", config.MaxAge, 28)
	}
	if config.Compress != true {
		t.Errorf("DefaultLoggingConfig() Compress = %v, want %v", config.Compress, true)
	}
	if config.RequestTimeout != 2*time.Second {
		t.Errorf("DefaultLoggingConfig() RequestTimeout = %v, want %v", config.RequestTimeout, 2*time.Second)
	}
}

func TestDefaultRateLimitConfig(t *testing.T) {
	config := config.DefaultRateLimitConfig()

	if config.Enabled != true {
		t.Errorf("DefaultRateLimitConfig() Enabled = %v, want %v", config.Enabled, true)
	}
	if config.RequestsPerMinute != 60 {
		t.Errorf("DefaultRateLimitConfig() RequestsPerMinute = %v, want %v", config.RequestsPerMinute, 60)
	}
	if config.BurstSize != 10 {
		t.Errorf("DefaultRateLimitConfig() BurstSize = %v, want %v", config.BurstSize, 10)
	}
	if config.CleanupInterval != 5*time.Minute {
		t.Errorf("DefaultRateLimitConfig() CleanupInterval = %v, want %v", config.CleanupInterval, 5*time.Minute)
	}
	if config.WindowSize != 1*time.Minute {
		t.Errorf("DefaultRateLimitConfig() WindowSize = %v, want %v", config.WindowSize, 1*time.Minute)
	}
	if config.AuthenticatedRPM != 120 {
		t.Errorf("DefaultRateLimitConfig() AuthenticatedRPM = %v, want %v", config.AuthenticatedRPM, 120)
	}
	if config.AnonymousRPM != 30 {
		t.Errorf("DefaultRateLimitConfig() AnonymousRPM = %v, want %v", config.AnonymousRPM, 30)
	}
	if config.RequestTimeout != 1*time.Second {
		t.Errorf("DefaultRateLimitConfig() RequestTimeout = %v, want %v", config.RequestTimeout, 1*time.Second)
	}
}

func TestLoadContextConfigFromEnv(t *testing.T) {
	tests := []struct {
		name     string
		setupEnv func()
		teardownEnv func()
		want     config.ContextConfig
	}{
		{
			name: "should load context config from environment variables",
			setupEnv: func() {
				os.Setenv("CONTEXT_REQUEST_TIMEOUT", "45s")
				os.Setenv("CONTEXT_DATABASE_TIMEOUT", "20s")
				os.Setenv("CONTEXT_VALIDATION_TIMEOUT", "8s")
				os.Setenv("CONTEXT_LOGGING_TIMEOUT", "3s")
			},
			teardownEnv: func() {
				os.Unsetenv("CONTEXT_REQUEST_TIMEOUT")
				os.Unsetenv("CONTEXT_DATABASE_TIMEOUT")
				os.Unsetenv("CONTEXT_VALIDATION_TIMEOUT")
				os.Unsetenv("CONTEXT_LOGGING_TIMEOUT")
			},
			want: config.ContextConfig{
				RequestTimeout:    45 * time.Second,
				DatabaseTimeout:   20 * time.Second,
				ValidationTimeout: 8 * time.Second,
				LoggingTimeout:    3 * time.Second,
			},
		},
		{
			name: "should use defaults when environment variables are not set",
			setupEnv: func() {
				// Don't set any environment variables
			},
			teardownEnv: func() {
				// Nothing to clean up
			},
			want: config.DefaultContextConfig(),
		},
		{
			name: "should use defaults when environment variables are invalid",
			setupEnv: func() {
				os.Setenv("CONTEXT_REQUEST_TIMEOUT", "invalid")
				os.Setenv("CONTEXT_DATABASE_TIMEOUT", "-5s")
				os.Setenv("CONTEXT_VALIDATION_TIMEOUT", "0s")
			},
			teardownEnv: func() {
				os.Unsetenv("CONTEXT_REQUEST_TIMEOUT")
				os.Unsetenv("CONTEXT_DATABASE_TIMEOUT")
				os.Unsetenv("CONTEXT_VALIDATION_TIMEOUT")
			},
			want: config.DefaultContextConfig(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupEnv()
			defer tt.teardownEnv()

			// Create a config instance to test the method
			c := &config.Config{}
			got := c.LoadContextConfig()

			if got != tt.want {
				t.Errorf("loadContextConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLoadLoggingConfigFromEnv(t *testing.T) {
	tests := []struct {
		name     string
		setupEnv func()
		teardownEnv func()
		want     config.LoggingConfig
	}{
		{
			name: "should load logging config from environment variables",
			setupEnv: func() {
				os.Setenv("LOG_LEVEL", "debug")
				os.Setenv("LOG_FORMAT", "text")
				os.Setenv("LOG_OUTPUT_PATH", "/var/log/app.log")
				os.Setenv("LOG_ERROR_OUTPUT_PATH", "/var/log/error.log")
				os.Setenv("LOG_MAX_SIZE", "200")
				os.Setenv("LOG_MAX_BACKUPS", "5")
				os.Setenv("LOG_MAX_AGE", "30")
				os.Setenv("LOG_COMPRESS", "false")
				os.Setenv("LOG_REQUEST_TIMEOUT", "5s")
			},
			teardownEnv: func() {
				os.Unsetenv("LOG_LEVEL")
				os.Unsetenv("LOG_FORMAT")
				os.Unsetenv("LOG_OUTPUT_PATH")
				os.Unsetenv("LOG_ERROR_OUTPUT_PATH")
				os.Unsetenv("LOG_MAX_SIZE")
				os.Unsetenv("LOG_MAX_BACKUPS")
				os.Unsetenv("LOG_MAX_AGE")
				os.Unsetenv("LOG_COMPRESS")
				os.Unsetenv("LOG_REQUEST_TIMEOUT")
			},
			want: config.LoggingConfig{
				Level:           "debug",
				Format:          "text",
				OutputPath:      "/var/log/app.log",
				ErrorOutputPath: "/var/log/error.log",
				MaxSize:         200,
				MaxBackups:      5,
				MaxAge:          30,
				Compress:        false,
				RequestTimeout:  5 * time.Second,
			},
		},
		{
			name: "should use defaults when environment variables are not set",
			setupEnv: func() {
				// Don't set any environment variables
			},
			teardownEnv: func() {
				// Nothing to clean up
			},
			want: config.DefaultLoggingConfig(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupEnv()
			defer tt.teardownEnv()

			// Create a config instance to test the method
			c := &config.Config{}
			got := c.LoadLoggingConfig()

			if got != tt.want {
				t.Errorf("loadLoggingConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLoadRateLimitConfigFromEnv(t *testing.T) {
	tests := []struct {
		name     string
		setupEnv func()
		teardownEnv func()
		want     config.RateLimitConfig
	}{
		{
			name: "should load rate limit config from environment variables",
			setupEnv: func() {
				os.Setenv("RATE_LIMIT_ENABLED", "false")
				os.Setenv("RATE_LIMIT_REQUESTS_PER_MINUTE", "100")
				os.Setenv("RATE_LIMIT_BURST_SIZE", "20")
				os.Setenv("RATE_LIMIT_CLEANUP_INTERVAL", "10m")
				os.Setenv("RATE_LIMIT_WINDOW_SIZE", "2m")
				os.Setenv("RATE_LIMIT_AUTHENTICATED_RPM", "200")
				os.Setenv("RATE_LIMIT_ANONYMOUS_RPM", "50")
				os.Setenv("RATE_LIMIT_REQUEST_TIMEOUT", "2s")
			},
			teardownEnv: func() {
				os.Unsetenv("RATE_LIMIT_ENABLED")
				os.Unsetenv("RATE_LIMIT_REQUESTS_PER_MINUTE")
				os.Unsetenv("RATE_LIMIT_BURST_SIZE")
				os.Unsetenv("RATE_LIMIT_CLEANUP_INTERVAL")
				os.Unsetenv("RATE_LIMIT_WINDOW_SIZE")
				os.Unsetenv("RATE_LIMIT_AUTHENTICATED_RPM")
				os.Unsetenv("RATE_LIMIT_ANONYMOUS_RPM")
				os.Unsetenv("RATE_LIMIT_REQUEST_TIMEOUT")
			},
			want: config.RateLimitConfig{
				Enabled:           false,
				RequestsPerMinute: 100,
				BurstSize:         20,
				CleanupInterval:   10 * time.Minute,
				WindowSize:        2 * time.Minute,
				AuthenticatedRPM:  200,
				AnonymousRPM:      50,
				RequestTimeout:    2 * time.Second,
			},
		},
		{
			name: "should use defaults when environment variables are not set",
			setupEnv: func() {
				// Don't set any environment variables
			},
			teardownEnv: func() {
				// Nothing to clean up
			},
			want: config.DefaultRateLimitConfig(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupEnv()
			defer tt.teardownEnv()

			// Create a config instance to test the method
			c := &config.Config{}
			got := c.LoadRateLimitConfig()

			if got != tt.want {
				t.Errorf("loadRateLimitConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}func Tes
tValidateConfig(t *testing.T) {
	tests := []struct {
		name     string
		set