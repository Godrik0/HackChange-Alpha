package storage

import (
	"fmt"
	"time"

	"github.com/Godrik0/HackChange-Alpha/backend/internal/config"
	"github.com/Godrik0/HackChange-Alpha/backend/internal/domain/interfaces"
	"github.com/Godrik0/HackChange-Alpha/backend/internal/domain/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

func NewPostgresRepository(cfg *config.Config, logger interfaces.Logger) (*gorm.DB, error) {
	logger.Info("Connecting to database", "host", cfg.Database.Host, "port", cfg.Database.Port)

	gormLogger := gormlogger.Default.LogMode(gormlogger.Info)

	db, err := gorm.Open(postgres.Open(cfg.Database.GetDSN()), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		logger.Error("Failed to connect to database", "error", err)
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		logger.Error("Failed to get database instance", "error", err)
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	logger.Info("Database connected successfully")
	return db, nil
}

func RunMigrations(db *gorm.DB, logger interfaces.Logger) error {
	logger.Info("Running database migrations")

	if err := db.AutoMigrate(&models.Client{}); err != nil {
		logger.Error("Failed to run migrations", "error", err)
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	logger.Info("Database migrations completed successfully")
	return nil
}
