package config

import (
	"fmt"
	"log"
	"time"

	"rest-api-ipk-mahasiswa-its/internal/model"

	"github.com/glebarez/sqlite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// InitDatabase initializes the GORM database connection
// Supports PostgreSQL as primary database, with automatic fallback or explicit SQLite mode
func InitDatabase(cfg *Config) (*gorm.DB, error) {
	var gormDB *gorm.DB
	var err error

	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}

	if cfg.DBDriver == "sqlite" {
		log.Printf("[INFO] Initializing SQLite database: %s\n", cfg.DBName)
		gormDB, err = gorm.Open(sqlite.Open(cfg.DBName), gormConfig)
		if err != nil {
			return nil, fmt.Errorf("failed to open sqlite database: %w", err)
		}
	} else {
		// PostgreSQL connection string
		dsn := fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
			cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort, cfg.DBSSLMode, cfg.DBTimeZone,
		)

		log.Printf("[INFO] Attempting connection to PostgreSQL at %s:%s (database: %s)...\n", cfg.DBHost, cfg.DBPort, cfg.DBName)
		gormDB, err = gorm.Open(postgres.Open(dsn), gormConfig)
		if err != nil {
			log.Printf("[WARNING] Unable to connect to PostgreSQL (%v).", err)
			log.Println("[INFO] Auto-fallback to local SQLite 'its_academic.db' so the application can run immediately without requiring a running PostgreSQL server.")
			gormDB, err = gorm.Open(sqlite.Open("its_academic.db"), gormConfig)
			if err != nil {
				return nil, fmt.Errorf("failed fallback to sqlite database: %w", err)
			}
		}
	}

	// Configure connection pool
	sqlDB, err := gormDB.DB()
	if err == nil {
		sqlDB.SetMaxIdleConns(10)
		sqlDB.SetMaxOpenConns(100)
		sqlDB.SetConnMaxLifetime(time.Hour)

		// If using SQLite, enable foreign keys
		if gormDB.Dialector.Name() == "sqlite" {
			sqlDB.Exec("PRAGMA foreign_keys = ON;")
		}
	}

	// Auto-migrate tables and foreign key constraints
	log.Println("[INFO] Running database migrations (AutoMigrate)...")
	if err := gormDB.AutoMigrate(
		&model.Student{},
		&model.Course{},
		&model.Grade{},
	); err != nil {
		return nil, fmt.Errorf("failed to run database auto-migration: %w", err)
	}

	log.Println("[INFO] Database connected and schema migrated successfully.")
	return gormDB, nil
}
