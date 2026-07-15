package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	AppPort          string `mapstructure:"APP_PORT"`
	DBHost           string `mapstructure:"MYSQL_HOST"`
	DBPort           string `mapstructure:"MYSQL_PORT"`
	DBUser           string `mapstructure:"MYSQL_USER"`
	DBPass           string `mapstructure:"MYSQL_PASSWORD"`
	DBName           string `mapstructure:"MYSQL_DATABASE"`
	RedisHost        string `mapstructure:"REDIS_HOST"`
	RedisPort        string `mapstructure:"REDIS_PORT"`
	RedisPassword    string `mapstructure:"REDIS_PASSWORD"`
	CORSAllowOrigins string `mapstructure:"CORS_ALLOW_ORIGINS"` // Tambahan keamanan CORS
}

func LoadConfig() (*Config, error) {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	// Default values (hanya yang tidak sensitif)
	viper.SetDefault("APP_PORT", "8084")
	viper.SetDefault("MYSQL_HOST", "localhost")
	viper.SetDefault("MYSQL_PORT", "3306")
	viper.SetDefault("MYSQL_USER", "root")
	viper.SetDefault("MYSQL_DATABASE", "cart_db")
	viper.SetDefault("REDIS_HOST", "localhost")
	viper.SetDefault("REDIS_PORT", "6379")
	viper.SetDefault("CORS_ALLOW_ORIGINS", "http://localhost:8080")

	// Secrets tidak ada defaultnya (wajib diisi environment variable):
	// Tapi harus di-BindEnv agar Viper bisa baca dari OS env var saat Unmarshal
	viper.BindEnv("MYSQL_PASSWORD")
	viper.BindEnv("REDIS_PASSWORD")

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
	// REDIS_PASSWORD disarankan diisi di production, 
	// namun untuk pengembangan lokal sering kali redis tidak dipasword.
	// Jika mau diwajibkan:
	// if cfg.RedisPassword == "" {
	// 	return nil, fmt.Errorf("REDIS_PASSWORD is required but not set")
	// }

	return cfg, nil
}
