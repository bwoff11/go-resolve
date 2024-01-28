package config

type LocalConfig struct {
	Standard []StandardRecord `yaml:"standard"`
	Wildcard []WildcardRecord `yaml:"wildcard"`
}

type StandardRecord struct {
	Domain string `yaml:"domain"`
	Type   string `yaml:"type"`
	Value  string `yaml:"value"`
	TTL    int    `yaml:"ttl"`
}

type WildcardRecord struct {
	Pattern string `yaml:"pattern"`
	Type    string `yaml:"type"`
	Value   string `yaml:"value"`
	TTL     int    `yaml:"ttl"`
}
