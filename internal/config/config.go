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
	DBConfig     DBConfig     `yaml:"db"`
	Mode         Mode         `yaml:"mode"`
	MemoryConfig MemoryConfig `yaml:"memory_mode"`
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

type MemoryConfig struct {
	Auth     bool   `yaml:"auth"`
	UserName string `yaml:"username"`
	Password string `yaml:"password"`
}
