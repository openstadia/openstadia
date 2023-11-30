package rtc

import (
	"fmt"
	o "github.com/openstadia/openstadia/offer"
	"github.com/openstadia/openstadia/terminal"
	"github.com/pion/webrtc/v3"
)

func OfferShell(offer o.Offer) *webrtc.SessionDescription {
	s := webrtc.SettingEngine{}
	s.DetachDataChannels()

	api := webrtc.NewAPI(webrtc.WithSettingEngine(s))

	config := webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: []string{"stun:stun.l.google.com:19302"},
			},
		},
	}

	peerConnection, err := api.NewPeerConnection(config)
	if err != nil {
		panic(err)
	}

	peerConnection.OnConnectionStateChange(func(s webrtc.PeerConnectionState) {
		fmt.Printf("Peer Connection State has changed: %s\n", s.String())

		if s == webrtc.PeerConnectionStateFailed {
			fmt.Println("Peer Connection has gone to failed exiting")
		}

		if s == webrtc.PeerConnectionStateClosed {
			fmt.Println("Peer Connection has gone to closed exiting")
		}
	})

	peerConnection.OnDataChannel(func(d *webrtc.DataChannel) {
		fmt.Printf("New DataChannel %s %d\n", d.Label(), d.ID())

		d.OnOpen(func() {
			fmt.Printf("Data channel '%s'-'%d' open.\n", d.Label(), d.ID())

			detach, err := d.Detach()
			if err != nil {
				return
			}

			pty, err := terminal.New(80, 30)
			if err != nil {
				return
			}

			pty.Start(detach)
		})
	})

	err = peerConnection.SetRemoteDescription(offer.SessionDescription)
	if err != nil {
		panic(err)
	}

	answer, err := peerConnection.CreateAnswer(nil)
	if err != nil {
		panic(err)
	}

	gatherComplete := webrtc.GatheringCompletePromise(peerConnection)

	err = peerConnection.SetLocalDescription(answer)
	if err != nil {
		panic(err)
	}

	<-gatherComplete

	return peerConnection.LocalDescription()
}
