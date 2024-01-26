package config

import (
	"io/ioutil"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

var configPaths = []string{
	"/etc/go-resolve/config.yaml",
	os.Getenv("HOME") + "/.go-resolve/config.yaml",
	"./config.yaml",
}

type Config struct {
	Web struct {
		Port int `yaml:"port"`
	} `yaml:"web"`
	Logging struct {
		Level string `yaml:"level"`
	} `yaml:"logging"`
	DNS struct {
		DefaultTTL int      `yaml:"defaultTTL"`
		Upstream   []string `yaml:"upstream"`
	} `yaml:"dns"`
}

func findConfig() (string, error) {
	for _, path := range configPaths {
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}
	return "", os.ErrNotExist
}

func Load() (*Config, error) {
	// Find the config file...
	path, err := findConfig()
	if err != nil {
		return nil, err
	}

	// Read the config file...
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// Unmarshal the YAML data into the Config struct...
	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	// Log the loaded configuration...
	log.Printf("Web server will start at port %d\n", config.Web.Port)
	log.Printf("Logging level: %s\n", config.Logging.Level)
	log.Printf("DNS default TTL: %d\n", config.DNS.DefaultTTL)
	log.Printf("Upstream DNS servers: %v\n", config.DNS.Upstream)

	return &config, nil
}
