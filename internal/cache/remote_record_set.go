package cache

import "github.com/miekg/dns"

type RemoteRecordSet struct {
	Records []dns.RR
}

func NewRemoteRecordSet() *RemoteRecordSet {
	return &RemoteRecordSet{}
}

func (r *RemoteRecordSet) Add(records []dns.RR) {
	r.Records = append(r.Records, records...)
}

func (r *RemoteRecordSet) Query(domain string, recordType uint16) ([]dns.RR, bool) {
	return nil, false
}
