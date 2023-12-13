package terminal

import (
	"github.com/kr/pty"
	"io"
	"log"
	"os"
	"os/exec"
	"sync"
	"syscall"
)

const (
	DEFAULT_SHELL string = "sh"
)

type Terminal struct {
	f   *os.File
	tty *os.File
}

func New(width int, height int) (*Terminal, error) {
	f, tty, err := pty.Open()
	if err != nil {
		log.Fatalf("could not start pty (%s)", err)
	}

	return &Terminal{
		f:   f,
		tty: tty,
	}, nil
}

func (t *Terminal) Start(channel io.ReadWriteCloser) {
	var shell string
	shell = os.Getenv("SHELL")
	if shell == "" {
		shell = DEFAULT_SHELL
	}

	cmd := exec.Command(shell)
	cmd.Env = []string{"TERM=xterm"}
	err := ptyRun(cmd, t.tty)
	if err != nil {
		log.Printf("%s", err)
	}

	var once sync.Once
	close := func() {
		channel.Close()
		log.Printf("session closed")
	}

	go func() {
		io.Copy(channel, t.f)
		once.Do(close)
	}()

	go func() {
		io.Copy(t.f, channel)
		once.Do(close)
	}()

}

func (t *Terminal) Close() error {
	return t.f.Close()
}

func ptyRun(c *exec.Cmd, tty *os.File) (err error) {
	c.Stdout = tty
	c.Stdin = tty
	c.Stderr = tty
	c.SysProcAttr = &syscall.SysProcAttr{
		Setctty: true,
		Setsid:  true,
	}
	return c.Start()
}
