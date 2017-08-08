package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/imdario/mergo"
)

type Config struct {
	ListenAddr     string
	ListenSecurely bool
	TLS            *TLSConfig
	Database       *DatabaseConfig
	Retention      *RetentionConfig
}

type DatabaseConfig struct {
	Path string
}

type RetentionConfig struct {
	Jobs             int
	Builds           int
	BuildsPerProject map[string]int
}

type TLSConfig struct {
	CertFile string
	KeyFile  string
}

func newConfig() *Config {
	return &Config{
		ListenAddr:     ":8000",
		ListenSecurely: false,
		Database: &DatabaseConfig{
			Path: "/var/lib/deploid/deploid.db",
		},
		Retention: &RetentionConfig{
			Jobs:             200,
			Builds:           50,
			BuildsPerProject: map[string]int{},
		},
	}
}

func newConfigFromFile(path string) (*Config, error) {
	conf := newConfig()
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("Failed to open config file: %s", err)
	}
	fileConf := &Config{}
	if err = json.NewDecoder(file).Decode(fileConf); err != nil {
		return nil, fmt.Errorf("Failed to parse config file: %s", err)
	}
	if err = mergo.MergeWithOverwrite(conf, fileConf); err != nil {
		return nil, fmt.Errorf("Failed to merge default config and file config: %s", err)
	}
	return conf, nil
}
