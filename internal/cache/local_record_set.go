package cache

import (
	"github.com/bwoff11/go-resolve/internal/config"
	"github.com/miekg/dns"
	"github.com/rs/zerolog/log"
)

type LocalRecordSet struct {
	Records []dns.RR
}

func NewLocalRecordSet(lr []config.StandardRecord) *LocalRecordSet {

	log.Debug().
		Str("msg", "Creating new local reset set").
		Int("record_count", len(lr)).
		Send()

	var lrs LocalRecordSet
	for _, record := range lr {

		// Convert the local record to a DNS RR.
		rr, err := createRRFromValue(record.Domain, record.Type, record.Value, record.TTL)
		if err != nil {
			log.Error().
				Err(err).
				Msg("Failed to convert local record to RR")
		} else {
			log.Debug().
				Str("msg", "Adding local record").
				Str("domain", rr.Header().Name).
				Str("type", dns.TypeToString[rr.Header().Rrtype]).
				//Str("value", getRRValue(rr)).
				Int("ttl", int(rr.Header().Ttl)).
				Send()
			lrs.Records = append(lrs.Records, rr)
		}
	}

	return &lrs
}

func (l *LocalRecordSet) Query(domain string, recordType uint16) ([]dns.RR, bool) {
	var records []dns.RR
	var ok bool

	for _, record := range l.Records {
		if record.Header().Name == domain && record.Header().Rrtype == recordType {
			records = append(records, record)
			ok = true
		}
	}

	if len(records) > 0 {
		for _, record := range records {
			log.Debug().
				Str("msg", "Found record in local record set").
				Str("domain", record.Header().Name).
				Str("type", dns.TypeToString[record.Header().Rrtype]).
				Str("value", extractRRValue(record)).
				Int("ttl", int(record.Header().Ttl)).
				Send()
		}
	} else {
		log.Debug().
			Str("msg", "No record found in local record set").
			Str("domain", domain).
			Str("type", dns.TypeToString[recordType]).
			Send()
	}

	return records, ok
}
