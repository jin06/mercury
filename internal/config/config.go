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
	Listeners  []Listener `yaml:"listeners"`
	MQTTConfig MQTTConfig `yaml:"mqtt"`
	DBConfig   DBConfig   `yaml:"db"`
}

type Listener struct {
	Type string `yaml:"type"`
	Addr string `yaml:"addr"`
}

type MQTTConfig struct {
	// MaxConnections is the maximum number of connections the server will accept.
	MaxConnections int `yaml:"max_connections"`
	// MessageDeliveryTimeout is the maximum time in seconds the server will wait for a message to be delivered.
	MessageDeliveryTimeout time.Duration `yaml:"message_delivery_timeout"`
	MessageExpiryInterval  time.Duration `yaml:"message_expiry_interval"`
}

type DBConfig struct {
	Driver   string `yaml:"driver"` //  mysql,postgres
	User     string `yaml:"user"`   //
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Database string `yaml:"database"`
}

type BadgerConfig struct {
	Dir string `yaml:"dir"`
}
