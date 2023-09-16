package application

import (
	"github.com/pion/mediadevices"
	"github.com/pion/mediadevices/pkg/frame"
	"github.com/pion/mediadevices/pkg/prop"
)

type screenApp struct {
}

func NewScreen() Application {
	return &screenApp{}
}

func (s *screenApp) Start() error {
	return nil
}

func (s *screenApp) Stop() {

}

func (s *screenApp) GetMedia(codecSelector *mediadevices.CodecSelector) (mediadevices.MediaStream, error) {
	return mediadevices.GetDisplayMedia(mediadevices.MediaStreamConstraints{
		Video: func(c *mediadevices.MediaTrackConstraints) {
			c.FrameFormat = prop.FrameFormat(frame.FormatRGBA)
			c.Width = prop.Int(640)
			c.Height = prop.Int(480)
			c.FrameRate = prop.Float(30)
		},
		Codec: codecSelector,
	})
}
