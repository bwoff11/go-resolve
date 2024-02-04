package config

// Local groups DNS records by their type for easy management and parsing.
type Local struct {
	Standard []StandardRecord `yaml:"standard"`
}

// DNSRecord represents a DNS record with common fields.
type StandardRecord struct {
	Type   string `yaml:"type"`
	Domain string `yaml:"domain"`
	Value  string `yaml:"value"`
	TTL    int    `yaml:"ttl"`
}
