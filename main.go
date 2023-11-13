package main

import (
	"errors"
	"fmt"
	c "github.com/openstadia/openstadia/config"
	"github.com/openstadia/openstadia/hooks"
	h "github.com/openstadia/openstadia/hub"
	g "github.com/openstadia/openstadia/inputs/gamepad"
	"github.com/openstadia/openstadia/inputs/keyboard"
	"github.com/openstadia/openstadia/inputs/mouse"
	l "github.com/openstadia/openstadia/local"
	r "github.com/openstadia/openstadia/rtc"
	"github.com/openstadia/openstadia/runtime"
	s "github.com/openstadia/openstadia/store"
	"log"
	"os"
	"os/signal"
)

//Hook
//sudo apt install xcb libxcb-xkb-dev x11-xkb-utils libx11-xcb-dev libxkbcommon-x11-dev libxkbcommon-dev

//https://github.com/intel/media-driver

func main() {
	hooks.Before()

	config, err := c.Load()
	if err != nil {
		if errors.Is(err, c.ErrNoConfigFile) {
			fmt.Println("Config file not found")
		} else {
			log.Fatal(err)
		}
	}
	fmt.Printf("Config %#v\n", config)

	store, err := s.CreateStore()
	if err != nil {
		log.Fatal(err)
	}
	store.SetConfig(config)
	fmt.Printf("Store %#v\n", store.Config())

	remoteGamepad := true
	var gamepad_ g.Gamepad

	if remoteGamepad {
		gamepad, err := g.CreateGamepad()
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

	rtc := r.New(store, mouse_, keyboard_, gamepad_)
	local := l.New(store, rtc)
	hub := h.New(store, rtc)

	if hubObj := store.Hub(); hubObj != nil && hubObj.Enabled {
		interrupt := make(chan os.Signal, 1)
		signal.Notify(interrupt, os.Interrupt)

		go hub.Start(interrupt)
	}

	if localObj := store.Local(); localObj != nil && localObj.Enabled {
		go local.ServeHttp()
	}

	runtime.Run(store)
}
