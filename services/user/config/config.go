package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	AppPort   string `mapstructure:"APP_PORT"`
	DBHost    string `mapstructure:"MYSQL_HOST"`
	DBPort    string `mapstructure:"MYSQL_PORT"`
	DBUser    string `mapstructure:"MYSQL_USER"`
	DBPass    string `mapstructure:"MYSQL_PASSWORD"`
	DBName    string `mapstructure:"MYSQL_DATABASE"`
	JWTSecret            string `mapstructure:"JWT_SECRET"`
	JWTExpiration        int    `mapstructure:"JWT_EXPIRATION"`
	JWTRefreshExpiration int    `mapstructure:"JWT_REFRESH_EXPIRATION"`
	RedisHost            string `mapstructure:"REDIS_HOST"`
	RedisPort            string `mapstructure:"REDIS_PORT"`
	RedisPassword        string `mapstructure:"REDIS_PASSWORD"`
	GoogleClientID       string `mapstructure:"GOOGLE_CLIENT_ID"`
	GoogleClientSecret   string `mapstructure:"GOOGLE_CLIENT_SECRET"`
	GoogleRedirectURL    string `mapstructure:"GOOGLE_REDIRECT_URL"`
	FrontendCallbackURL  string `mapstructure:"FRONTEND_CALLBACK_URL"`
	RabbitMQURL          string `mapstructure:"RABBITMQ_URL"`
	MinioEndpoint        string `mapstructure:"MINIO_ENDPOINT"`
	MinioAccessKey       string `mapstructure:"MINIO_ACCESS_KEY"`
	MinioSecretKey       string `mapstructure:"MINIO_SECRET_KEY"`
	MinioBucket          string `mapstructure:"MINIO_BUCKET"`
	MinioPublicURL       string `mapstructure:"MINIO_PUBLIC_URL"` // Tambahan untuk production
	CORSAllowOrigins     string `mapstructure:"CORS_ALLOW_ORIGINS"` // Tambahan untuk keamanan CORS
}

func LoadConfig() (*Config, error) {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	// Default values (hanya yang tidak sensitif)
	viper.SetDefault("APP_PORT", "8082")
	viper.SetDefault("MYSQL_HOST", "localhost")
	viper.SetDefault("MYSQL_PORT", "3306")
	viper.SetDefault("MYSQL_USER", "root")
	viper.SetDefault("MYSQL_DATABASE", "user_db")
	viper.SetDefault("JWT_EXPIRATION", 24)
	viper.SetDefault("JWT_REFRESH_EXPIRATION", 168)
	viper.SetDefault("REDIS_HOST", "localhost")
	viper.SetDefault("REDIS_PORT", "6379")
	viper.SetDefault("GOOGLE_REDIRECT_URL", "http://localhost:8080/api/v1/auth/google/callback")
	viper.SetDefault("FRONTEND_CALLBACK_URL", "http://localhost:3000/auth/callback")
	viper.SetDefault("MINIO_ENDPOINT", "localhost:9000")
	viper.SetDefault("MINIO_BUCKET", "users")
	viper.SetDefault("MINIO_PUBLIC_URL", "http://localhost:9000") // fallback default
	viper.SetDefault("CORS_ALLOW_ORIGINS", "http://localhost:8080")

	// Secrets TIDAK BOLEH punya default — wajib dari environment
	// Tapi harus di-BindEnv agar Viper bisa baca dari OS env var saat Unmarshal
	viper.BindEnv("MYSQL_PASSWORD")
	viper.BindEnv("JWT_SECRET")
	viper.BindEnv("REDIS_PASSWORD")
	viper.BindEnv("GOOGLE_CLIENT_ID")
	viper.BindEnv("GOOGLE_CLIENT_SECRET")
	viper.BindEnv("RABBITMQ_URL")
	viper.BindEnv("MINIO_ACCESS_KEY")
	viper.BindEnv("MINIO_SECRET_KEY")

	if err := viper.ReadInConfig(); err != nil {
		// .env file not found, rely on env vars
	}

	cfg := &Config{}
	if err := viper.Unmarshal(cfg); err != nil {
		return nil, err
	}

	// Validasi wajib (Credentials & Secrets)
	if cfg.JWTSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET is required but not set")
	}
	if len(cfg.JWTSecret) < 32 {
		return nil, fmt.Errorf("JWT_SECRET must be at least 32 characters long for security")
	}
	if cfg.DBPass == "" {
		return nil, fmt.Errorf("MYSQL_PASSWORD is required but not set")
	}
	if cfg.RabbitMQURL == "" {
		return nil, fmt.Errorf("RABBITMQ_URL is required but not set")
	}
	if cfg.MinioAccessKey == "" {
		return nil, fmt.Errorf("MINIO_ACCESS_KEY is required but not set")
	}
	if cfg.MinioSecretKey == "" {
		return nil, fmt.Errorf("MINIO_SECRET_KEY is required but not set")
	}

	return cfg, nil
}
