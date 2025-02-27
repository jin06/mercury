package config

import (
	"os"
	"time"

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
	Listeners    []Listener   `yaml:"listeners"`
	ServerConfig ServerConfig `yaml:"server_config"`
}

type Listener struct {
	Type string `yaml:"type"`
	Addr string `yaml:"addr"`
}

type ServerConfig struct {
	// MaxConnections is the maximum number of connections the server will accept.
	MaxConnections int `yaml:"max_connections"`
	// MessageDeliveryTimeout is the maximum time in seconds the server will wait for a message to be delivered.
	MessageDeliveryTimeout time.Duration `yaml:"message_delivery_timeout"`
}
