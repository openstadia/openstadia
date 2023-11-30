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
}

func New(width int, height int) *Terminal {
	f, tty, err := pty.Open()
	if err != nil {
		log.Fatalf("could not start pty (%s)", err)
	}

	var shell string
	shell = os.Getenv("SHELL")
	if shell == "" {
		shell = DEFAULT_SHELL
	}

	cmd := exec.Command(shell)
	cmd.Env = []string{"TERM=xterm"}
	err = ptyRun(cmd, tty)
	if err != nil {
		log.Printf("%s", err)
	}

	w := newDcw(d)

	var once sync.Once
	close := func() {
		d.Close()
		log.Printf("session closed")
	}

	go func() {
		io.Copy(w, f)
		once.Do(close)
	}()

	go func() {
		io.Copy(f, w)
		once.Do(close)
	}()
}

func ptyRun(c *exec.Cmd, tty *os.File) (err error) {
	defer tty.Close()
	c.Stdout = tty
	c.Stdin = tty
	c.Stderr = tty
	c.SysProcAttr = &syscall.SysProcAttr{
		Setctty: true,
		Setsid:  true,
	}
	return c.Start()
}
