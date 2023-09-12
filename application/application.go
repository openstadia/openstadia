package application

import "github.com/openstadia/openstadia/config"

type Application interface {
	Start() error
	Stop()
}

func IsScreen(app *config.App) bool {
	return app.Name == "screen"
}
