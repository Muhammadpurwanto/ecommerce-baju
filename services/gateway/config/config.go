package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	AppPort             string `mapstructure:"APP_PORT"`
	JWTSecret           string `mapstructure:"JWT_SECRET"`
	AuthServiceURL         string `mapstructure:"AUTH_SERVICE_URL"`
	AuthServiceHttpURL     string `mapstructure:"AUTH_SERVICE_HTTP_URL"`
	UserServiceURL         string `mapstructure:"USER_SERVICE_URL"`
	ProductServiceURL      string `mapstructure:"PRODUCT_SERVICE_URL"`
	ProductServiceHttpURL  string `mapstructure:"PRODUCT_SERVICE_HTTP_URL"`
	CartServiceURL         string `mapstructure:"CART_SERVICE_URL"`
	OrderServiceURL        string `mapstructure:"ORDER_SERVICE_URL"`
	PaymentServiceURL      string `mapstructure:"PAYMENT_SERVICE_URL"`
	RedisHost              string `mapstructure:"REDIS_HOST"`
	RedisPort              string `mapstructure:"REDIS_PORT"`
	RedisPassword          string `mapstructure:"REDIS_PASSWORD"`
	RateLimitMax           int    `mapstructure:"RATE_LIMIT_MAX"`
	RateLimitExpiration    int    `mapstructure:"RATE_LIMIT_EXPIRATION"`
	CORSAllowOrigins       string `mapstructure:"CORS_ALLOW_ORIGINS"`
}

func LoadConfig() (*Config, error) {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	// Port & networking defaults (non-sensitive)
	viper.SetDefault("APP_PORT", "8080")
	viper.SetDefault("AUTH_SERVICE_URL", "localhost:50051")
	viper.SetDefault("AUTH_SERVICE_HTTP_URL", "http://localhost:8081")
	viper.SetDefault("USER_SERVICE_URL", "localhost:50052")
	viper.SetDefault("PRODUCT_SERVICE_URL", "localhost:50053")
	viper.SetDefault("PRODUCT_SERVICE_HTTP_URL", "http://localhost:8083")
	viper.SetDefault("CART_SERVICE_URL", "localhost:50054")
	viper.SetDefault("ORDER_SERVICE_URL", "localhost:50055")
	viper.SetDefault("PAYMENT_SERVICE_URL", "localhost:50056")
	viper.SetDefault("REDIS_HOST", "localhost")
	viper.SetDefault("REDIS_PORT", "6379")
	viper.SetDefault("RATE_LIMIT_MAX", 100)
	viper.SetDefault("RATE_LIMIT_EXPIRATION", 60) // in seconds
	viper.SetDefault("CORS_ALLOW_ORIGINS", "http://localhost:3000")

	// Secrets TIDAK BOLEH punya default — wajib dari environment
	// Tapi harus di-BindEnv agar Viper bisa baca dari OS env var saat Unmarshal
	viper.BindEnv("JWT_SECRET")
	viper.BindEnv("REDIS_PASSWORD")

	if err := viper.ReadInConfig(); err != nil {
		// Abaikan jika .env tidak ada, env vars tetap terbaca
	}

	cfg := &Config{}
	if err := viper.Unmarshal(cfg); err != nil {
		return nil, err
	}

	// Validasi: JWT_SECRET wajib di-set
	if cfg.JWTSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET is required but not set. Please provide it via environment variable or .env file")
	}
	if len(cfg.JWTSecret) < 32 {
		return nil, fmt.Errorf("JWT_SECRET must be at least 32 characters long for security")
	}

	// Validasi: REDIS_PASSWORD wajib di-set
	if cfg.RedisPassword == "" {
		return nil, fmt.Errorf("REDIS_PASSWORD is required but not set. Please provide it via environment variable or .env file")
	}

	return cfg, nil
}
