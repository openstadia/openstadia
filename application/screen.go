package application

import (
	"github.com/openstadia/openstadia/driver/loopback"
	"github.com/pion/mediadevices"
	"github.com/pion/mediadevices/pkg/driver/camera"
	"github.com/pion/mediadevices/pkg/frame"
	"github.com/pion/mediadevices/pkg/prop"
	_ "unsafe"
)

//go:linkname selectScreen github.com/pion/mediadevices.selectScreen
func selectScreen(constraints mediadevices.MediaTrackConstraints, selector *mediadevices.CodecSelector) (mediadevices.Track, error)

//go:linkname selectAudio github.com/pion/mediadevices.selectAudio
func selectAudio(constraints mediadevices.MediaTrackConstraints, selector *mediadevices.CodecSelector) (mediadevices.Track, error)

func ScreenInitialize() {
	camera.Initialize()
	loopback.Initialize()
}

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
	return GetDisplayMedia(mediadevices.MediaStreamConstraints{
		Video: func(c *mediadevices.MediaTrackConstraints) {
			c.FrameFormat = prop.FrameFormat(frame.FormatRGBA)
			c.Width = prop.Int(640)
			c.Height = prop.Int(480)
			c.FrameRate = prop.Float(30)
		},
		Audio: func(c *mediadevices.MediaTrackConstraints) {
		},
		Codec: codecSelector,
	})
}

func GetDisplayMedia(constraints mediadevices.MediaStreamConstraints) (mediadevices.MediaStream, error) {
	trackers := make([]mediadevices.Track, 0)

	cleanTrackers := func() {
		for _, t := range trackers {
			t.Close()
		}
	}

	var videoConstraints, audioConstraints mediadevices.MediaTrackConstraints
	if constraints.Video != nil {
		constraints.Video(&videoConstraints)
		tracker, err := selectScreen(videoConstraints, constraints.Codec)
		if err != nil {
			cleanTrackers()
			return nil, err
		}

		trackers = append(trackers, tracker)
	}

	if constraints.Audio != nil {
		constraints.Audio(&audioConstraints)
		tracker, err := selectAudio(audioConstraints, constraints.Codec)
		if err != nil {
			cleanTrackers()
			return nil, err
		}

		trackers = append(trackers, tracker)
	}

	s, err := mediadevices.NewMediaStream(trackers...)
	if err != nil {
		cleanTrackers()
		return nil, err
	}

	return s, nil
}
