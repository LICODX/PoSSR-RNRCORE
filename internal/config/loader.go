package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config represents the structure of mainnet.yaml
type Config struct {
	Network  NetworkConfig `yaml:"network"`
	Sharding ShardConfig   `yaml:"sharding"`
}

type NetworkConfig struct {
	ListenPort int      `yaml:"listen_port"`
	SeedNodes  []string `yaml:"seed_nodes"`
	MaxPeers   int      `yaml:"max_peers"`
}

type ShardConfig struct {
	Role     string `yaml:"role"`      // "FullNode", "ShardNode"
	ShardIDs []int  `yaml:"shard_ids"` // List of shards to sync (0-9)
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
