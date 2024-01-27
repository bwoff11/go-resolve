package cache

import (
	"log"
	"os"

	"github.com/bwoff11/go-resolve/internal/config"
	"github.com/bwoff11/go-resolve/internal/models"
	"github.com/miekg/dns"
	"gorm.io/driver/sqlite"
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
	// Check if the database file exists
	if _, err := os.Stat(dbName); err == nil || !os.IsNotExist(err) {
		// Delete the file if it exists
		if err := os.Remove(dbName); err != nil {
			log.Fatalf("Failed to remove existing database file: %v", err)
		}
	}

	// Open a new database connection
	db, err := gorm.Open(sqlite.Open(dbName), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}

	// AutoMigrate tables
	err = db.AutoMigrate(&models.Domain{}, &models.Record{})
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	addRecords := func(domainName string, recordType uint16, value string) {
		var domain models.Domain
		if err := db.Where("full_domain = ?", domainName).FirstOrCreate(&domain, models.Domain{FullDomain: domainName}).Error; err != nil {
			log.Fatalf("Failed to add/find domain: %v", err)
		}

		record := models.Record{
			DomainID: domain.ID,
			Type:     recordType,
			Value:    value,
			TTL:      300,
		}

		if dbResult := db.Create(&record); dbResult.Error != nil {
			log.Fatalf("Failed to add local DNS record: %v", dbResult.Error)
		}
	}

	// Add local DNS A, AAAA, CNAME records
	for _, record := range localDNSConfig.Records.A {
		addRecords(record.Domain+".", dns.TypeA, record.IP)
	}
	for _, record := range localDNSConfig.Records.AAAA {
		addRecords(record.Domain+".", dns.TypeAAAA, record.IP)
	}
	for _, record := range localDNSConfig.Records.CNAME {
		addRecords(record.Domain, dns.TypeCNAME, record.Target)
	}

	return &Cache{db: db}, nil
}

// AddDomainRecord adds a new domain record to the database.
func AddDomainRecord(db *gorm.DB, domain *models.Domain) error {
	result := db.Create(domain)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// RemoveDomainRecord removes a domain record from the database.
func RemoveDomainRecord() error {
	return nil
}

// QueryDomainRecord queries for a domain record in the database.
func (c *Cache) Query(domain string, recordType uint16) (*models.Record, error) {
	var record models.Record
	var domainObj models.Domain

	if err := c.db.Where("full_domain = ?", domain).First(&domainObj).Error; err != nil {
		return nil, err
	}

	if err := c.db.Where("domain_id = ? AND type = ?", domainObj.ID, recordType).First(&record).Error; err != nil {
		return nil, err
	}

	return &record, nil
}

func (c *Cache) Add(msg *dns.Msg) error {
	for _, answer := range msg.Answer {
		var recordType uint16
		var value string

		switch rr := answer.(type) {
		case *dns.A:
			recordType = dns.TypeA
			value = rr.A.String()
		case *dns.AAAA:
			recordType = dns.TypeAAAA
			value = rr.AAAA.String()
		case *dns.CNAME:
			recordType = dns.TypeCNAME
			value = rr.Target
		default:
			// Handle other record types if necessary
			continue
		}

		domainName := dns.Fqdn(answer.Header().Name)
		ttl := int(answer.Header().Ttl)

		// Add record to the cache
		if err := c.addRecord(domainName, recordType, value, ttl); err != nil {
			return err
		}
	}

	return nil
}

func (c *Cache) addRecord(domainName string, recordType uint16, value string, ttl int) error {
	var domain models.Domain
	// Find or create the domain
	if err := c.db.Where("full_domain = ?", domainName).FirstOrCreate(&domain, models.Domain{FullDomain: domainName}).Error; err != nil {
		return err
	}

	// Create the record
	record := models.Record{
		DomainID: domain.ID,
		Type:     recordType,
		Value:    value,
		TTL:      ttl,
	}

	// Add the record to the database
	if dbResult := c.db.Create(&record); dbResult.Error != nil {
		return dbResult.Error
	}

	return nil
}
