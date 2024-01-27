package models

import (
	"net"

	"github.com/miekg/dns"
	"gorm.io/gorm"
)

type Record struct {
	gorm.Model
	DomainID uint
	Type     uint16
	Value    string
	TTL      int
}

func (r *Record) ToDNSMsg(id int) *dns.Msg {
	msg := new(dns.Msg)
	msg.SetReply(&dns.Msg{})
	msg.Id = uint16(id)

	fqdn := dns.Fqdn(r.Value) // Ensure the domain name is fully qualified

	switch r.Type {
	case dns.TypeA:
		rr := &dns.A{
			Hdr: dns.RR_Header{
				Name:   fqdn,
				Rrtype: dns.TypeA,
				Class:  dns.ClassINET,
				Ttl:    uint32(r.TTL),
			},
			A: net.ParseIP(r.Value),
		}
		msg.Answer = append(msg.Answer, rr)

	case dns.TypeAAAA:
		rr := &dns.AAAA{
			Hdr: dns.RR_Header{
				Name:   fqdn,
				Rrtype: dns.TypeAAAA,
				Class:  dns.ClassINET,
				Ttl:    uint32(r.TTL),
			},
			AAAA: net.ParseIP(r.Value),
		}
		msg.Answer = append(msg.Answer, rr)

	case dns.TypeCNAME:
		rr := &dns.CNAME{
			Hdr: dns.RR_Header{
				Name:   fqdn,
				Rrtype: dns.TypeCNAME,
				Class:  dns.ClassINET,
				Ttl:    uint32(r.TTL),
			},
			Target: fqdn,
		}
		msg.Answer = append(msg.Answer, rr)
	}

	return msg
}
