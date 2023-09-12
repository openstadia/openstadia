package rtc

import (
	"encoding/binary"
	"fmt"
	"github.com/openstadia/openstadia/application"
	c "github.com/openstadia/openstadia/config"
	"github.com/openstadia/openstadia/display"
	_ "github.com/openstadia/openstadia/driver/screen"
	"github.com/openstadia/openstadia/inputs/gamepad"
	"github.com/openstadia/openstadia/inputs/keyboard"
	"github.com/openstadia/openstadia/inputs/mouse"
	o "github.com/openstadia/openstadia/offer"
	"github.com/pion/mediadevices"
	"github.com/pion/mediadevices/pkg/frame"
	"github.com/pion/mediadevices/pkg/prop"
	"github.com/pion/webrtc/v3"
	"time"
)

type Rtc struct {
	config   *c.Openstadia
	track    *mediadevices.VideoTrack
	mouse    mouse.Mouse
	keyboard keyboard.Keyboard
	gamepad  gamepad.Gamepad
}

func New(config *c.Openstadia, mouse mouse.Mouse, keyboard keyboard.Keyboard, gamepad gamepad.Gamepad) *Rtc {
	return &Rtc{config: config, mouse: mouse, keyboard: keyboard, gamepad: gamepad}
}

func (r *Rtc) IsBusy() bool {
	return r.track != nil
}

func (r *Rtc) Offer(offer o.Offer) *webrtc.SessionDescription {
	webrtcConfig := webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: []string{"stun:stun.l.google.com:19302"},
			},
		},
	}

	codecParams := codecParams(offer)

	codecSelector := mediadevices.NewCodecSelector(
		mediadevices.WithVideoEncoders(codecParams),
	)

	mediaEngine := webrtc.MediaEngine{}
	codecSelector.Populate(&mediaEngine)
	api := webrtc.NewAPI(webrtc.WithMediaEngine(&mediaEngine))
	peerConnection, err := api.NewPeerConnection(webrtcConfig)
	if err != nil {
		panic(err)
	}

	name := offer.App.Name
	appConfig, err := r.config.GetAppByName(name)
	if err != nil {
		panic(err)
	}

	var app application.Application
	var display_ display.Display = nil

	if application.IsScreen(appConfig) {
		app = application.NewScreen()
	} else {
		display_ := display.Create(appConfig.Width, appConfig.Height)
		display_.Start()

		//TODO Add display creation check
		time.Sleep(time.Second * 5)

		env := display_.AppEnv()
		app = application.NewCmd(appConfig.Command[0], appConfig.Command[1:], env)
	}

	err = app.Start()
	if err != nil {
		panic(err)
	}

	peerConnection.OnConnectionStateChange(func(state webrtc.PeerConnectionState) {
		fmt.Printf("Connection State has changed %s \n", state.String())
		if state == webrtc.PeerConnectionStateDisconnected {
			if closeErr := peerConnection.Close(); closeErr != nil {
				fmt.Println(closeErr)
			}
		} else if state == webrtc.PeerConnectionStateClosed {
			closeErr := r.track.Close()
			if closeErr != nil {
				panic(closeErr)
			}

			if display_ != nil {
				display_.Stop()
			}

			app.Stop()
			r.track = nil
		}
	})

	// Set the handler for ICE connection state
	// This will notify you when the peer has connected/disconnected
	peerConnection.OnICEConnectionStateChange(func(connectionState webrtc.ICEConnectionState) {
		fmt.Printf("Connection State has changed %s \n", connectionState.String())
	})

	//markerEnable := false
	//marker := false

	//scale := video.Scale(640, 480, video.ScalerFastNearestNeighbor)

	peerConnection.OnDataChannel(func(d *webrtc.DataChannel) {
		fmt.Printf("New DataChannel %s %d\n", d.Label(), d.ID())

		// Register channel opening handling
		d.OnOpen(func() {
			fmt.Printf("Data channel '%s'-'%d' open\n", d.Label(), d.ID())
		})

		// Register text message handling
		d.OnMessage(func(msg webrtc.DataChannelMessage) {
			switch l := len(msg.Data); l {
			case 1:
				value := msg.Data[0] != 0
				fmt.Printf("Message from DataChannel '%s': '%t'\n", d.Label(), value)
				//marker = value
			case 12:
				event := int32(binary.LittleEndian.Uint32(msg.Data[:4]))
				x := int32(binary.LittleEndian.Uint32(msg.Data[4:8]))
				y := int32(binary.LittleEndian.Uint32(msg.Data[8:]))

				switch event {
				case 0:
					r.mouse.Move(int(x), int(y))
				case 1:
					r.mouse.Click()
				case 2:
					r.mouse.Scroll(int(x), int(y))
				}
			case 20:
				if r.gamepad != nil {
					parseGamepadData(r.gamepad, msg.Data)
				}
			}
		})
	})

	s, err := mediadevices.GetDisplayMedia(mediadevices.MediaStreamConstraints{
		Video: func(c *mediadevices.MediaTrackConstraints) {
			c.FrameFormat = prop.FrameFormat(frame.FormatRGBA)
			c.Width = prop.Int(640)
			c.Height = prop.Int(480)
			c.FrameRate = prop.Float(30)
		},
		Codec: codecSelector,
	})

	if err != nil {
		panic(err)
	}

	track := s.GetVideoTracks()[0]
	switch v := track.(type) {
	case *mediadevices.VideoTrack:
		//if markerEnable {
		//	mark := Mark(&marker)
		//	v.Transform(mark)
		//}

		r.track = v
	default:
		fmt.Printf("unexpected type %T\n", v)
	}

	track.OnEnded(func(err error) {
		fmt.Printf("Track (ID: %s) ended with error: %v\n",
			track.ID(), err)
	})

	_, err = peerConnection.AddTransceiverFromTrack(track,
		webrtc.RTPTransceiverInit{
			Direction: webrtc.RTPTransceiverDirectionSendonly,
		},
	)
	if err != nil {
		panic(err)
	}

	// Set the remote SessionDescription
	err = peerConnection.SetRemoteDescription(offer.SessionDescription)
	if err != nil {
		panic(err)
	}

	// Create an answer
	answer, err := peerConnection.CreateAnswer(nil)
	if err != nil {
		panic(err)
	}

	// Create channel that is blocked until ICE Gathering is complete
	gatherComplete := webrtc.GatheringCompletePromise(peerConnection)

	// Sets the LocalDescription, and starts our UDP listeners
	err = peerConnection.SetLocalDescription(answer)
	if err != nil {
		panic(err)
	}

	// Block until ICE Gathering is complete, disabling trickle ICE
	// we do this because we only can exchange one signaling message
	// in a production application you should exchange ICE Candidates via OnICECandidate
	<-gatherComplete

	return peerConnection.LocalDescription()
}
