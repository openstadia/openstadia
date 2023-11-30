package rtc

import (
	"fmt"
	"github.com/openstadia/openstadia/fs"
	o "github.com/openstadia/openstadia/offer"
	"github.com/openstadia/openstadia/packet"
	"github.com/pion/webrtc/v3"
	"os"
	"path/filepath"
)

func OfferFs(offer o.Offer) *webrtc.SessionDescription {
	config := webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: []string{"stun:stun.l.google.com:19302"},
			},
		},
	}

	peerConnection, err := webrtc.NewPeerConnection(config)
	if err != nil {
		panic(err)
	}

	peerConnection.OnConnectionStateChange(func(s webrtc.PeerConnectionState) {
		fmt.Printf("Peer Connection State has changed: %s\n", s.String())

		if s == webrtc.PeerConnectionStateFailed {
			fmt.Println("Peer Connection has gone to failed exiting")
		}

		if s == webrtc.PeerConnectionStateClosed {
			// PeerConnection was explicitly closed. This usually happens from a DTLS CloseNotify
			fmt.Println("Peer Connection has gone to closed exiting")
		}
	})

	peerConnection.OnDataChannel(func(d *webrtc.DataChannel) {
		fmt.Printf("New DataChannel %s %d\n", d.Label(), d.ID())

		d.OnOpen(func() {
			fmt.Printf("Data channel '%s'-'%d' open.\n", d.Label(), d.ID())
		})

		d.OnMessage(func(msg webrtc.DataChannelMessage) {
			fmt.Printf("Message from DataChannel '%s': '%s'\n", d.Label(), string(msg.Data))

			header := packet.Header{}
			err := header.DecodeFromPacket(msg.Data)
			if err != nil {
				return
			}

			if fs.Command(header.Name) == fs.Ls {
				pack := packet.Packet[fs.LsReq]{}
				err := pack.Decode(msg.Data)
				if err != nil {
					return
				}

				res, err := handleLs(&pack.Payload)
				if err != nil {
					return
				}

				ack := packet.MakeAck[fs.LsReq, *fs.LsRes](&pack, res)
				data, err := ack.Encode()
				fmt.Printf("Answer: %s\n", data)

				err = d.Send(data)
				if err != nil {
					return
				}
			}
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

func handleLs(lsReq *fs.LsReq) (*fs.LsRes, error) {
	path := filepath.Join(lsReq.Path...)
	rootedPath := filepath.Join(string(os.PathSeparator), path)

	dirs, err := os.ReadDir(rootedPath)
	if err != nil {
		return nil, err
	}

	files := make([]fs.File, 0)

	for _, dir := range dirs {
		info, err := dir.Info()
		if err != nil {
			continue
		}

		isHidden, err := fs.IsHidden(filepath.Join(rootedPath, info.Name()))
		if err != nil {
			continue
		}

		if isHidden {
			continue
		}

		file := fs.File{
			Name:     dir.Name(),
			Size:     info.Size(),
			IsDir:    dir.IsDir(),
			IsHidden: isHidden,
		}

		files = append(files, file)
	}

	lsRes := fs.LsRes{Files: files, Path: lsReq.Path}

	return &lsRes, nil
}
