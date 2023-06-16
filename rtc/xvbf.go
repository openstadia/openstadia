package rtc

import (
	"fmt"
	"golang.org/x/net/context"
	"log"
	"os"
	"os/exec"
	"syscall"
	"time"
)

type Xvfb interface {
	Start()
	Stop()
}

type xvfb struct {
	cmd        *exec.Cmd
	display    string
	displayNum uint
	cancel     context.CancelFunc
	width      uint
	height     uint
}

func NewXvfb(displayNum uint, width uint, height uint) Xvfb {
	return &xvfb{
		cmd:        nil,
		display:    fmt.Sprintf(":%d", displayNum),
		displayNum: displayNum,
		cancel:     nil,
		width:      width,
		height:     height,
	}
}

func (x *xvfb) Start() {
	if x.cmd != nil {
		return
	}

	x.setDisplayEnvVariable()
	x.spawnProcess()

	for !x.isDisplayExists() {
		time.Sleep(time.Millisecond * 100)
	}
}

func (x *xvfb) Stop() {
	if x.cmd == nil {
		return
	}

	x.cancel()
}

func (x *xvfb) isDisplayExists() bool {
	lockFile := x.getLockFile()
	_, err := os.Stat(lockFile)
	return err == nil
}

func (x *xvfb) setDisplayEnvVariable() {
	err := os.Setenv("DISPLAY", x.display)
	if err != nil {
		panic(err)
	}
}

func (x *xvfb) spawnProcess() {
	display := x.display
	ctx, cancel := context.WithCancel(context.Background())
	cmd := exec.CommandContext(ctx, "Xvfb", display, "-screen", "0", fmt.Sprintf("%dx%dx24", x.width, x.height))
	cmd.Cancel = func() error {
		return cmd.Process.Signal(os.Interrupt)
	}
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	x.cmd = cmd
	x.cancel = cancel

	err := cmd.Start()
	if err != nil {
		panic(err)
	}

	go func() {
		err := cmd.Wait()
		if err != nil {
			log.Println(cmd.Output())
			log.Printf("process finished with error = %v\n", err)
		}
		log.Print("process finished successfully")
	}()
}

func (x *xvfb) getLockFile() string {
	return fmt.Sprintf("/tmp/.X%d-lock", x.displayNum)
}
