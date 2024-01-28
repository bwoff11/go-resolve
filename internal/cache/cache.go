package cache

import (
	"errors"
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
// If the record already exists, it is updated.
func (c *Cache) Add(rrs []dns.RR) error {
	log.Printf("Adding %d records to cache", len(rrs))
	for _, rr := range rrs {
		// Ensure the domain exists or create it
		domain := models.Domain{FullDomain: rr.Header().Name}
		if err := c.db.FirstOrCreate(&domain, models.Domain{FullDomain: rr.Header().Name}).Error; err != nil {
			return err
		}

		var record models.Record
		switch rr.Header().Rrtype {
		case dns.TypeA:
			record = record.FromA(rr.(*dns.A), domain.ID)
		case dns.TypeAAAA:
			record = record.FromAAAA(rr.(*dns.AAAA), domain.ID)
		case dns.TypeCNAME:
			record = record.FromCNAME(rr.(*dns.CNAME), domain.ID)
		default:
			return errors.New("unsupported record type")
		}

		// Update or create the record
		if err := c.db.Where(models.Record{DomainID: domain.ID, Type: record.Type, Value: record.Value}).
			Assign(record).
			FirstOrCreate(&models.Record{}).Error; err != nil {
			return err
		}
	}
	return nil
}

// Query accepts a dns.Msg and responds with the RRs, whether the
// response is authoritative, and a possible error.
func (c *Cache) Query(question dns.Question) ([]dns.RR, bool, error) {

	// Get the domain record
	var domain models.Domain
	if err := c.db.Where("full_domain = ?", question.Name).First(&domain).Error; err != nil {
		return nil, false, err
	}

	// Get the dns records
	var records []models.Record
	if err := c.db.Where("domain_id = ?", domain.ID).Find(&records).Error; err != nil {
		return nil, false, err
	}

	// Convert the records to dns.RRs
	var rrs []dns.RR
	for _, record := range records {
		if !record.IsExpired() {
			rr := record.ToRR(domain.FullDomain)
			rrs = append(rrs, rr)
		}
	}

	// If no records were found, return an error
	if len(rrs) == 0 {
		return nil, false, gorm.ErrRecordNotFound
	}

	// Return the records
	return rrs, true, nil
}
