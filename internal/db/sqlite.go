package db

import (
	"fmt"
	"log"

	"github.com/glebarez/sqlite"
	"github.com/wbw1537/synapse/internal/models"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Database struct {
	Conn *gorm.DB
}

func Connect(dbPath string) (*Database, error) {
	// Open the database using the pure-Go sqlite driver compatible with GORM
	// We use Silent logger to avoid noise, but you can change this to Info for debugging
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Optimize SQLite settings (Direct SQL execution via GORM)
	// WAL mode for better concurrency
	if err := db.Exec("PRAGMA journal_mode=WAL;").Error; err != nil {
		return nil, fmt.Errorf("failed to set WAL mode: %w", err)
	}
	// Normal synchronous mode
	if err := db.Exec("PRAGMA synchronous=NORMAL;").Error; err != nil {
		return nil, fmt.Errorf("failed to set synchronous mode: %w", err)
	}

	log.Printf("Connected to SQLite database at %s (GORM)", dbPath)

	return &Database{Conn: db}, nil
}

func (d *Database) InitSchema() error {
	// AutoMigrate creates tables, missing columns, and indexes automatically
	err := d.Conn.AutoMigrate(&models.Service{})
	if err != nil {
		return fmt.Errorf("failed to auto-migrate schema: %w", err)
	}
	return nil
}

func (d *Database) Close() error {
	sqlDB, err := d.Conn.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}