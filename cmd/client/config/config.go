package config

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/imdario/mergo"
	homedir "github.com/mitchellh/go-homedir"
)

var (
	DeploiConfiguration *Configuration
)

type Configuration struct {
	Location       string
	Host           string
	DialSecurely   bool
	TLSCertificate string
	Token          string
}

func NewConfig() *Configuration {
	return &Configuration{
		Host:     "localhost:8000",
		Token:    "",
		Location: "",
	}
}

func NewConfigFromFile(path string) (*Configuration, error) {
	conf := NewConfig()
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("Failed to open config file: %s", err)
	}
	fileConf := &Configuration{}
	if err = json.NewDecoder(file).Decode(fileConf); err != nil {
		return nil, fmt.Errorf("Failed to parse config file: %s", err)
	}
	if err = mergo.MergeWithOverwrite(conf, fileConf); err != nil {
		return nil, fmt.Errorf("Failed to merge default config and file config: %s", err)
	}
	conf.Location = path
	return conf, nil
}

func WriteConfig(conf *Configuration) error {
	if conf.Location == "" {
		return fmt.Errorf("No location specified, cannot write")
	}
	file, err := os.Create(conf.Location)
	if err != nil {
		return fmt.Errorf("Failed to open/create file: %s", err)
	}
	// Location is not stored in the file but instead always set on runtime
	conf.Location = ""
	if err = json.NewEncoder(file).Encode(conf); err != nil {
		return fmt.Errorf("Failed to encode config: %s", err)
	}
	return nil
}

func GetDefaultConfLocation() string {
	home, err := homedir.Dir()
	if err != nil {
		home = "/"
	}
	return fmt.Sprintf("%s/%s", home, ".deploi")
}
