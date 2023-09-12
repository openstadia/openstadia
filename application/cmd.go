package application

import (
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
