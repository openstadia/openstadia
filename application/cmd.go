package application

import (
	"github.com/pion/mediadevices"
	"github.com/pion/mediadevices/pkg/frame"
	"github.com/pion/mediadevices/pkg/prop"
	"golang.org/x/net/context"
	"log"
	"os"
	"os/exec"
)

type cmdApp struct {
	cmd    *exec.Cmd
	cancel context.CancelFunc
}

func NewCmd(name string, arg []string, env []string) Application {
	ctx, cancel := context.WithCancel(context.Background())
	cmd := exec.CommandContext(ctx, name, arg...)
	cmd.Cancel = func() error {
		return cmd.Process.Signal(os.Interrupt)
	}
	cmd.Env = append(os.Environ(), env...)

	return &cmdApp{
		cmd:    cmd,
		cancel: cancel,
	}
}

func (a *cmdApp) Start() error {
	err := a.cmd.Start()
	if err != nil {
		return err
	}

	go func() {
		err := a.cmd.Wait()
		if err != nil {
			log.Printf("process finished with error = %v\n", err)
		}
		log.Print("process finished successfully")
	}()

	return nil
}

func (a *cmdApp) Stop() {
	if a.cmd == nil {
		return
	}

	a.cancel()
}

func (a *cmdApp) GetMedia(codecSelector *mediadevices.CodecSelector) (mediadevices.MediaStream, error) {
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
