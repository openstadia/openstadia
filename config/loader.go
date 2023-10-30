package config

import (
	"errors"
	"gopkg.in/yaml.v3"
	"os"
)

const ConfigFileName = "openstadia.yaml"

var (
	ErrNoConfigFile = errors.New("config not provided")
)

func Load() (*Openstadia, error) {
	config := Openstadia{}

	if _, err := os.Stat(ConfigFileName); err == nil {
		yamlFile, err := os.ReadFile(ConfigFileName)
		if err != nil {
			return nil, err
		}

		err = yaml.Unmarshal(yamlFile, &config)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, ErrNoConfigFile
	}

	return &config, nil
}
