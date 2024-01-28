package config

import (
	"fmt"
	"log"

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
	Enabled bool      `yaml:"enabled"`
	Port    int       `yaml:"port"`
	TLS     TLSConfig `yaml:"tls"`
}

type TLSConfig struct {
	Enabled  bool   `yaml:"enabled"`
	CertFile string `yaml:"certFile"`
	KeyFile  string `yaml:"keyFile"`
}

type LoggingConfig struct {
	Level    string `yaml:"level"`
	Output   string `yaml:"output"`
	FilePath string `yaml:"filePath"`
}

type DNSConfig struct {
	TTL            int                `yaml:"TTL"`
	MaxMessageSize int                `yaml:"maxMessageSize"`
	Cache          CacheConfig        `yaml:"cache"`
	Upstream       UpstreamConfig     `yaml:"upstream"`
	Protocols      ProtocolConfigs    `yaml:"protocols"`
	Local          LocalDNSConfig     `yaml:"local"`
	RateLimiting   RateLimitingConfig `yaml:"rateLimiting"`
	BlockList      BlockListConfig    `yaml:"blockList"`
}

type CacheConfig struct {
	Enabled       bool   `yaml:"enabled"`
	Size          int    `yaml:"size"`
	TTL           string `yaml:"ttl"`
	PurgeInterval string `yaml:"purgeInterval"`
}

type UpstreamConfig struct {
	Enabled  bool                  `yaml:"enabled"`
	Timeout  string                `yaml:"timeout"`
	Strategy LoadBalancingStrategy `yaml:"strategy"`
	Servers  []string              `yaml:"servers"`
}

type ProtocolConfigs struct {
	UDP ProtocolConfig `yaml:"udp"`
	TCP ProtocolConfig `yaml:"tcp"`
	Dot ProtocolConfig `yaml:"dot"`
}

type ProtocolConfig struct {
	Enabled     bool   `yaml:"enabled"`
	Port        int    `yaml:"port"`
	TLSCertFile string `yaml:"tlsCertFile,omitempty"`
	TLSKeyFile  string `yaml:"tlsKeyFile,omitempty"`
	StrictSNI   bool   `yaml:"strictSNI,omitempty"`
}

type LocalDNSConfig struct {
	Enabled bool       `yaml:"enabled"`
	Records DNSRecords `yaml:"records"`
}

type DNSRecords struct {
	A     []ARecord     `yaml:"a"`
	AAAA  []AAAARecord  `yaml:"aaaa"`
	CNAME []CNAMERecord `yaml:"cname"`
}

type ARecord struct {
	Domain string `yaml:"domain"`
	IP     string `yaml:"ip"`
}

type AAAARecord struct {
	Domain string `yaml:"domain"`
	IP     string `yaml:"ip"`
}

type CNAMERecord struct {
	Domain string `yaml:"domain"`
	Target string `yaml:"target"`
}

type RateLimitingConfig struct {
	Enabled           bool `yaml:"enabled"`
	RequestsPerSecond int  `yaml:"requestsPerSecond"`
	BurstSize         int  `yaml:"burstSize"`
}

type BlockListConfig struct {
	Local  BlockListDetail `yaml:"local"`
	Remote BlockListDetail `yaml:"remote"`
}

type BlockListDetail struct {
	Enabled bool     `yaml:"enabled"`
	Sources []string `yaml:"sources"`
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
