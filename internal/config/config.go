package config

import (
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type ProtocolType string
type LoadBalancingStrategy string

const (
	ProtocolTypeDOT ProtocolType = "dot"
	ProtocolTypeUDP ProtocolType = "udp"
	ProtocolTypeTCP ProtocolType = "tcp"

	LoadBalancingStrategyRandom     LoadBalancingStrategy = "random"
	LoadBalancingStrategyRoundRobin LoadBalancingStrategy = "roundRobin"
	LoadBalancingStrategyLatency    LoadBalancingStrategy = "latency"
)

type Config struct {
	Web     WebConfig     `yaml:"web"`
	Logging LoggingConfig `yaml:"logging"`
	DNS     DNSConfig     `yaml:"dns"`
	Metrics MetricsConfig `yaml:"metrics"`
}

type WebConfig struct {
	Enabled bool         `yaml:"enabled"`
	Port    int          `yaml:"port"`
	TLS     WebTLSConfig `yaml:"tls"`
}

type WebTLSConfig struct {
	Enabled  bool   `yaml:"enabled"`
	CertFile string `yaml:"certFile"`
	KeyFile  string `yaml:"keyFile"`
}

type LoggingConfig struct {
	Level    string `yaml:"level"`
	Output   string `yaml:"output"`
	FilePath string `yaml:"filePath"`
}

type MetricsConfig struct {
	Enabled bool   `yaml:"enabled"`
	Route   string `yaml:"route"`
	Port    int    `yaml:"port"`
}

// Load reads the configuration file and unmarshals it into the Config struct.
func Load() (*Config, error) {
	v := viper.New()

	// Set the file name of the configurations file
	v.SetConfigName("config")
	v.SetConfigType("yml")

	// Add the path to look for the configurations file
	v.AddConfigPath("/etc/go-resolve/")
	v.AddConfigPath("$HOME/.go-resolve/")
	v.AddConfigPath(".")

	// Enable environment variable override of file settings
	v.AutomaticEnv()

	// Read the config file
	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	// Unmarshal into the Config struct
	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}

	log.Debug().
		Str("msg", "Loaded config").
		Str("config", fmt.Sprintf("%+v", cfg)).
		Send()

	return &cfg, nil
}
