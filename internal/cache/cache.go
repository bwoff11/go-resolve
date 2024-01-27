package cache

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

const dbName = "local_cache.db"

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

func init() {
	// Delete the existing database file if it exists.
	if err := os.Remove(dbName); err != nil && !os.IsNotExist(err) {
		log.Fatalf("Failed to remove existing database file: %v", err)
	}

	// Create a new SQLite database.
	db, err := sql.Open("sqlite3", dbName)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	// Create tables.
	createTables(db)
}

func createTables(db *sql.DB) {

	log.Println("Creating tables...")

	createDomainsSQL := `
    CREATE TABLE IF NOT EXISTS domains (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT NOT NULL,
        parent_id INTEGER,
        type TEXT NOT NULL,
        full_domain TEXT NOT NULL
    );`

	createRecordsSQL := `
    CREATE TABLE IF NOT EXISTS records (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        domain_id INTEGER,
        record_type TEXT NOT NULL,
        value TEXT NOT NULL,
        ttl INTEGER,
        FOREIGN KEY(domain_id) REFERENCES domains(id)
    );`

	_, err := db.Exec(createDomainsSQL)
	if err != nil {
		log.Fatalf("Failed to create domains table: %v", err)
	}

	_, err = db.Exec(createRecordsSQL)
	if err != nil {
		log.Fatalf("Failed to create records table: %v", err)
	}
}

// AddDomainRecord adds a new domain record to the database.
func AddDomainRecord(dr *DomainRecord) error {
	// Implement the logic to insert a record into the domains table.
	return nil
}

// RemoveDomainRecord removes a domain record from the database.
func RemoveDomainRecord(db *sql.DB, id int) error {
	// Implement the logic to remove a record from the domains table.
	return nil
}

// QueryDomainRecord queries for a domain record in the database.
func QueryDomainRecord(db *sql.DB, name string, recordType string) (*DomainRecord, error) {
	// Implement the logic to query a record from the domains table.
	return nil, nil
}

func AddLocalRecord(db *sql.DB, name string, recordType string, value string) error {
	// Implement the logic to insert a record into the domains table.
	return nil
}
