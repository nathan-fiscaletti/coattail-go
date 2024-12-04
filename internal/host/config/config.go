package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type HostConfig struct {
	ServiceAddress string `yaml:"service_address"`
	LogPackets     bool   `yaml:"log_packets"`
	WebAddress     string `yaml:"web_address"`
	ApiAddress     string `yaml:"api_address"`
}

func GetHostConfig() (*HostConfig, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	hostConfigFile, err := os.ReadFile(filepath.Join(cwd, "host-config.yaml"))
	if err != nil {
		return nil, err
	}

	cfg := HostConfig{}
	err = yaml.Unmarshal(hostConfigFile, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
