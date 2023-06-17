package main

import (
	"fmt"
	"github.com/go-vgo/robotgo"
	c "github.com/openstadia/openstadia/config"
	h "github.com/openstadia/openstadia/hub"
	l "github.com/openstadia/openstadia/local"
	r "github.com/openstadia/openstadia/rtc"
	"github.com/openstadia/openstadia/uinput"
	_ "github.com/pion/mediadevices/pkg/driver/screen"
	"os"
	"os/signal"
)

//sudo apt-get install libx11-dev libxext-dev libvpx-dev

//x11
//sudo apt install libx11-dev xorg-dev libxtst-dev

//Hook
//sudo apt install xcb libxcb-xkb-dev x11-xkb-utils libx11-xcb-dev libxkbcommon-x11-dev libxkbcommon-dev

//https://github.com/intel/media-driver

func main() {
	robotgo.MouseSleep = 0

	config, err := c.Load()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Config %#v\n", config)

	remoteGamepad := true
	var pGamepad *uinput.Gamepad

	if remoteGamepad {
		gamepad, err := uinput.CreateGamepad("/dev/uinput", []byte("Xbox One Wireless Controller"), 0x045E, 0x02EA)
		if err != nil {
			panic(err)
		}
		defer func(gamepad uinput.Gamepad) {
			err := gamepad.Close()
			if err != nil {
				panic(err)
			}
		}(gamepad)
		pGamepad = &gamepad
	}

	rtc := r.New(config, pGamepad)
	local := l.New(config, rtc)
	hub := h.New(config, rtc)

	if config.Hub != nil {
		interrupt := make(chan os.Signal, 1)
		signal.Notify(interrupt, os.Interrupt)

		go hub.Start(interrupt)
	}

	if config.Local != nil {
		go local.ServeHttp()
	}

	select {}
}
