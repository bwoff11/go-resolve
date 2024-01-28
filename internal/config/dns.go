package config

type DNSConfig struct {
	TTL            int                `yaml:"TTL"`
	MaxMessageSize int                `yaml:"maxMessageSize"`
	BlockList      BlockListConfig    `yaml:"blockList"`
	Cache          CacheConfig        `yaml:"cache"`
	Local          LocalConfig        `yaml:"local"`
	Protocols      ProtocolConfigs    `yaml:"protocols"`
	RateLimiting   RateLimitingConfig `yaml:"rateLimiting"`
	Upstream       UpstreamConfig     `yaml:"upstream"`
}

type BlockListConfig struct {
	Local  []string `yaml:"local"`
	Remote []string `yaml:"remote"`
}

type CacheConfig struct {
	Enabled       bool   `yaml:"enabled"`
	Size          int    `yaml:"size"`
	PurgeInterval string `yaml:"purgeInterval"`
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

type RateLimitingConfig struct {
	Enabled           bool `yaml:"enabled"`
	RequestsPerSecond int  `yaml:"requestsPerSecond"`
	BurstSize         int  `yaml:"burstSize"`
}

type UpstreamConfig struct {
	Enabled  bool                  `yaml:"enabled"`
	Strategy LoadBalancingStrategy `yaml:"strategy"`
	Servers  []UpstreamServer      `yaml:"servers"`
}

type UpstreamServer struct {
	Name    string `yaml:"name"`
	IP      string `yaml:"ip"`
	Port    int    `yaml:"port"`
	Timeout int    `yaml:"timeout"`
}
