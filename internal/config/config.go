package config

import (
	"os"
	"slices"
	"time"

	"github.com/jin06/mercury/internal/utils"
	"gopkg.in/yaml.v3"
)

var Def *Config

const (
	MemoryMode Mode = "memory"
)

func Init(path string) (err error) {
	Def, err = Parse(path)
	if err != nil {
		return
	}
	Def.MQTTConfig.MessageExpiryInterval = time.Hour * 24
	return
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
	MQTTConfig   MQTTConfig   `yaml:"mqtt"`
	Database     Database     `yaml:"database"`
	Mode         Mode         `yaml:"mode"`
	MessageStore MessageStore `yaml:"message_store"`
}

func (cfg *Config) Valid() (err error) {
	if err = cfg.Mode.Valid(); err != nil {
		return err
	}
	return nil
}

type Mode string

func (m Mode) Valid() error {
	if slices.Contains([]Mode{MemoryMode}, m) {
		return nil
	}
	return utils.ErrNotValidMode
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

type Database struct {
	Type string `json:"type"`
	DSN  string `json:"dsn"`
}

type BadgerConfig struct {
	Dir string `yaml:"dir"`
}

type MemoryConfig struct {
	Auth     bool   `yaml:"auth"`
	UserName string `yaml:"username"`
	Password string `yaml:"password"`
}

type MessageStore struct {
	Mode         string       `yaml:"mode"`
	BadgerConfig BadgerConfig `yaml:"badger"`
	MemoryConfig MemoryConfig `yaml:"moeory"`
}
