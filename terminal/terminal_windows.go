package terminal

import (
	"fmt"
	"github.com/UserExistsError/conpty"
	"golang.org/x/net/context"
	"io"
)

type Terminal struct {
	pty *conpty.ConPty
}

func New(width int, height int) (*Terminal, error) {
	commandLine := `c:\windows\system32\cmd.exe`

	conpty.ConPtyDimensions(width, height)
	pty, err := conpty.Start(commandLine)
	if err != nil {
		return nil, err
	}

	return &Terminal{
		pty: pty,
	}, nil
}

func (t *Terminal) Start(channel io.ReadWriteCloser) {
	go io.Copy(channel, t.pty)
	go io.Copy(t.pty, channel)

	go func() {
		wait, err := t.pty.Wait(context.Background())
		fmt.Printf("Terminal exited with %d\n", wait)
		if err != nil {
			return
		}
	}()
}

func (t *Terminal) Close() error {
	return t.pty.Close()
}
