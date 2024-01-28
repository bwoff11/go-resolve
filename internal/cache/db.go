package cache

import (
	"log"
	"os"

	"github.com/bwoff11/go-resolve/internal/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// removeExistingDatabase checks and removes the existing database file if it exists.
func removeExistingDatabase() {
	if _, err := os.Stat(dbName); err == nil || !os.IsNotExist(err) {
		if err := os.Remove(dbName); err != nil {
			log.Fatalf("Failed to remove existing database file: %v", err)
		}
	}
}

// createNewDatabase initializes and returns a new database connection.
func createNewDatabase() *gorm.DB {
	removeExistingDatabase()

	db, err := gorm.Open(sqlite.Open(dbName), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to create database: %v", err)
	}

	if err := db.AutoMigrate(&models.Record{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	go housekeeping(db, 10)

	return db
}
