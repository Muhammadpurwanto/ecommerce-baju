package database

import (
	"fmt"
	"time"

	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/Muhammadpurwanto/ecommerce-baju/services/order/config"
	"github.com/Muhammadpurwanto/ecommerce-baju/services/order/internal/model"
)

func NewMySQL(cfg *config.Config, log *zap.Logger) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.DBUser,
		cfg.DBPass,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn), // Gunakan Warn agar log query tidak bocor
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql.DB: %w", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// UNTUK DEVELOPMENT: AutoMigrate dinyalakan agar tabel otomatis dibuat
	if err := db.AutoMigrate(&model.Order{}, &model.OrderItem{}); err != nil {
		return nil, fmt.Errorf("failed to auto migrate: %w", err)
	}

	log.Info("Connected to MySQL database", zap.String("database", cfg.DBName))
	return db, nil
}
