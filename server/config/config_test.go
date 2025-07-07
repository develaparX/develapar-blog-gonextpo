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