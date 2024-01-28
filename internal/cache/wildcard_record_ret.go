package cache

import (
	"github.com/bwoff11/go-resolve/internal/config"
	"github.com/miekg/dns"
)

type WildcardRecordSet struct {
}

func NewWildcardRecordSet(wr []config.WildcardRecord) *WildcardRecordSet {
	return &WildcardRecordSet{}
}

func (w *WildcardRecordSet) Query(domain string, recordType uint16) ([]dns.RR, bool) {
	return nil, false
}
