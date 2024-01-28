package cache

import (
	"github.com/bwoff11/go-resolve/internal/config"
	"github.com/bwoff11/go-resolve/internal/models"
	"github.com/miekg/dns"
	"gorm.io/gorm"
)

const dbName = "local_cache.db"

type Cache struct {
	db *gorm.DB
}

// DomainRecord represents a DNS record in the database.
type DomainRecord struct {
	ID         int
	FullDomain string
}

// Record represents a DNS record in the database.
type Record struct {
	ID       int
	DomainID int
	Type     string
	Value    string
	TTL      int
}

func New(cacheConfig config.CacheConfig, localDNSConfig config.LocalDNSConfig) (*Cache, error) {

	db := createNewDatabase()

	return &Cache{
		db: db,
	}, nil
}

// Add accepts a slice of dns.RRs and adds them to the cache.
func (c *Cache) Add(rrs []dns.RR) error {
	for _, rr := range rrs {
		var domain models.Domain
		if err := c.db.Where("full_domain = ?", rr.Header().Name).First(&domain).Error; err != nil {
			return err
		}

		var record models.Record
		if err := c.db.Where("domain_id = ? AND type = ? AND value = ?", domain.ID, rr.Header().Rrtype, rr.String()).First(&record).Error; err != nil {
			return err
		}

		if err := c.db.Create(&record).Error; err != nil {
			return err
		}
	}
	return nil
}

// Query accepts a dns.Msg and responds with the RRs, whether the
// response is authoritative, and a possible error.
func (c *Cache) Query(question dns.Question) ([]dns.RR, bool, error) {
	var domain models.Domain
	if err := c.db.Where("full_domain = ?", question.Name).First(&domain).Error; err != nil {
		return nil, false, err
	}

	var records []models.Record
	if err := c.db.Where("domain_id = ?", domain.ID).Find(&records).Error; err != nil {
		return nil, false, err
	}

	var rrs []dns.RR
	for _, record := range records {
		rr := record.ToRR(domain.FullDomain)
		rrs = append(rrs, rr)
	}

	return rrs, true, nil
}
