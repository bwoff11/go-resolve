package config

import (
	"errors"
	"fmt"

	"github.com/miekg/dns"
)

// Local groups DNS records by their type for easy management and parsing.
type Local struct {
	Standard []StandardRecord `yaml:"standard"`
}

// DNSRecord represents a DNS record with common fields.
type StandardRecord struct {
	Domain string `yaml:"domain"`
	Type   string `yaml:"type"`
	Value  string `yaml:"value"`
	TTL    int    `yaml:"ttl"`
}

func (sr StandardRecord) ToMsg() (*dns.Msg, error) {
	msg := new(dns.Msg)

	// Set the question section of the message.
	if rType, err := dns.StringToType[sr.Type]; err != true {
		return nil, errors.New("invalid record type: " + sr.Type)
	} else {
		msg.SetQuestion(sr.Domain+".", rType)
	}

	// Add the answer section of the message.
	rr, err := dns.NewRR(sr.Domain + "." + " " + fmt.Sprint(sr.TTL) + " IN " + sr.Type + " " + sr.Value)
	if err != nil {
		return nil, err
	}
	msg.Answer = append(msg.Answer, rr)

	return msg, nil
}
