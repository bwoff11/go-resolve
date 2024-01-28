package cache

import "github.com/miekg/dns"

func getRRValue(record dns.RR) string {
	switch rr := record.(type) {
	case *dns.A:
		return rr.A.String()
	case *dns.AAAA:
		return rr.AAAA.String()
	case *dns.CNAME:
		return rr.Target
	case *dns.MX:
		return rr.Mx
	case *dns.NS:
		return rr.Ns
	case *dns.PTR:
		return rr.Ptr
	case *dns.SOA:
		return rr.Ns
	case *dns.TXT:
		return rr.Txt[0]
	default:
		return ""
	}
}
