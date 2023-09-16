package application

import (
	"github.com/pion/mediadevices"
	_ "github.com/pion/mediadevices/pkg/driver/camera"
	_ "github.com/pion/mediadevices/pkg/driver/microphone"
	"github.com/pion/mediadevices/pkg/frame"
	"github.com/pion/mediadevices/pkg/prop"
)

type cameraApp struct {
}

func CreateCamera() Application {
	return &cameraApp{}
}

func (c *cameraApp) Start() error {
	return nil
}

func (c *cameraApp) Stop() {

}

func (c *cameraApp) GetMedia(codecSelector *mediadevices.CodecSelector) (mediadevices.MediaStream, error) {
	return mediadevices.GetUserMedia(mediadevices.MediaStreamConstraints{
		Video: func(c *mediadevices.MediaTrackConstraints) {
			c.FrameFormat = prop.FrameFormat(frame.FormatMJPEG)
			c.Width = prop.Int(640)
			c.Height = prop.Int(480)
		},
		Audio: func(c *mediadevices.MediaTrackConstraints) {
		},
		Codec: codecSelector,
	})
}
