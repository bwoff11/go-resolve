package config

type Protocols struct {
	UDP Protocol    `yaml:"udp"`
	TCP Protocol    `yaml:"tcp"`
	Dot DOTProtocol `yaml:"dot"`
}

type Protocol struct {
	Enabled bool `yaml:"enabled"`
	Port    int  `yaml:"port"`
}

type DOTProtocol struct {
	Enabled     bool   `yaml:"enabled"`
	Port        int    `yaml:"port"`
	TLSCertFile string `yaml:"tlsCertFile"`
	TLSKeyFile  string `yaml:"tlsKeyFile"`
	StrictSNI   bool   `yaml:"strictSNI"`
}
