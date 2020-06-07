package config

import (
	"testing"
)

func TestConfigFileLoad(t *testing.T) {
	config := "../test/scm_config.json"
	c, err := LoadChangeLogConfig(config)
	if err != nil {
		t.Errorf("error %s", err)
	} else {
		t.Log(c)
	}
}
