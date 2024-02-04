package local

import (
	"sync"
	"time"

	"github.com/bwoff11/go-resolve/internal/config"
	"github.com/miekg/dns"
	"github.com/rs/zerolog/log"
)

type LocalRecords struct {
	mutex   sync.RWMutex
	Records []Record
}

type Record struct {
	Question *dns.Question
	Answer   []dns.RR
	Expiry   time.Time
}

// New creates a new LocalRecords instance.
// It converts the configuration records to
// DNS messages and adds them to the local
// records cache.
func New(cfg *config.Local) *LocalRecords {
	l := &LocalRecords{}
	l.addRecords(cfg.Standard)
	return l
}

func (l *LocalRecords) addRecords(records []config.StandardRecord) {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	for _, r := range records {
		msg, err := r.ToMsg()
		if err != nil {
			log.Error().Err(err).Msg("failed to convert record to message")
			continue
		}

		l.Records = append(l.Records, Record{
			Question: &msg.Question[0],
			Answer:   msg.Answer,
			Expiry:   time.Now().Add(time.Duration(r.TTL) * time.Second),
		})

		log.Debug().Str("domain", r.Domain).Str("type", r.Type).Str("value", r.Value).Int("ttl", r.TTL).Msg("added local record")
	}
}

// Query returns the DNS message for the given question.
func (l *LocalRecords) Query(q *dns.Question) []dns.RR {
	l.mutex.RLock()
	defer l.mutex.RUnlock()

	for _, r := range l.Records {
		if r.Question.Name == q.Name && (r.Question.Qtype == q.Qtype || r.Question.Qtype == dns.TypeCNAME) {
			log.Debug().Str("domain", q.Name).Str("type", dns.TypeToString[q.Qtype]).Msg("local record found")
			return r.Answer
		}
	}

	log.Debug().Str("domain", q.Name).Str("type", dns.TypeToString[q.Qtype]).Msg("local record not found")
	return nil
}
