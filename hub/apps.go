package hub

import (
	"github.com/gorilla/websocket"
	"github.com/openstadia/openstadia/config"
	p "github.com/openstadia/openstadia/packet"
	"log"
)

type AppsAnswer struct {
	Apps []string `json:"apps"`
}

func handleApps(conn *websocket.Conn, header *p.Header, config *config.Openstadia) {
	apps := make([]string, 0)
	for _, app := range config.GetApps() {
		apps = append(apps, app.Name)
	}

	packetRes := p.Packet[AppsAnswer]{
		Header: p.Header{
			Type: p.TypeAck,
			Id:   header.Id,
		},
		Payload: AppsAnswer{Apps: apps},
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
