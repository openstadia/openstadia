package main

import (
	"golang.org/x/net/context"
	"log"
	"os"
	"os/exec"
)

type Application interface {
	Start() error
	Stop()
}

type application struct {
	cmd    *exec.Cmd
	cancel context.CancelFunc
}

func NewApplication(name string, arg []string, env []string) Application {
	ctx, cancel := context.WithCancel(context.Background())
	cmd := exec.CommandContext(ctx, name, arg...)
	cmd.Cancel = func() error {
		return cmd.Process.Signal(os.Interrupt)
	}
	cmd.Env = append(os.Environ(), env...)

	return &application{
		cmd:    cmd,
		cancel: cancel,
	}
}

func (a *application) Start() error {
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

func (a *application) Stop() {
	if a.cmd == nil {
		return
	}

	a.cancel()
}
