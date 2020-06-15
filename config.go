package main

import (
	"fmt"
	"io/ioutil"
	"time"

	"gopkg.in/yaml.v2"
)

type blockerConfig struct {
	Email               string `yaml:"email"`                // cloudfare email address
	APIKey              string `yaml:"api_key"`              // cloudfare API Key
	ZoneID              string `yaml:"zone_id,omitempty"`    // cloudfare zone ID
	AccountID           string `yaml:"account_id,omitempty"` // cloudfare account ID
	Scope               string `yaml:"scope"`
	DBPath              string `yaml:"dbpath"`
	PidDir              string `yaml:"piddir"`
	updateFrequency     time.Duration
	UpdateFrequencyYAML string `yaml:"update_frequency"`
	Daemon              bool   `yaml:"daemonize"`
	LogMode             string `yaml:"log_mode"`
	LogDir              string `yaml:"log_dir"`
}

func NewConfig(configPath string) (*blockerConfig, error) {
	config := &blockerConfig{}

	configBuff, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read %s : %v", configPath, err)
	}

	err = yaml.UnmarshalStrict(configBuff, &config)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal %s : %v", configPath, err)
	}

	config.updateFrequency, err = time.ParseDuration(config.UpdateFrequencyYAML)
	if err != nil {
		return nil, fmt.Errorf("invalid update frequency %s : %s", config.UpdateFrequencyYAML, err)
	}
	if config.DBPath == "" || config.Email == "" || config.PidDir == "" || config.LogMode == "" || config.APIKey == "" {
		return nil, fmt.Errorf("invalid configuration in %s, missing fields", configPath)
	}

	return config, nil
}
