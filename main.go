package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-vgo/robotgo"
	"github.com/gorilla/websocket"
	"github.com/openstadia/openstadia/packet"
	"github.com/openstadia/openstadia/uinput"
	"github.com/pion/mediadevices"
	_ "github.com/pion/mediadevices/pkg/driver/screen"
	"github.com/pion/mediadevices/pkg/io/video"
	"github.com/pion/webrtc/v3"
	"golang.org/x/image/colornames"
	"image"
	"io"
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

func Mark(show *bool) video.TransformFunc {
	return func(r video.Reader) video.Reader {
		return video.ReaderFunc(func() (image.Image, func(), error) {
			for {
				img, _, err := r.Read()
				if err != nil {
					return nil, func() {}, err
				}

				switch v := img.(type) {
				case *image.RGBA:
					for yi := 0; yi < 16; yi++ {
						for xi := 0; xi < 16; xi++ {
							if *show {
								v.Set(xi, yi, colornames.Red)
							} else {
								v.Set(xi, yi, colornames.White)
							}
						}
					}
				default:
					fmt.Printf("unexpected type %T\n", v)
				}

				if *show {

				}

				return img, func() {}, nil
			}
		})
	}
}

func rtcOfferHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(404)
		return
	}

	if pTrack != nil {
		w.WriteHeader(404)
		return
	}

	body, err := io.ReadAll(r.Body)
	offer := webrtc.SessionDescription{}

	err = json.Unmarshal(body, &offer)
	if err != nil {
		panic(err)
	}

	answer := rtcOffer(offer)

	marshal, err := json.Marshal(answer)
	if err != nil {
		panic(err)
	}

	_, err = fmt.Fprintf(w, string(marshal))
	if err != nil {
		return
	}
}

func ws(interrupt <-chan os.Signal) {
	addr := "192.168.1.162:8001"

	u := url.URL{Scheme: "ws", Host: addr, Path: "/ws"}
	log.Printf("connecting to %s", u.String())

	requestHeader := http.Header{}
	requestHeader.Add("Authorization", "aaa")

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

			answer := rtcOffer(packetReq.Data)

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

	interrupt := make(chan os.Signal, 1)
	//signal.Notify(interrupt, os.Interrupt)

	go ws(interrupt)

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

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)
	http.HandleFunc("/rtcOffer", rtcOfferHandler)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
