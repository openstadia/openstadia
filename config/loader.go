package config

import (
	"gopkg.in/yaml.v3"
	"os"
)

func Load() (*Openstadia, error) {
	config := Openstadia{}

	yamlFile, err := os.ReadFile("openstadia.yaml")
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
