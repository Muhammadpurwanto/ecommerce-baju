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
	MinioEndpoint    string `mapstructure:"MINIO_ENDPOINT"`
	MinioAccessKey   string `mapstructure:"MINIO_ACCESS_KEY"`
	MinioSecretKey   string `mapstructure:"MINIO_SECRET_KEY"`
	MinioBucket      string `mapstructure:"MINIO_BUCKET"`
	MinioPublicURL   string `mapstructure:"MINIO_PUBLIC_URL"` // Tambahan untuk production
	CORSAllowOrigins string `mapstructure:"CORS_ALLOW_ORIGINS"` // Tambahan keamanan CORS
}

func LoadConfig() (*Config, error) {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	// Hanya nilai yang non-sensitif yang boleh punya default
	viper.SetDefault("APP_PORT", "8083")
	viper.SetDefault("MYSQL_HOST", "localhost")
	viper.SetDefault("MYSQL_PORT", "3306")
	viper.SetDefault("MYSQL_USER", "root")
	viper.SetDefault("MYSQL_DATABASE", "product_db")
	viper.SetDefault("REDIS_HOST", "localhost")
	viper.SetDefault("REDIS_PORT", "6379")
	viper.SetDefault("MINIO_ENDPOINT", "http://localhost:9000")
	viper.SetDefault("MINIO_BUCKET", "products")
	viper.SetDefault("MINIO_PUBLIC_URL", "http://localhost:9000")
	viper.SetDefault("CORS_ALLOW_ORIGINS", "http://localhost:8080")

	// Secrets tidak ada defaultnya (wajib diisi environment variable):
	// Tapi harus di-BindEnv agar Viper bisa baca dari OS env var saat Unmarshal
	viper.BindEnv("MYSQL_PASSWORD")
	viper.BindEnv("REDIS_PASSWORD")
	viper.BindEnv("MINIO_ACCESS_KEY")
	viper.BindEnv("MINIO_SECRET_KEY")

	if err := viper.ReadInConfig(); err != nil {
		// Abaikan jika .env tidak ada (ambil dari sys env var)
	}

	cfg := &Config{}
	if err := viper.Unmarshal(cfg); err != nil {
		return nil, err
	}

	// Validasi Wajib
	if cfg.DBPass == "" {
		return nil, fmt.Errorf("MYSQL_PASSWORD is required but not set")
	}
	if cfg.MinioAccessKey == "" {
		return nil, fmt.Errorf("MINIO_ACCESS_KEY is required but not set")
	}
	if cfg.MinioSecretKey == "" {
		return nil, fmt.Errorf("MINIO_SECRET_KEY is required but not set")
	}

	return cfg, nil
}
