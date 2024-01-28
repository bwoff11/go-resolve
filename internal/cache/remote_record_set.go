package cache

import "github.com/miekg/dns"

type RemoteRecordSet struct {
}

func NewRemoteRecordSet() *RemoteRecordSet {
	return &RemoteRecordSet{}
}

func (r *RemoteRecordSet) Query(domain string, recordType uint16) ([]dns.RR, bool) {
	return nil, false
}
