package config

type Transport struct {
	UDP Protocol    `yaml:"udp"`
	TCP Protocol    `yaml:"tcp"`
	DOT DOTProtocol `yaml:"dot"`
	DOH DOHProtocol `yaml:"doh"`
}

type Protocol struct {
	Enabled bool `yaml:"enabled"`
	Port    int  `yaml:"port"`
}

type DOTProtocol struct {
	Protocol    `yaml:",inline"`
	TLSCertFile string `yaml:"tlsCertFile"`
	TLSKeyFile  string `yaml:"tlsKeyFile"`
	StrictSNI   bool   `yaml:"strictSNI"`
}

type DOHProtocol struct {
	Protocol    `yaml:",inline"`
	TLSCertFile string `yaml:"tlsCertFile"`
	TLSKeyFile  string `yaml:"tlsKeyFile"`
	Endpoint    string `yaml:"endpoint"`
}
