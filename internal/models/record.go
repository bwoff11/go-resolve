package models

import (
	"net"
	"time"

	"github.com/miekg/dns"
)

type Record struct {
	ID        uint
	CreatedAt time.Time
	UpdatedAt time.Time
	ExpiresAt *time.Time
	DomainID  uint
	Type      uint16
	Value     string
	TTL       int
}

func (r *Record) IsExpired() bool {
	return r.ExpiresAt.Before(time.Now())
}

func (r *Record) ToRR(domainName string) dns.RR {
	switch r.Type {
	case dns.TypeA:
		return &dns.A{
			Hdr: dns.RR_Header{
				Name:   domainName,
				Rrtype: dns.TypeA,
				Class:  dns.ClassINET,
				Ttl:    uint32(r.TTL),
			},
			A: net.ParseIP(r.Value),
		}
	case dns.TypeAAAA:
		return &dns.AAAA{
			Hdr: dns.RR_Header{
				Name:   domainName,
				Rrtype: dns.TypeAAAA,
				Class:  dns.ClassINET,
				Ttl:    uint32(r.TTL),
			},
			AAAA: net.ParseIP(r.Value),
		}
	case dns.TypeCNAME:
		return &dns.CNAME{
			Hdr: dns.RR_Header{
				Name:   domainName,
				Rrtype: dns.TypeCNAME,
				Class:  dns.ClassINET,
				Ttl:    uint32(r.TTL),
			},
			Target: r.Value,
		}
	default:
		return nil
	}
}
