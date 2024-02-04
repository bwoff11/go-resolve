package config

type Cache struct {
	Enabled       bool          `yaml:"enabled"`
	PruneInterval int           `yaml:"pruneInterval"`
	LocalRecords  []LocalRecord `yaml:"localRecords"`
}

type LocalRecord struct {
	Domain string `yaml:"domain"`
	Type   string `yaml:"type"`
	Value  string `yaml:"value"`
	TTL    int    `yaml:"ttl"`
}
