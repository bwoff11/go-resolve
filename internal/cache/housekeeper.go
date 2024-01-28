package cache

import (
	"log"
	"time"

	"github.com/bwoff11/go-resolve/internal/models"
	"gorm.io/gorm"
)

func housekeeping(db *gorm.DB, interval int) {
	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	defer ticker.Stop()

	for {
		<-ticker.C
		log.Println("Running housekeeping...")
		deleteExpiredRecords(db)
	}
}

func deleteExpiredRecords(db *gorm.DB) {
	now := time.Now()
	result := db.Where("expires_at <= ?", now).Delete(&models.Record{})
	if result.Error != nil {
		log.Printf("Failed to delete expired records: %v\n", result.Error)
	} else {
		log.Printf("Deleted %d expired records\n", result.RowsAffected)
	}
}
