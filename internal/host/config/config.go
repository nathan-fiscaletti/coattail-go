package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Address struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

func (h Address) String() string {
	return fmt.Sprintf("%s:%d", h.Host, h.Port)
}

type ServiceConfig struct {
	LogPackets bool    `yaml:"log_packets"`
	Address    Address `yaml:"address"`
}

type ApiConfig struct {
	Enabled bool    `yaml:"enabled"`
	Address Address `yaml:"address"`
}

type WebConfig struct {
	Enabled bool    `yaml:"enabled"`
	Address Address `yaml:"address"`
}

type HostConfig struct {
	ServiceConfig ServiceConfig `yaml:"service"`
	ApiConfig     ApiConfig     `yaml:"api"`
	WebConfig     WebConfig     `yaml:"web"`
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
