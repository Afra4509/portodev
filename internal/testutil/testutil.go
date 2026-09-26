package testutil

import (
	"fmt"
	"time"

	"rest-api-ipk-mahasiswa-its/internal/model"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// SetupTestDB initializes an isolated in-memory SQLite database for unit & integration testing
func SetupTestDB() *gorm.DB {
	dsn := fmt.Sprintf("file:test_%d?mode=memory&cache=shared", time.Now().UnixNano())
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		panic("failed to connect to in-memory test database: " + err.Error())
	}

	sqlDB, err := db.DB()
	if err == nil {
		sqlDB.Exec("PRAGMA foreign_keys = ON;")
	}

	// AutoMigrate all tables
	err = db.AutoMigrate(
		&model.Student{},
		&model.Course{},
		&model.Grade{},
	)
	if err != nil {
		panic("failed to migrate test database: " + err.Error())
	}

	return db
}
