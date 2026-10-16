package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"satset-app/internal/config"
)

// Connect establishes a connection to MySQL using GORM.
func Connect(cfg *config.Config) (*gorm.DB, error) {
	// Configure GORM logger based on environment
	var gormLogger logger.Interface
	if cfg.IsDevelopment() {
		gormLogger = logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			logger.Config{
				SlowThreshold:             200 * time.Millisecond,
				LogLevel:                  logger.Info,
				IgnoreRecordNotFoundError: true,
				Colorful:                  true,
			},
		)
	} else {
		gormLogger = logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			logger.Config{
				SlowThreshold:             200 * time.Millisecond,
				LogLevel:                  logger.Warn,
				IgnoreRecordNotFoundError: true,
				Colorful:                  false,
			},
		)
	}

	var db *gorm.DB
	var lastErr error
	maxRetries := 10

	for attempt := 1; attempt <= maxRetries; attempt++ {
		db, lastErr = gorm.Open(mysql.Open(cfg.DB.DSN), &gorm.Config{
			Logger:                                   gormLogger,
			DisableForeignKeyConstraintWhenMigrating: false,
			TranslateError:                           true,
		})
		if lastErr == nil {
			sqlDB, pingErr := db.DB()
			if pingErr == nil {
				if err := sqlDB.Ping(); err == nil {
					// Configure connection pool
					sqlDB.SetMaxIdleConns(10)
					sqlDB.SetMaxOpenConns(100)
					sqlDB.SetConnMaxLifetime(time.Hour)
					return db, nil
				} else {
					lastErr = fmt.Errorf("failed to ping database: %w", err)
				}
			} else {
				lastErr = fmt.Errorf("failed to get database instance: %w", pingErr)
			}
		}

		if attempt < maxRetries {
			log.Printf("Menunggu database siap (percobaan %d/%d, error: %v)...", attempt, maxRetries, lastErr)
			time.Sleep(2 * time.Second)
		}
	}

	return nil, fmt.Errorf("failed to connect to database after %d attempts: %w", maxRetries, lastErr)
}
