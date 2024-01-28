package cache

import (
	"github.com/bwoff11/go-resolve/internal/config"
	"github.com/miekg/dns"
)

type Cache struct {
	LocalRecords    *LocalRecordSet
	WildcardRecords *WildcardRecordSet
	RemoteRecords   *RemoteRecordSet
}

func New(lr []config.StandardRecord, wr []config.WildcardRecord) *Cache {
	return &Cache{
		LocalRecords:    NewLocalRecordSet(lr),
		WildcardRecords: NewWildcardRecordSet(wr),
		RemoteRecords:   NewRemoteRecordSet(),
	}
}

// The only time records should be added outside of New() is when the upstream
// servers are queried. Therefore, this function only supports additions to the
// remote record set.
func (c *Cache) Add(records []dns.RR) {
	c.RemoteRecords.Add(records)
}

func (c *Cache) Query(domain string, recordType uint16) ([]dns.RR, bool) {

	// Check local records first.
	records, ok := c.LocalRecords.Query(domain, recordType)
	if ok {
		return records, true
	}

	// Check wildcard records next.
	records, ok = c.WildcardRecords.Query(domain, recordType)
	if ok {
		return records, true
	}

	// Check remote records last.
	records, ok = c.RemoteRecords.Query(domain, recordType)
	if ok {
		return records, true
	}

	// No records found.
	return nil, false
}
