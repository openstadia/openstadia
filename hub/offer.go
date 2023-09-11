package hub

import (
	"github.com/gorilla/websocket"
	"github.com/openstadia/openstadia/offer"
	p "github.com/openstadia/openstadia/packet"
	"github.com/openstadia/openstadia/rtc"
	"github.com/pion/webrtc/v3"
	"log"
)

func handleOffer(conn *websocket.Conn, message []byte, rtc *rtc.Rtc) {
	packet := p.Packet[offer.Offer]{}
	packet.Decode(message)

	answer := rtc.Offer(packet.Payload)

	packetRes := p.Packet[webrtc.SessionDescription]{
		Header: p.Header{
			Type: p.TypeAck,
			Id:   packet.Header.Id,
		},
		Payload: *answer,
	}

	log.Printf("return package: %#v", packetRes)

	encoded, err := packetRes.Encode()
	if err != nil {
		panic(err)
	}

	err = conn.WriteMessage(websocket.TextMessage, encoded)
	if err != nil {
		panic(err)
	}
}
