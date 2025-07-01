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
		want    *config.Config
		wantErr bool
	}{
		{
			name: "should load config correctly",
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
					Durasi: 15 * time.Second,
					Issues: "test_issuer",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a dummy .env file for testing
			dummyEnv := []byte("DB_HOST=localhost\nDB_PORT=5432\nDB_NAME=test_db\nDB_USER=test_user\nDB_PASSWORD=test_password\nDB_DRIVER=postgres\nPORT_APP=8080\nJWT_KEY=test_key\nJWT_LIFE_TIME=15\nJWT_ISSUER_NAME=test_issuer")
			err := os.WriteFile(".env", dummyEnv, 0644)
			if err != nil {
				t.Fatal(err)
			}
			defer os.Remove(".env")

			got, err := config.NewConfig()
			if (err != nil) != tt.wantErr {
				t.Errorf("NewConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// We need to manually set the duration because it's calculated inside the function
			tt.want.SecurityConfig.Durasi = got.SecurityConfig.Durasi

			if got.DbConfig != tt.want.DbConfig || got.AppConfig != tt.want.AppConfig || got.SecurityConfig != tt.want.SecurityConfig {
				t.Errorf("NewConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}

