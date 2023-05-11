package main

import (
	"fmt"
	"github.com/go-vgo/robotgo"
	"github.com/gorilla/websocket"
	"github.com/openstadia/openstadia/packet"
	"github.com/openstadia/openstadia/types"
	"github.com/openstadia/openstadia/uinput"
	"github.com/pion/mediadevices"
	_ "github.com/pion/mediadevices/pkg/driver/screen"
	"github.com/pion/webrtc/v3"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

//sudo apt-get install libx11-dev libxext-dev libvpx-dev

//x11
//sudo apt install libx11-dev xorg-dev libxtst-dev

//Hook
//sudo apt install xcb libxcb-xkb-dev x11-xkb-utils libx11-xcb-dev libxkbcommon-x11-dev libxkbcommon-dev

//https://github.com/intel/media-driver

var pGamepad *uinput.Gamepad
var pTrack *mediadevices.VideoTrack

func ws(config *types.Openstadia, interrupt <-chan os.Signal) {
	u := url.URL{Scheme: "ws", Host: config.Hub, Path: "/ws"}
	log.Printf("connecting to %s", u.String())

	requestHeader := http.Header{}
	requestHeader.Add("Authorization", config.Token)

	c, _, err := websocket.DefaultDialer.Dial(u.String(), requestHeader)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Printf("recv: %s", message)

			packetReq := packet.Packet[webrtc.SessionDescription]{}
			packetReq.Decode(message)

			log.Printf("recv package: %#v", packetReq)

			answer := rtcOffer(config, packetReq.Data)

			packetRes := packet.Packet[webrtc.SessionDescription]{
				Type: packet.TypeAck,
				Data: *answer,
				Id:   packetReq.Id,
			}

			log.Printf("return package: %#v", packetRes)

			err = c.WriteMessage(websocket.TextMessage, packetRes.Encode())
			if err != nil {
				panic(err)
			}
		}
	}()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return
		case <-interrupt:
			log.Println("interrupt")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}

func main() {
	robotgo.MouseSleep = 0

	config, err := types.Load()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Config %#v\n", config)

	interrupt := make(chan os.Signal, 1)
	//signal.Notify(interrupt, os.Interrupt)

	go ws(config, interrupt)

	remoteGamepad := true
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

	select {}
}
