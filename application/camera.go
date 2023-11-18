package application

import (
	"github.com/pion/mediadevices"
	"github.com/pion/mediadevices/pkg/driver/camera"
	"github.com/pion/mediadevices/pkg/frame"
	"github.com/pion/mediadevices/pkg/prop"
)

func CameraInitialize() {
	camera.Initialize()
}

type cameraApp struct {
}

func NewCamera() Application {
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
		Codec: codecSelector,
	})
}
