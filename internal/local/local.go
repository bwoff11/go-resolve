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

	var finalAnswers []dns.RR

	for _, r := range l.Records {
		if r.Question.Name == q.Name && r.Question.Qtype == q.Qtype {
			log.Debug().Str("domain", q.Name).Str("type", dns.TypeToString[q.Qtype]).Msg("local record found")
			finalAnswers = append(finalAnswers, r.Answer...)
		} else if r.Question.Name == q.Name && r.Question.Qtype == dns.TypeCNAME {
			// Found a CNAME for the queried name. Need to do a recursive lookup for the CNAME target.
			log.Debug().Str("domain", q.Name).Str("type", dns.TypeToString[q.Qtype]).Msg("CNAME record found, performing recursive lookup")
			finalAnswers = append(finalAnswers, r.Answer...)
			cnameRecord := r.Answer[0].(*dns.CNAME)
			cnameQuestion := &dns.Question{Name: cnameRecord.Target, Qtype: q.Qtype, Qclass: q.Qclass}
			additionalAnswers := l.recursiveLookup(cnameQuestion)
			finalAnswers = append(finalAnswers, additionalAnswers...)
		}
	}

	if len(finalAnswers) > 0 {
		return finalAnswers
	}

	log.Debug().Str("domain", q.Name).Str("type", dns.TypeToString[q.Qtype]).Msg("local record not found")
	return nil
}

func (l *LocalRecords) recursiveLookup(q *dns.Question) []dns.RR {
	for _, r := range l.Records {
		if r.Question.Name == q.Name && (r.Question.Qtype == q.Qtype || r.Question.Qtype == dns.TypeCNAME) {
			return r.Answer
		}
	}
	return nil
}
