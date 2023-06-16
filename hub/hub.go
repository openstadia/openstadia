package hub

import (
	"github.com/gorilla/websocket"
	"github.com/openstadia/openstadia/config"
	"github.com/openstadia/openstadia/offer"
	"github.com/openstadia/openstadia/packet"
	"github.com/openstadia/openstadia/rtc"
	"github.com/pion/webrtc/v3"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

type Hub struct {
	config *config.Openstadia
	rtc    *rtc.Rtc
}

func New(config *config.Openstadia, rtc *rtc.Rtc) *Hub {
	return &Hub{
		config: config,
		rtc:    rtc,
	}
}

func (h *Hub) Start(interrupt <-chan os.Signal) {
	u := url.URL{Scheme: "ws", Host: h.config.Hub.Addr, Path: "/ws"}
	log.Printf("connecting to %s", u.String())

	requestHeader := http.Header{}
	requestHeader.Add("Authorization", h.config.Hub.Token)

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

			packetReq := packet.Packet[offer.Offer]{}
			packetReq.Decode(message)

			log.Printf("recv package: %#v", packetReq)

			answer := h.rtc.Offer(packetReq.Data)

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
