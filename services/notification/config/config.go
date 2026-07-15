package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	AppPort         string `mapstructure:"APP_PORT"`
	SMTPHost        string `mapstructure:"SMTP_HOST"`
	SMTPPort        int    `mapstructure:"SMTP_PORT"`
	SMTPUser        string `mapstructure:"SMTP_USER"`
	SMTPPass        string `mapstructure:"SMTP_PASSWORD"`
	SenderName      string `mapstructure:"SENDER_NAME"`
	RabbitMQURL     string `mapstructure:"RABBITMQ_URL"`
	CORSAllowOrigins string `mapstructure:"CORS_ALLOW_ORIGINS"` // Tambahan keamanan CORS
}

func LoadConfig() (*Config, error) {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	// Default values (hanya yang tidak sensitif)
	viper.SetDefault("APP_PORT", "8087")
	viper.SetDefault("SMTP_HOST", "sandbox.smtp.mailtrap.io")
	viper.SetDefault("SMTP_PORT", 2525)
	viper.SetDefault("SENDER_NAME", "Ecommerce Baju")
	viper.SetDefault("CORS_ALLOW_ORIGINS", "http://localhost:8080")

	// Secrets tidak ada defaultnya (wajib diisi environment variable):
	// Tapi harus di-BindEnv agar Viper bisa baca dari OS env var saat Unmarshal
	viper.BindEnv("SMTP_USER")
	viper.BindEnv("SMTP_PASSWORD")
	viper.BindEnv("RABBITMQ_URL")

	if err := viper.ReadInConfig(); err != nil {
		// abaikan jika tidak ada .env
	}

	cfg := &Config{}
	if err := viper.Unmarshal(cfg); err != nil {
		return nil, err
	}

	// Validasi Wajib
	if cfg.RabbitMQURL == "" {
		return nil, fmt.Errorf("RABBITMQ_URL is required but not set")
	}
	if cfg.SMTPUser == "" {
		return nil, fmt.Errorf("SMTP_USER is required but not set")
	}
	if cfg.SMTPPass == "" {
		return nil, fmt.Errorf("SMTP_PASSWORD is required but not set")
	}

	return cfg, nil
}
