package application

import (
	"github.com/openstadia/openstadia/driver/loopback"
	"github.com/openstadia/openstadia/driver/screen"
)

func ScreenInitialize() {
	screen.Initialize()
	loopback.Initialize()
}
