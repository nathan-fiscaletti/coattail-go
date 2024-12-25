package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type HostAddress struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

func (h HostAddress) String() string {
	return fmt.Sprintf("%s:%d", h.Host, h.Port)
}

type HostConfig struct {
	ServiceAddress HostAddress `yaml:"service_address"`
	WebAddress     HostAddress `yaml:"web_address"`
	ApiAddress     HostAddress `yaml:"api_address"`
	LogPackets     bool        `yaml:"log_packets"`
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
