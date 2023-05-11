package types

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

func Load() (*Openstadia, error) {
	config := Openstadia{}

	yamlFile, err := ioutil.ReadFile("openstadia.yaml")
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
