package broker

import (
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	DataDir string `yaml:"data_dir"`
	Port    int    `yaml:"port"`
}

func LoadConfig(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var cfg Config
	if err := yaml.NewDecoder(f).Decode(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
