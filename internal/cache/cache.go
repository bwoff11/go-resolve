package cache

import (
	"log"

	"github.com/bwoff11/go-resolve/internal/config"
	"github.com/bwoff11/go-resolve/internal/models"
	"github.com/miekg/dns"
	"gorm.io/gorm"
)

const dbName = "local_cache.db"

type Cache struct {
	db *gorm.DB
}

func New(cacheConfig config.CacheConfig, localRecords []config.DNSRecord) (*Cache, error) {
	db := createNewDatabase()
	cache := &Cache{db: db}

	if err := cache.addLocalRecords(localRecords); err != nil {
		return nil, err
	}

	return cache, nil
}

func (c *Cache) addLocalRecords(cfg []config.DNSRecord) error {
	for _, localRecord := range cfg {
		var record models.Record
		record.FromConfig(localRecord)
		if err := c.db.Create(&record).Error; err != nil {
			return err
		}
	}
	return nil
}

func (c *Cache) Add(rrs []dns.RR) error {
	for _, rr := range rrs {
		record := new(models.Record)
		*record = record.FromRR(rr)
		if err := c.db.Create(record).Error; err != nil {
			return err
		}
	}
	return nil
}

func (c *Cache) Query(question dns.Question) ([]dns.RR, bool, error) {
	var records []models.Record
	if err := c.db.Where("domain = ? AND type = ?", question.Name, question.Qtype).Find(&records).Error; err != nil {
		log.Printf("Failed to query cache: %v", err)
		return nil, false, err
	}

	// Check if no records were found
	if len(records) == 0 {
		return nil, false, nil // No records found
	}

	rrs := make([]dns.RR, 0, len(records))
	for _, record := range records {
		if !record.IsExpired() {
			rrs = append(rrs, record.ToRR())
		}
	}

	return rrs, true, nil // Records found and returned
}

// Additional helper functions as required
