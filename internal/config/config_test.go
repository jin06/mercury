package config

import "testing"

func TestParse(t *testing.T) {
	cfg, err := Parse("mercury_test.yaml")
	if err != nil {
		t.Error(err)
	}
	t.Log(cfg)
}
