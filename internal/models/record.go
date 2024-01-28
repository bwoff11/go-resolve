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
	if r.ExpiresAt == nil {
		return false
	}
	return r.ExpiresAt.Before(time.Now())
}

func (r *Record) FromRR(rr dns.RR, domainID uint) Record {
	switch rr.Header().Rrtype {
	case dns.TypeA:
		return r.FromA(rr.(*dns.A), domainID)
	case dns.TypeAAAA:
		return r.FromAAAA(rr.(*dns.AAAA), domainID)
	case dns.TypeCNAME:
		return r.FromCNAME(rr.(*dns.CNAME), domainID)
	default:
		return Record{}
	}
}

func (r *Record) FromA(a *dns.A, domainID uint) Record {
	return Record{
		DomainID:  domainID,
		Type:      dns.TypeA,
		Value:     a.A.String(),
		TTL:       int(a.Hdr.Ttl),
		ExpiresAt: calculateExpiresAt(int(a.Hdr.Ttl)),
	}
}

func (r *Record) FromAAAA(aaaa *dns.AAAA, domainID uint) Record {
	return Record{
		DomainID:  domainID,
		Type:      dns.TypeAAAA,
		Value:     aaaa.AAAA.String(),
		TTL:       int(aaaa.Hdr.Ttl),
		ExpiresAt: calculateExpiresAt(int(aaaa.Hdr.Ttl)),
	}
}

func (r *Record) FromCNAME(cname *dns.CNAME, domainID uint) Record {
	return Record{
		DomainID:  domainID,
		Type:      dns.TypeCNAME,
		Value:     cname.Target,
		TTL:       int(cname.Hdr.Ttl),
		ExpiresAt: calculateExpiresAt(int(cname.Hdr.Ttl)),
	}
}

func calculateExpiresAt(ttl int) *time.Time {
	expiresAt := time.Now().Add(time.Duration(ttl) * time.Second)
	return &expiresAt
}

func (r *Record) ToRR(domainName string) dns.RR {
	switch r.Type {
	case dns.TypeA:
		return r.toA(domainName)
	case dns.TypeAAAA:
		return r.toAAAA(domainName)
	case dns.TypeCNAME:
		return r.toCNAME(domainName)
	default:
		return nil
	}
}

func (r *Record) toA(domainName string) *dns.A {
	return &dns.A{
		Hdr: dns.RR_Header{
			Name:   domainName,
			Rrtype: dns.TypeA,
			Class:  dns.ClassINET,
			Ttl:    uint32(r.TTL),
		},
		A: net.ParseIP(r.Value),
	}
}

func (r *Record) toAAAA(domainName string) *dns.AAAA {
	return &dns.AAAA{
		Hdr: dns.RR_Header{
			Name:   domainName,
			Rrtype: dns.TypeAAAA,
			Class:  dns.ClassINET,
			Ttl:    uint32(r.TTL),
		},
		AAAA: net.ParseIP(r.Value),
	}
}

func (r *Record) toCNAME(domainName string) *dns.CNAME {
	return &dns.CNAME{
		Hdr: dns.RR_Header{
			Name:   domainName,
			Rrtype: dns.TypeCNAME,
			Class:  dns.ClassINET,
			Ttl:    uint32(r.TTL),
		},
		Target: r.Value,
	}
}
