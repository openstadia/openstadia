package types

import "fmt"

type Openstadia struct {
	Hub          string        `yaml:"hub" json:"hub"`
	Applications []Application `yaml:"applications" json:"applications"`
	Token        string        `yaml:"token" json:"token"`
}

func (o *Openstadia) GetApplications() []Application {
	return o.Applications
}

func (o *Openstadia) GetApplicationByName(name string) (*Application, error) {
	for _, app := range o.Applications {
		if app.Name == name {
			return &app, nil
		}
	}

	return nil, fmt.Errorf("no such application: %s", name)
}
