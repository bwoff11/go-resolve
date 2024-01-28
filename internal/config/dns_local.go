package config

type LocalConfig struct {
	Enabled bool     `yaml:"enabled"`
	Records []Record `yaml:"records"`
}

type Record struct {
	Domain string `yaml:"domain"`
	Type   string `yaml:"type"`
	Value  string `yaml:"value"`
	TTL    int    `yaml:"ttl"`
}
