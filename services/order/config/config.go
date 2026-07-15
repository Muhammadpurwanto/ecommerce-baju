package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	AppPort           string `mapstructure:"APP_PORT"`
	DBHost            string `mapstructure:"MYSQL_HOST"`
	DBPort            string `mapstructure:"MYSQL_PORT"`
	DBUser            string `mapstructure:"MYSQL_USER"`
	DBPass            string `mapstructure:"MYSQL_PASSWORD"`
	DBName            string `mapstructure:"MYSQL_DATABASE"`
	RabbitMQURL       string `mapstructure:"RABBITMQ_URL"`
	CORSAllowOrigins  string `mapstructure:"CORS_ALLOW_ORIGINS"` // Tambahan keamanan CORS
	ProductServiceURL string `mapstructure:"PRODUCT_SERVICE_URL"`
	CartServiceURL    string `mapstructure:"CART_SERVICE_URL"`
	PaymentServiceURL string `mapstructure:"PAYMENT_SERVICE_URL"`
}

func LoadConfig() (*Config, error) {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	// Default values (hanya yang tidak sensitif)
	viper.SetDefault("APP_PORT", "8085")
	viper.SetDefault("MYSQL_HOST", "localhost")
	viper.SetDefault("MYSQL_PORT", "3306")
	viper.SetDefault("MYSQL_USER", "root")
	viper.SetDefault("MYSQL_DATABASE", "order_db")
	viper.SetDefault("CORS_ALLOW_ORIGINS", "http://localhost:8080")
	viper.SetDefault("PRODUCT_SERVICE_URL", "localhost:50053")
	viper.SetDefault("CART_SERVICE_URL", "localhost:50054")
	viper.SetDefault("PAYMENT_SERVICE_URL", "localhost:50056")

	// Secrets tidak ada defaultnya (wajib diisi environment variable):
	// Tapi harus di-BindEnv agar Viper bisa baca dari OS env var saat Unmarshal
	viper.BindEnv("MYSQL_PASSWORD")
	viper.BindEnv("RABBITMQ_URL")

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

	return cfg, nil
}
