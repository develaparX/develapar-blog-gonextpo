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

type Config struct {
	DbConfig
	AppConfig
	SecurityConfig
	PoolConfig
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

	if c.DbConfig.Host == "" || c.DbConfig.Port == "" || c.DbConfig.Name == "" || c.DbConfig.User == "" || c.DbConfig.Password == "" || c.DbConfig.Driver == "" || c.SecurityConfig.Key == "" || c.SecurityConfig.Durasi < 0 || c.SecurityConfig.Issues == "" {
		return errors.New("environtment is empty")
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

func NewConfig() (*Config, error) {

	config := &Config{}
	if err := config.readConfig(); err != nil {
		return nil, err
	}

	return config, nil
}
