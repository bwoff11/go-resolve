package config

import (
	"fmt"
	"log"
	"time"

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
	Web     WebConfig     `mapstructure:"web"`
	Logging LoggingConfig `mapstructure:"logging"`
	DNS     DNSConfig     `mapstructure:"dns"`
}

type WebConfig struct {
	Enabled bool      `mapstructure:"enabled"`
	Port    int       `mapstructure:"port"`
	TLS     TLSConfig `mapstructure:"tls"`
}

type TLSConfig struct {
	Enabled  bool   `mapstructure:"enabled"`
	CertFile string `mapstructure:"certFile"`
	KeyFile  string `mapstructure:"keyFile"`
}

type LoggingConfig struct {
	Level    string `mapstructure:"level"`
	Output   string `mapstructure:"output"`
	FilePath string `mapstructure:"filePath"`
}

type DNSConfig struct {
	TTL            int             `mapstructure:"ttl"`
	MaxMessageSize int             `mapstructure:"maxMessageSize"`
	Cache          CacheConfig     `mapstructure:"cache"`
	Upstream       UpstreamConfig  `mapstructure:"upstream"`
	Protocols      ProtocolConfigs `mapstructure:"protocols"`
	LocalDNSConfig LocalDNSConfig  `mapstructure:"local"`
}

type LocalDNSConfig struct {
	Enabled bool        `mapstructure:"enabled"`
	Records []DNSRecord `mapstructure:"records"`
}

type DNSRecord struct {
	Type    string `mapstructure:"type"`
	Domain  string `mapstructure:"domain"`
	Address string `mapstructure:"address,omitempty"`
	Target  string `mapstructure:"target,omitempty"`
}

type CacheConfig struct {
	Enabled       bool          `mapstructure:"enabled"`
	Size          int           `mapstructure:"size"`
	TTL           time.Duration `mapstructure:"ttl"`
	PurgeInterval time.Duration `mapstructure:"purgeInterval"`
}

type UpstreamConfig struct {
	Enabled  bool                  `mapstructure:"enabled"`
	Timeout  time.Duration         `mapstructure:"timeout"`
	Strategy LoadBalancingStrategy `mapstructure:"strategy"`
	Servers  []string              `mapstructure:"servers"`
}

type ProtocolConfigs struct {
	UDP ProtocolConfig `mapstructure:"udp"`
	TCP ProtocolConfig `mapstructure:"tcp"`
	DOT ProtocolConfig `mapstructure:"dot"`
}

type ProtocolConfig struct {
	Enabled     bool   `mapstructure:"enabled"`
	Port        int    `mapstructure:"port"`
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
