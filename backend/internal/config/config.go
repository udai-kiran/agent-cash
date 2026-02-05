package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

// Config holds all configuration for the application
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
	CORS     CORSConfig
}

// CORSConfig holds CORS configuration
type CORSConfig struct {
	AllowedOrigins []string `mapstructure:"allowedOrigins"`
}

// ServerConfig holds HTTP server configuration
type ServerConfig struct {
	Port         int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

// DatabaseConfig holds database connection configuration
type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
	MaxConns int32
	MinConns int32
}

// JWTConfig holds JWT configuration
type JWTConfig struct {
	Secret           string
	AccessTokenTTL   time.Duration
	RefreshTokenTTL  time.Duration
}

// ConnectionString returns a PostgreSQL connection string
func (d *DatabaseConfig) ConnectionString() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		d.Host, d.Port, d.User, d.Password, d.DBName, d.SSLMode,
	)
}

// Load loads configuration from file and environment variables
func Load(configPath string) (*Config, error) {
	// Set defaults
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("server.readTimeout", "15s")
	viper.SetDefault("server.writeTimeout", "15s")
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", 5432)
	viper.SetDefault("database.user", "gnucash")
	viper.SetDefault("database.password", "gnucash")
	viper.SetDefault("database.dbname", "gnucash")
	viper.SetDefault("database.sslmode", "disable")
	viper.SetDefault("database.maxConns", 10)
	viper.SetDefault("database.minConns", 2)
	viper.SetDefault("cors.allowedOrigins", []string{"http://localhost:3000"})
	viper.SetDefault("jwt.secret", "")
	viper.SetDefault("jwt.accessTokenTTL", "15m")
	viper.SetDefault("jwt.refreshTokenTTL", "168h") // 7 days

	// Enable environment variable override
	viper.AutomaticEnv()
	viper.SetEnvPrefix("") // Allow env vars without prefix

	// Explicitly bind env vars so Unmarshal picks them up
	viper.BindEnv("database.host", "DATABASE_HOST")
	viper.BindEnv("database.port", "DATABASE_PORT")
	viper.BindEnv("database.user", "DATABASE_USER")
	viper.BindEnv("database.password", "DATABASE_PASSWORD")
	viper.BindEnv("database.dbname", "DATABASE_NAME")
	viper.BindEnv("database.sslmode", "DATABASE_SSLMODE")
	viper.BindEnv("server.port", "SERVER_PORT")
	viper.BindEnv("jwt.secret", "JWT_SECRET")
	viper.BindEnv("cors.allowedOrigins", "CORS_ALLOWED_ORIGINS")

	// Try to read config file (optional in Docker)
	if configPath != "" {
		viper.SetConfigFile(configPath)
		if err := viper.ReadInConfig(); err != nil {
			// Config file is optional if env vars are provided
			fmt.Printf("Warning: Could not read config file: %v\n", err)
		}
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	if len(config.JWT.Secret) < 32 {
		return nil, fmt.Errorf("JWT_SECRET must be at least 32 characters (set via JWT_SECRET env var)")
	}

	return &config, nil
}
