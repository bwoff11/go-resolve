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
	Domain    string
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

func (r *Record) FromRR(rr dns.RR) Record {
	switch rr.Header().Rrtype {
	case dns.TypeA:
		return r.FromA(rr.(*dns.A))
	case dns.TypeAAAA:
		return r.FromAAAA(rr.(*dns.AAAA))
	case dns.TypeCNAME:
		return r.FromCNAME(rr.(*dns.CNAME))
	default:
		return Record{}
	}
}

func (r *Record) FromA(a *dns.A) Record {
	return Record{
		Domain:    a.Header().Name,
		Type:      dns.TypeA,
		Value:     a.A.String(),
		TTL:       int(a.Hdr.Ttl),
		ExpiresAt: calculateExpiresAt(int(a.Hdr.Ttl)),
	}
}

func (r *Record) FromAAAA(aaaa *dns.AAAA) Record {
	return Record{
		Domain:    aaaa.Header().Name,
		Type:      dns.TypeAAAA,
		Value:     aaaa.AAAA.String(),
		TTL:       int(aaaa.Hdr.Ttl),
		ExpiresAt: calculateExpiresAt(int(aaaa.Hdr.Ttl)),
	}
}

func (r *Record) FromCNAME(cname *dns.CNAME) Record {
	return Record{
		Domain:    cname.Header().Name,
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

func (r *Record) ToRR() dns.RR {
	switch r.Type {
	case dns.TypeA:
		return r.toA()
	case dns.TypeAAAA:
		return r.toAAAA()
	case dns.TypeCNAME:
		return r.toCNAME()
	default:
		return nil
	}
}

func (r *Record) toA() *dns.A {
	return &dns.A{
		Hdr: dns.RR_Header{
			Name:   r.Domain,
			Rrtype: dns.TypeA,
			Class:  dns.ClassINET,
			Ttl:    uint32(r.TTL),
		},
		A: net.ParseIP(r.Value),
	}
}

func (r *Record) toAAAA() *dns.AAAA {
	return &dns.AAAA{
		Hdr: dns.RR_Header{
			Name:   r.Domain,
			Rrtype: dns.TypeAAAA,
			Class:  dns.ClassINET,
			Ttl:    uint32(r.TTL),
		},
		AAAA: net.ParseIP(r.Value),
	}
}

func (r *Record) toCNAME() *dns.CNAME {
	return &dns.CNAME{
		Hdr: dns.RR_Header{
			Name:   r.Domain,
			Rrtype: dns.TypeCNAME,
			Class:  dns.ClassINET,
			Ttl:    uint32(r.TTL),
		},
		Target: r.Value,
	}
}
