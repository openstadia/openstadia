package store

import (
	"errors"
	"fmt"
	c "github.com/openstadia/openstadia/config"
)

type ConfigStore struct {
	config *c.Openstadia
}

func CreateConfigStore(config *c.Openstadia) *ConfigStore {
	return &ConfigStore{config: config}
}

func (s *ConfigStore) Apps() []c.DbApp {
	var apps []c.DbApp

	for i, a := range s.config.Apps {
		dbApp := c.DbApp{
			Id:      i,
			BaseApp: a,
		}
		apps = append(apps, dbApp)
	}

	return apps
}

func (s *ConfigStore) Config() *c.DbOpenstadia {
	apps := s.Apps()
	hub := s.Hub()
	local := s.Local()

	config := c.DbOpenstadia{
		Openstadia: c.Openstadia{
			Hub:   hub,
			Local: local,
		},
		Apps: apps,
	}

	return &config
}

func (s *ConfigStore) GetAppById(id int) (*c.DbApp, error) {
	for _, app := range s.Apps() {
		if app.Id == id {
			return &app, nil
		}
	}

	return nil, fmt.Errorf("no such application: %d", id)
}

func (s *ConfigStore) Hub() *c.Hub {
	return s.config.Hub
}

func (s *ConfigStore) Local() *c.Local {
	return s.config.Local
}

func (s *ConfigStore) AddApp(app *c.BaseApp) error {
	return errors.New("store is not mutable")
}

func (s *ConfigStore) SetHub(hub *c.Hub) error {
	return errors.New("store is not mutable")
}

func (s *ConfigStore) SetLocal(local *c.Local) error {
	return errors.New("store is not mutable")
}
