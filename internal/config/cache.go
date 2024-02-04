package config

type Cache struct {
	Enabled       bool `yaml:"enabled"`
	PruneInterval int  `yaml:"pruneInterval"`
}
