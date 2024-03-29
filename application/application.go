package application

import (
	"github.com/openstadia/openstadia/config"
	"github.com/pion/mediadevices"
)

type Application interface {
	Start() error
	Stop()
	GetMedia(codecSelector *mediadevices.CodecSelector) (mediadevices.MediaStream, error)
}

func IsScreen(app *config.DbApp) bool {
	return app.Name == "screen"
}

func IsCamera(app *config.DbApp) bool {
	return app.Name == "camera"
}

func IsApp(app *config.DbApp) bool {
	return true
}
