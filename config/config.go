package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Application struct {
		Server struct {
			Host string `yaml:"host"`
			Port int    `yaml:"port"`
      Protocol string `yaml:"protocol"`
		} `yaml:"server"`
	} `yaml:"application"`
}

func ReadConfig(filename string) (*Config, error) {
	buf, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var cfg Config
	err = yaml.Unmarshal(buf, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
