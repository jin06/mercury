package config

import "testing"

func TestParse(t *testing.T) {
	cfg, err := Parse("mercury_test.yaml")
	if err != nil {
		t.Error(err)
	}
	if cfg.Listeners[0].Type != "tcp" {
		t.Fail()
	}
	t.Log(cfg)
}
