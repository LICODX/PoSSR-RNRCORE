package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config represents the structure of mainnet.yaml
type Config struct {
	Network NetworkConfig `yaml:"network"`
}

type NetworkConfig struct {
	ListenPort int      `yaml:"listen_port"`
	SeedNodes  []string `yaml:"seed_nodes"`
	MaxPeers   int      `yaml:"max_peers"`
}

// LoadConfig reads and parses a YAML configuration file
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %v", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %v", err)
	}

	return &cfg, nil
}
