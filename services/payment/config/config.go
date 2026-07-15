package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	AppPort         string `mapstructure:"APP_PORT"`
	DBHost          string `mapstructure:"MYSQL_HOST"`
	DBPort          string `mapstructure:"MYSQL_PORT"`
	DBUser          string `mapstructure:"MYSQL_USER"`
	DBPass          string `mapstructure:"MYSQL_PASSWORD"`
	DBName          string `mapstructure:"MYSQL_DATABASE"`
	MidtransServer  string `mapstructure:"MIDTRANS_SERVER_KEY"`
	MidtransClient  string `mapstructure:"MIDTRANS_CLIENT_KEY"`
	MidtransIsProd  bool   `mapstructure:"MIDTRANS_IS_PROD"`
	RabbitMQURL     string `mapstructure:"RABBITMQ_URL"`
	CORSAllowOrigins string `mapstructure:"CORS_ALLOW_ORIGINS"` // Tambahan keamanan CORS
	FrontendURL     string `mapstructure:"FRONTEND_URL"`       // URL Redirect Frontend
}

func LoadConfig() (*Config, error) {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	// Default values (hanya yang tidak sensitif)
	viper.SetDefault("APP_PORT", "8086")
	viper.SetDefault("MYSQL_HOST", "localhost")
	viper.SetDefault("MYSQL_PORT", "3306")
	viper.SetDefault("MYSQL_USER", "root")
	viper.SetDefault("MYSQL_DATABASE", "payment_db")
	viper.SetDefault("MIDTRANS_IS_PROD", false)
	viper.SetDefault("CORS_ALLOW_ORIGINS", "http://localhost:8080")
	viper.SetDefault("FRONTEND_URL", "http://localhost:3000")

	// Secrets tidak ada defaultnya (wajib diisi environment variable):
	// Tapi harus di-BindEnv agar Viper bisa baca dari OS env var saat Unmarshal
	viper.BindEnv("MYSQL_PASSWORD")
	viper.BindEnv("MIDTRANS_SERVER_KEY")
	viper.BindEnv("MIDTRANS_CLIENT_KEY")
	viper.BindEnv("RABBITMQ_URL")
	viper.BindEnv("FRONTEND_URL")

	if err := viper.ReadInConfig(); err != nil {
		// abaikan jika tidak ada .env
	}

	cfg := &Config{}
	if err := viper.Unmarshal(cfg); err != nil {
		return nil, err
	}

	// Validasi Wajib
	if cfg.DBPass == "" {
		return nil, fmt.Errorf("MYSQL_PASSWORD is required but not set")
	}
	if cfg.RabbitMQURL == "" {
		return nil, fmt.Errorf("RABBITMQ_URL is required but not set")
	}
	if cfg.MidtransServer == "" {
		return nil, fmt.Errorf("MIDTRANS_SERVER_KEY is required but not set")
	}
	if cfg.MidtransClient == "" {
		return nil, fmt.Errorf("MIDTRANS_CLIENT_KEY is required but not set")
	}

	return cfg, nil
}
