package cache

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/miekg/dns"
)

func (c *Cache) GenerateKey(q dns.Question) string {
	return fmt.Sprintf("%s:%d", q.Name, q.Qtype)
}

func (c *Cache) DecodeKey(key string) (domain string, rType uint16) {
	s := strings.Split(key, ":")
	d := s[0]
	t, _ := strconv.ParseUint(s[1], 10, 16)
	return d, uint16(t)
}
