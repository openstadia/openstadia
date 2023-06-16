package config

import "fmt"

type Openstadia struct {
	Hub   *Hub   `yaml:"hub" json:"hub"`
	Local *Local `yaml:"local" json:"local"`
	Apps  []App  `yaml:"apps" json:"apps"`
}

func (o *Openstadia) GetApps() []App {
	return o.Apps
}

func (o *Openstadia) GetAppByName(name string) (*App, error) {
	for _, app := range o.Apps {
		if app.Name == name {
			return &app, nil
		}
	}

	return nil, fmt.Errorf("no such application: %s", name)
}
