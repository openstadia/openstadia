package main

import (
	"fmt"
	"golang.org/x/net/context"
	"log"
	"os"
	"os/exec"
	"time"
)

type Ppsspp interface {
	Start()
	Stop()
}

type ppsspp struct {
	cmd        *exec.Cmd
	display    string
	displayNum uint
	cancel     context.CancelFunc
}

func NewPpsspp(displayNum uint) Ppsspp {
	return &ppsspp{
		cmd:        nil,
		display:    fmt.Sprintf(":%d", displayNum),
		displayNum: displayNum,
		cancel:     nil,
	}
}

func (x *ppsspp) Start() {
	if x.cmd != nil {
		return
	}

	x.setDisplayEnvVariable()
	x.spawnProcess()

	for !x.isDisplayExists() {
		time.Sleep(time.Millisecond * 100)
	}
}

func (x *ppsspp) Stop() {
	if x.cmd == nil {
		return
	}

	x.cancel()
}

func (x *ppsspp) isDisplayExists() bool {
	lockFile := x.getLockFile()
	_, err := os.Stat(lockFile)
	return err == nil
}

func (x *ppsspp) setDisplayEnvVariable() {
	err := os.Setenv("DISPLAY", x.display)
	if err != nil {
		panic(err)
	}
}

//480x272x24
//720x408x24
//960x544x24

//640x480x24

func (x *ppsspp) spawnProcess() {
	display := x.display
	ctx, cancel := context.WithCancel(context.Background())
	cmd := exec.CommandContext(ctx, "/home/user/ppsspp_build/PPSSPPSDL")
	cmd.Cancel = func() error {
		return cmd.Process.Signal(os.Interrupt)
	}

	cmd.Env = append(os.Environ(), "DISPLAY="+display)

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

func (x *ppsspp) getLockFile() string {
	return fmt.Sprintf("/tmp/.X%d-lock", x.displayNum)
}
