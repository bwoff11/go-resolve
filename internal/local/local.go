package local

import (
	"github.com/bwoff11/go-resolve/internal/config"
	"github.com/miekg/dns"
)

type LocalRecords struct {
}

func New(cfg config.Local) *LocalRecords {
	return &LocalRecords{}
}

func (l *LocalRecords) Query(q dns.Question) []dns.RR {
	return nil
}
