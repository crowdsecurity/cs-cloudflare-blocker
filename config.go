package main

import (
	"fmt"
	"io/ioutil"
	"time"

	"gopkg.in/yaml.v2"
)

type blockerConfig struct {
	Email     string `yaml:"email"`                // cloudflare email address
	APIKey    string `yaml:"api_key"`              // cloudflare API Key
	ZoneID    string `yaml:"zone_id,omitempty"`    // cloudflare zone ID
	AccountID string `yaml:"account_id,omitempty"` // cloudflare account ID
	Scope     string `yaml:"scope"`
	//DBPath              string `yaml:"dbpath"`
	PidDir              string `yaml:"piddir"`
	updateFrequency     time.Duration
	UpdateFrequencyYAML string            `yaml:"update_frequency"`
	Daemon              bool              `yaml:"daemonize"`
	LogMode             string            `yaml:"log_mode"`
	LogDir              string            `yaml:"log_dir"`
	DBConfig            map[string]string `yaml:"db_config"`
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
	if len(config.DBConfig) == 0 || config.Email == "" || config.PidDir == "" || config.LogMode == "" || config.APIKey == "" {
		return nil, fmt.Errorf("invalid configuration in %s, missing fields", configPath)
	}

	return config, nil
}
