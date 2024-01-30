package config

type Strategy string

const (
	StrategyRandom     Strategy = "random"
	StrategyRoundRobin Strategy = "round_robin"
	StrategyLatency    Strategy = "latency"
	StrategySequential Strategy = "sequential"
)

type Upstream struct {
	Strategy Strategy         `yaml:"strategy"`
	Servers  []UpstreamServer `yaml:"servers"`
}

type UpstreamServer struct {
	Name    string `yaml:"name"`
	IP      string `yaml:"ip"`
	Port    int    `yaml:"port"`
	Timeout int    `yaml:"timeout"`
}
