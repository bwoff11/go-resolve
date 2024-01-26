package config

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Web     WebConfig
	Logging LoggingConfig
	DNS     DNSConfig
	Metrics MetricsConfig
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
	DefaultTTL     int
	MaxMessageSize int
	QueryTimeout   int
	DOT            DOTConfig
	UDP            ProtocolConfig
	TCP            ProtocolConfig
	Upstream       []string
	DNSSEC         DNSSECConfig
	RateLimiting   RateLimitingConfig
	LoadBalancing  LoadBalancingConfig
	HealthChecks   HealthChecksConfig
	Blacklist      BlacklistConfig
	Whitelist      WhitelistConfig
}

type DOTConfig struct {
	Enabled     bool
	Port        int
	TLSCertFile string `mapstructure:"tlsCertFile"`
	TLSKeyFile  string `mapstructure:"tlsKeyFile"`
	StrictSNI   bool   `mapstructure:"strictSNI"`
}

type ProtocolConfig struct {
	Enabled bool
	Port    int
	Cache   CacheConfig
}

type CacheConfig struct {
	TTL     int
	MaxSize int `mapstructure:"maxSize"`
}

type DNSSECConfig struct {
	Enabled             bool
	TrustAnchors        string `mapstructure:"trustAnchors"`
	ValidateRecursively bool   `mapstructure:"validateRecursively"`
}

type RateLimitingConfig struct {
	Enabled           bool
	RequestsPerSecond int `mapstructure:"requestsPerSecond"`
	BurstSize         int `mapstructure:"burstSize"`
}

type LoadBalancingConfig struct {
	Enabled bool
	Method  string
}

type HealthChecksConfig struct {
	Enabled   bool
	Interval  int // in seconds
	Endpoints []string
}

type BlacklistConfig struct {
	Enabled bool
	Domains []string
}

type WhitelistConfig struct {
	Enabled bool
	Domains []string
}

type MetricsConfig struct {
	Enabled  bool
	Endpoint string
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
