package config

import (
	"net"

	"github.com/miekg/dns"
	"github.com/rs/zerolog/log"
)

type Cache struct {
	Enabled       bool          `yaml:"enabled"`
	PruneInterval int           `yaml:"pruneInterval"`
	LocalRecords  []LocalRecord `yaml:"localRecords"`
}

type LocalRecord struct {
	Domain string `yaml:"domain"`
	Type   string `yaml:"type"`
	Value  string `yaml:"value"`
	TTL    int    `yaml:"ttl"`
}

func (lr *LocalRecord) ToQuestion() dns.Question {
	return dns.Question{
		Name:   lr.Domain + ".",
		Qtype:  dns.StringToType[lr.Type],
		Qclass: dns.ClassINET,
	}
}

func (lr *LocalRecord) ToAnswer() []dns.RR {
	var rr []dns.RR
	switch lr.Type {
	case "A":
		rr = append(rr, &dns.A{
			Hdr: dns.RR_Header{
				Name:   lr.Domain + ".",
				Rrtype: dns.TypeA,
				Class:  dns.ClassINET,
				Ttl:    uint32(lr.TTL),
			},
			A: net.ParseIP(lr.Value),
		})
	case "AAAA":
		rr = append(rr, &dns.AAAA{
			Hdr: dns.RR_Header{
				Name:   lr.Domain + ".",
				Rrtype: dns.TypeAAAA,
				Class:  dns.ClassINET,
				Ttl:    uint32(lr.TTL),
			},
			AAAA: net.ParseIP(lr.Value),
		})
	case "CNAME":
		rr = append(rr, &dns.CNAME{
			Hdr: dns.RR_Header{
				Name:   lr.Domain + ".",
				Rrtype: dns.TypeCNAME,
				Class:  dns.ClassINET,
				Ttl:    uint32(lr.TTL),
			},
			Target: lr.Value + ".",
		})
	default:
		log.Error().Str("type", lr.Type).Msg("Unsupported record type")
	}

	if len(rr) > 0 {
		return rr
	}
	return nil
}
