package cache

import (
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/miekg/dns"
)

func NormalizeDomain(domain string) string {
	return strings.TrimSuffix(domain, ".") + "."
}

func extractRRValue(record dns.RR) string {
	switch rr := record.(type) {
	case *dns.A:
		return rr.A.String()
	case *dns.AAAA:
		return rr.AAAA.String()
	case *dns.CNAME:
		return rr.Target
	case *dns.MX:
		return fmt.Sprintf("%v %s", rr.Preference, rr.Mx)
	case *dns.NS:
		return rr.Ns
	case *dns.PTR:
		return rr.Ptr
	case *dns.SOA:
		return fmt.Sprintf("%s %s %d %d %d %d %d", rr.Ns, rr.Mbox, rr.Serial, rr.Refresh, rr.Retry, rr.Expire, rr.Minttl)
	case *dns.TXT:
		return strings.Join(rr.Txt, " ")
	default:
		return ""
	}
}

func createRRFromValue(domain, recordType, value string, ttl int) (dns.RR, error) {
	header := dns.RR_Header{
		Name:   dns.Fqdn(domain),
		Rrtype: dns.StringToType[recordType],
		Class:  dns.ClassINET,
		Ttl:    uint32(ttl),
	}

	switch dns.StringToType[recordType] {
	case dns.TypeA:
		ip := net.ParseIP(value)
		if ip == nil {
			return nil, fmt.Errorf("invalid IP address")
		}
		return &dns.A{Hdr: header, A: ip}, nil
	case dns.TypeAAAA:
		ip := net.ParseIP(value)
		if ip == nil {
			return nil, fmt.Errorf("invalid IP address")
		}
		return &dns.AAAA{Hdr: header, AAAA: ip}, nil
	case dns.TypeCNAME:
		return &dns.CNAME{Hdr: header, Target: dns.Fqdn(value)}, nil
	case dns.TypeMX:
		parts := strings.Fields(value)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid MX record format")
		}
		pref, err := strconv.Atoi(parts[0])
		if err != nil {
			return nil, fmt.Errorf("invalid MX preference")
		}
		return &dns.MX{Hdr: header, Preference: uint16(pref), Mx: dns.Fqdn(parts[1])}, nil
	case dns.TypeNS:
		return &dns.NS{Hdr: header, Ns: dns.Fqdn(value)}, nil
	case dns.TypePTR:
		return &dns.PTR{Hdr: header, Ptr: dns.Fqdn(value)}, nil
	case dns.TypeSOA:
		// Placeholders for now.
	case dns.TypeTXT:
		return &dns.TXT{Hdr: header, Txt: strings.Fields(value)}, nil
	// Add other cases as needed.
	default:
		return nil, fmt.Errorf("unsupported record type: %s", recordType)
	}
	return nil, nil
}
