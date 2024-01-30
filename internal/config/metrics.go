package config

type Metrics struct {
	Enabled bool   `yaml:"enabled"`
	Route   string `yaml:"route"`
	Port    int    `yaml:"port"`
}
