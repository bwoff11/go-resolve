package common

import (
	"fmt"
	"strings"

	"github.com/miekg/dns"
)

type LocalRecord struct {
	Domain string
	Type   string
	Value  []string
	TTL    uint32
}

func (r *LocalRecord) FromRR(rr dns.RR) error {
	r.Domain = rr.Header().Name
	r.Type = dns.TypeToString[rr.Header().Rrtype]
	r.TTL = rr.Header().Ttl

	switch rr := rr.(type) {
	case *dns.A:
		r.Value = []string{rr.A.String()}
	case *dns.AAAA:
		r.Value = []string{rr.AAAA.String()}
	case *dns.CNAME:
		r.Value = []string{rr.Target}
	case *dns.TXT:
		r.Value = rr.Txt
	case *dns.MX:
		r.Value = []string{fmt.Sprintf("%v %s", rr.Preference, rr.Mx)}
	case *dns.NS:
		r.Value = []string{rr.Ns}
	case *dns.PTR:
		r.Value = []string{rr.Ptr}
	case *dns.SOA:
		r.Value = []string{fmt.Sprintf("%s %s %d %d %d %d %d", rr.Ns, rr.Mbox, rr.Serial, rr.Refresh, rr.Retry, rr.Expire, rr.Minttl)}
	// Add other cases as necessary.
	default:
		return fmt.Errorf("unsupported DNS record type: %v", dns.TypeToString[rr.Header().Rrtype])
	}

	return nil
}

func (r *LocalRecord) ToRR() (dns.RR, error) {
	switch r.Type {
	case "A":
		if len(r.Value) != 1 {
			return nil, fmt.Errorf("invalid value for A record")
		}
		return dns.NewRR(fmt.Sprintf("%s A %s", r.Domain, r.Value[0]))
	case "AAAA":
		if len(r.Value) != 1 {
			return nil, fmt.Errorf("invalid value for AAAA record")
		}
		return dns.NewRR(fmt.Sprintf("%s AAAA %s", r.Domain, r.Value[0]))
	case "CNAME":
		if len(r.Value) != 1 {
			return nil, fmt.Errorf("invalid value for CNAME record")
		}
		return dns.NewRR(fmt.Sprintf("%s CNAME %s", r.Domain, r.Value[0]))
	case "TXT":
		// TXT records can have multiple strings, so we join them.
		txtValue := strings.Join(r.Value, " ")
		return dns.NewRR(fmt.Sprintf("%s TXT \"%s\"", r.Domain, txtValue))
	// Add other cases as necessary.
	default:
		return nil, fmt.Errorf("unsupported DNS record type")
	}
}
