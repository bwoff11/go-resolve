package cache

import (
	"fmt"
	"strings"

	"github.com/miekg/dns"
)

type Record struct {
	Domain string
	Type   uint16
	Value  []string
}

func FromRR(rr dns.RR) (*Record, error) {
	record := &Record{
		Domain: rr.Header().Name,
		Type:   rr.Header().Rrtype,
	}

	switch rr := rr.(type) {
	case *dns.A:
		record.Value = []string{rr.A.String()}
	case *dns.AAAA:
		record.Value = []string{rr.AAAA.String()}
	case *dns.CNAME:
		record.Value = []string{rr.Target}
	case *dns.TXT:
		record.Value = rr.Txt
	case *dns.MX:
		record.Value = []string{fmt.Sprintf("%v %s", rr.Preference, rr.Mx)}
	case *dns.NS:
		record.Value = []string{rr.Ns}
	case *dns.PTR:
		record.Value = []string{rr.Ptr}
	case *dns.SOA:
		record.Value = []string{fmt.Sprintf("%s %s %d %d %d %d %d", rr.Ns, rr.Mbox, rr.Serial, rr.Refresh, rr.Retry, rr.Expire, rr.Minttl)}
	// Add other cases as necessary.
	default:
		return nil, fmt.Errorf("unsupported DNS record type: %v", dns.TypeToString[rr.Header().Rrtype])
	}

	return record, nil
}

func (r *Record) ToRR() (dns.RR, error) {
	switch r.Type {
	case dns.TypeA:
		if len(r.Value) != 1 {
			return nil, fmt.Errorf("invalid value for A record")
		}
		return dns.NewRR(fmt.Sprintf("%s A %s", r.Domain, r.Value[0]))
	case dns.TypeAAAA:
		if len(r.Value) != 1 {
			return nil, fmt.Errorf("invalid value for AAAA record")
		}
		return dns.NewRR(fmt.Sprintf("%s AAAA %s", r.Domain, r.Value[0]))
	case dns.TypeCNAME:
		if len(r.Value) != 1 {
			return nil, fmt.Errorf("invalid value for CNAME record")
		}
		return dns.NewRR(fmt.Sprintf("%s CNAME %s", r.Domain, r.Value[0]))
	case dns.TypeTXT:
		// TXT records can have multiple strings, so we join them.
		txtValue := strings.Join(r.Value, " ")
		return dns.NewRR(fmt.Sprintf("%s TXT \"%s\"", r.Domain, txtValue))
	// Add other cases as necessary.
	default:
		return nil, fmt.Errorf("unsupported DNS record type")
	}
}
