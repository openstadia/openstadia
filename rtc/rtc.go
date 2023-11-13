package rtc

import (
	"fmt"
	"github.com/openstadia/openstadia/application"
	"github.com/openstadia/openstadia/display"
	"github.com/openstadia/openstadia/inputs/gamepad"
	"github.com/openstadia/openstadia/inputs/keyboard"
	"github.com/openstadia/openstadia/inputs/mouse"
	o "github.com/openstadia/openstadia/offer"
	s "github.com/openstadia/openstadia/store"
	"github.com/pion/mediadevices"
	"github.com/pion/mediadevices/pkg/codec/opus"
	"github.com/pion/webrtc/v3"
	"time"
)

type Rtc struct {
	store    *s.Store
	tracks   []mediadevices.Track
	mouse    mouse.Mouse
	keyboard keyboard.Keyboard
	gamepad  gamepad.Gamepad
}

func New(store *s.Store, mouse mouse.Mouse, keyboard keyboard.Keyboard, gamepad gamepad.Gamepad) *Rtc {
	return &Rtc{store: store, mouse: mouse, keyboard: keyboard, gamepad: gamepad}
}

func (r *Rtc) IsBusy() bool {
	return r.tracks != nil
}

func (r *Rtc) Offer(offer o.Offer) *webrtc.SessionDescription {
	webrtcConfig := webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: []string{"stun:stun.l.google.com:19302"},
			},
		},
	}

	codecParams := getCodecParams(offer)

	// TODO Add check for audio requires
	opusParams, err := opus.NewParams()
	if err != nil {
		panic(err)
	}

	codecSelector := mediadevices.NewCodecSelector(
		mediadevices.WithVideoEncoders(codecParams),
		mediadevices.WithAudioEncoders(&opusParams),
	)

	mediaEngine := webrtc.MediaEngine{}
	codecSelector.Populate(&mediaEngine)
	api := webrtc.NewAPI(webrtc.WithMediaEngine(&mediaEngine))
	peerConnection, err := api.NewPeerConnection(webrtcConfig)
	if err != nil {
		panic(err)
	}

	appConfig, err := r.store.GetAppById(offer.AppId)

	if err != nil {
		panic(err)
	}

	var app application.Application
	var display_ display.Display = nil

	if application.IsScreen(appConfig) {
		app = application.NewScreen()
	} else {
		display_ = display.Create(appConfig.Width, appConfig.Height)
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
			for _, track := range r.tracks {
				closeErr := track.Close()
				if closeErr != nil {
					panic(closeErr)
				}
			}

			if display_ != nil {
				display_.Stop()
			}

			app.Stop()
			r.tracks = nil
		}
	})

	// Set the handler for ICE connection state
	// This will notify you when the peer has connected/disconnected
	peerConnection.OnICEConnectionStateChange(func(connectionState webrtc.ICEConnectionState) {
		fmt.Printf("Connection State has changed %s \n", connectionState.String())
	})

	//markerEnable := false
	//marker := false

	//scale := video.Scale(1080, 720, video.ScalerFastNearestNeighbor)

	peerConnection.OnDataChannel(func(d *webrtc.DataChannel) {
		fmt.Printf("New DataChannel %s %d\n", d.Label(), d.ID())

		// Register channel opening handling
		d.OnOpen(func() {
			fmt.Printf("Data channel '%s'-'%d' open\n", d.Label(), d.ID())
		})

		// Register text message handling
		d.OnMessage(func(msg webrtc.DataChannelMessage) {
			handleMessage(r, d, msg)
		})
	})

	s, err := app.GetMedia(codecSelector)

	if err != nil {
		panic(err)
	}

	//track := s.GetVideoTracks()[0]
	//switch v := track.(type) {
	//case *mediadevices.VideoTrack:
	//	//if markerEnable {
	//	//	mark := Mark(&marker)
	//	//	v.Transform(mark)
	//	//}
	//
	//	r.track = v
	//default:
	//	fmt.Printf("unexpected type %T\n", v)
	//}

	for _, track := range s.GetTracks() {
		fmt.Printf("Track %#v\n", track)
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

		r.tracks = append(r.tracks, track)
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
