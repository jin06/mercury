package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

var Def *Config

func Init(path string) error {
	cfg, err := Parse(path)
	if err != nil {
		return err
	}
	Def = cfg
	return nil
}

func Parse(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	cfg := &Config{}
	err = yaml.NewDecoder(file).Decode(cfg)
	return cfg, err
}

type Config struct {
	Listeners []Listener `yaml:"listeners"`
}

type Listener struct {
	Type string `yaml:"type"`
	Addr string
	Port int `yaml:"port"`
}
