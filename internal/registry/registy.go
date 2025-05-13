package registry

import (
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type ServiceConfig struct {
	Name       string   `yaml:"name"`
	PathPrefix string   `yaml:"path_prefix"`
	Targets    []string `yaml:"targets"`
}

type Config struct {
	BindAddress string `yaml:"bind_address"`
	HealthCheck struct {
		Path     string        `yaml:"path"`
		Interval time.Duration `yaml:"interval"`
	} `yaml:"health_check"`
	Services []ServiceConfig `yaml:"services"`
}

func LoadConfig(path string) (*Config, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(b, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
