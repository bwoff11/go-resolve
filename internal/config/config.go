package config

import (
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type Config struct {
	BlockLists []string  `yaml:"blockLists"`
	Cache      Cache     `yaml:"cache"`
	Local      Local     `yaml:"local"`
	Metrics    Metrics   `yaml:"metrics"`
	Protocols  Protocols `yaml:"protocols"`
	Upstream   Upstream  `yaml:"upstream"`
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
