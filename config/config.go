package config

var Def Config

func init() {
	Def = Config{}
}

func Parse(path string) Config {
	cfg := Config{}
	return cfg
}

type Config struct {
	Listeners []Listener `yaml:"listeners"`
}

type Listener struct {
	Type string `yaml:"type"`
	Addr string
	Port int `yaml:"port"`
}
