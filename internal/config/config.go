package config

import (
	"fmt"
	"log"
	"time"

	"github.com/spf13/viper"
)

type ProtocolType string

const (
	ProtocolTypeDOT ProtocolType = "dot"
	ProtocolTypeUDP ProtocolType = "udp"
	ProtocolTypeTCP ProtocolType = "tcp"
)

type LoadBalancingStrategy string

const (
	LoadBalancingStrategyRandom     LoadBalancingStrategy = "random"
	LoadBalancingStrategyRoundRobin LoadBalancingStrategy = "roundRobin"
	LoadBalancingStrategyLatency    LoadBalancingStrategy = "latency"
)

type Config struct {
	Web     WebConfig
	Logging LoggingConfig
	DNS     DNSConfig
}

type WebConfig struct {
	Enabled bool
	Port    int
	TLS     TLSConfig
}

type TLSConfig struct {
	Enabled  bool
	CertFile string
	KeyFile  string
}

type LoggingConfig struct {
	Level    string
	Output   string
	FilePath string
}

type DNSConfig struct {
	TTL            int
	MaxMessageSize int
	Cache          CacheConfig
	Upstream       UpstreamConfig
	Protocols      ProtocolConfigs
	LocalDNSConfig LocalDNSConfig
}

type LocalDNSConfig struct {
	Enabled bool
	Records DNSRecords
}

type DNSRecords struct {
	A     map[string]string
	AAAA  map[string]string
	CNAME map[string]string
}

type CacheConfig struct {
	Enabled       bool
	Size          int
	TTL           time.Duration
	PurgeInterval time.Duration
}

type UpstreamConfig struct {
	Enabled  bool
	Timeout  time.Duration
	Strategy LoadBalancingStrategy
	Servers  []string
}

type ProtocolConfigs struct {
	UDP ProtocolConfig
	TCP ProtocolConfig
	DOT ProtocolConfig
}

type ProtocolConfig struct {
	Enabled     bool
	Port        int
	TLSCertFile string `mapstructure:"tlsCertFile,omitempty"`
	TLSKeyFile  string `mapstructure:"tlsKeyFile,omitempty"`
	StrictSNI   bool   `mapstructure:"strictSNI,omitempty"`
}

// Load reads the configuration file and unmarshals it into the Config struct.
func Load() (*Config, error) {
	v := viper.New()

	// Set the file name of the configurations file
	v.SetConfigName("config")
	v.SetConfigType("yaml")

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

	log.Printf("Configuration loaded: %+v\n", cfg)
	return &cfg, nil
}
