//go:build windows && d3d

package screen

import (
	"errors"
	"fmt"
	"github.com/kirides/go-d3d/outputduplication"
	"github.com/kirides/go-d3d/win"
	"image"
	"io"

	"github.com/kbinani/screenshot"
	"github.com/kirides/go-d3d/d3d11"
	"github.com/pion/mediadevices/pkg/driver"
	"github.com/pion/mediadevices/pkg/frame"
	"github.com/pion/mediadevices/pkg/io/video"
	"github.com/pion/mediadevices/pkg/prop"
)

type screen struct {
	displayIndex int
	doneCh       chan struct{}
}

func init() {
	Initialize()
}

// Initialize finds and registers active displays. This is part of an experimental API.
func Initialize() {
	activeDisplays := screenshot.NumActiveDisplays()
	for i := 0; i < activeDisplays; i++ {
		priority := driver.PriorityNormal
		if i == 0 {
			priority = driver.PriorityHigh
		}

		s := newScreen(i)
		driver.GetManager().Register(s, driver.Info{
			Label:      fmt.Sprint(i),
			DeviceType: driver.Screen,
			Priority:   priority,
		})
	}
}

func newScreen(displayIndex int) *screen {
	s := screen{
		displayIndex: displayIndex,
	}
	return &s
}

func (s *screen) Open() error {
	s.doneCh = make(chan struct{})
	return nil
}

func (s *screen) Close() error {
	close(s.doneCh)
	return nil
}

func (s *screen) VideoRecord(selectedProp prop.Media) (video.Reader, error) {
	screenBounds := screenshot.GetDisplayBounds(s.displayIndex)

	if win.IsValidDpiAwarenessContext(win.DpiAwarenessContextPerMonitorAwareV2) {
		_, err := win.SetThreadDpiAwarenessContext(win.DpiAwarenessContextPerMonitorAwareV2)
		if err != nil {
			fmt.Printf("Could not set thread DPI awareness to PerMonitorAwareV2. %v\n", err)
		} else {
			fmt.Printf("Enabled PerMonitorAwareV2 DPI awareness.\n")
		}
	}

	device, deviceCtx, err := d3d11.NewD3D11Device()
	if err != nil {
		fmt.Printf("Could not create D3D11 Device. %v\n", err)
		return nil, err
	}

	ddup, err := outputduplication.NewIDXGIOutputDuplication(device, deviceCtx, uint(s.displayIndex))
	if err != nil {
		fmt.Printf("Err NewIDXGIOutputDuplication: %v\n", err)
		return nil, err
	}

	// TODO Add support for mouse enable
	//ddup.DrawPointer = true

	imgBuf := image.NewRGBA(screenBounds)

	r := video.ReaderFunc(func() (img image.Image, release func(), err error) {
		select {
		case <-s.doneCh:
			return nil, nil, io.EOF
		default:
		}

		err = ddup.GetImage(imgBuf, 0)
		if err != nil && !errors.Is(err, outputduplication.ErrNoImageYet) {
			fmt.Printf("Err ddup.GetImage: %v\n", err)
			//return nil, nil, err
		}
		err = nil
		img = imgBuf
		release = func() {}
		return
	})
	return r, nil
}

func (s *screen) Properties() []prop.Media {
	resolution := screenshot.GetDisplayBounds(s.displayIndex)
	supportedProp := prop.Media{
		Video: prop.Video{
			Width:       resolution.Dx(),
			Height:      resolution.Dy(),
			FrameFormat: frame.FormatRGBA,
		},
	}
	return []prop.Media{supportedProp}
}
