package main

import (
	"fmt"
	c "github.com/openstadia/openstadia/config"
	h "github.com/openstadia/openstadia/hub"
	g "github.com/openstadia/openstadia/inputs/gamepad"
	"github.com/openstadia/openstadia/inputs/keyboard"
	"github.com/openstadia/openstadia/inputs/mouse"
	l "github.com/openstadia/openstadia/local"
	r "github.com/openstadia/openstadia/rtc"
	_ "github.com/pion/mediadevices/pkg/driver/screen"
	"os"
	"os/signal"
)

//Hook
//sudo apt install xcb libxcb-xkb-dev x11-xkb-utils libx11-xcb-dev libxkbcommon-x11-dev libxkbcommon-dev

//https://github.com/intel/media-driver

func main() {
	config, err := c.Load()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Config %#v\n", config)

	remoteGamepad := true
	var gamepad_ g.Gamepad

	if remoteGamepad {
		gamepad, err := g.CreateGamepad("/dev/uinput", []byte("Xbox One Wireless Controller"), 0x045E, 0x02EA)
		if err != nil {
			panic(err)
		}
		defer func(gamepad g.Gamepad) {
			err := gamepad.Close()
			if err != nil {
				panic(err)
			}
		}(gamepad)
		gamepad_ = gamepad
	}

	mouse_, err := mouse.Create()
	if err != nil {
		panic(err)
	}

	keyboard_, err := keyboard.Create()
	if err != nil {
		panic(err)
	}

	rtc := r.New(config, mouse_, keyboard_, gamepad_)
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
