package store

import c "github.com/openstadia/openstadia/config"

type Store interface {
	Config() *c.DbOpenstadia

	Apps() []c.DbApp
	Hub() *c.Hub
	Local() *c.Local
	GetAppById(id int) (*c.DbApp, error)

	AddApp(app *c.BaseApp) error
	SetHub(hub *c.Hub) error
	SetLocal(local *c.Local) error
}
