package main

import (
	"encoding/binary"
	"fmt"
	"github.com/bendahl/uinput"
	"github.com/go-vgo/robotgo"
	"github.com/openstadia/openstadia/signal"
	"github.com/pion/mediadevices"
	"github.com/pion/mediadevices/pkg/codec/vpx"
	_ "github.com/pion/mediadevices/pkg/driver/screen"
	"github.com/pion/mediadevices/pkg/frame"
	"github.com/pion/mediadevices/pkg/io/video"
	"github.com/pion/mediadevices/pkg/prop"
	"github.com/pion/webrtc/v3"
	"golang.org/x/image/colornames"
	"image"
	"io"
	"log"
	"net/http"
)

//sudo apt-get install libx11-dev libxext-dev libvpx-dev

var pGamepad *uinput.Gamepad

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

func rtcOffer(w http.ResponseWriter, r *http.Request) {
	config := webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: []string{"stun:stun.l.google.com:19302"},
			},
		},
	}

	// Wait for the offer to be pasted
	offer := webrtc.SessionDescription{}
	body, err := io.ReadAll(r.Body)
	bodyString := string(body)
	signal.Decode(bodyString, &offer)

	vp8params, err := vpx.NewVP8Params()
	if err != nil {
		panic(err)
	}
	vp8params.BitRate = 10_000_000

	codecSelector := mediadevices.NewCodecSelector(
		mediadevices.WithVideoEncoders(&vp8params),
	)

	mediaEngine := webrtc.MediaEngine{}
	codecSelector.Populate(&mediaEngine)
	api := webrtc.NewAPI(webrtc.WithMediaEngine(&mediaEngine))
	peerConnection, err := api.NewPeerConnection(config)
	if err != nil {
		panic(err)
	}

	// Set the handler for ICE connection state
	// This will notify you when the peer has connected/disconnected
	peerConnection.OnICEConnectionStateChange(func(connectionState webrtc.ICEConnectionState) {
		fmt.Printf("Connection State has changed %s \n", connectionState.String())
	})

	marker := false
	mark := Mark(&marker)
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
				marker = value
			case 12:
				event := int32(binary.LittleEndian.Uint32(msg.Data[:4]))
				x := int32(binary.LittleEndian.Uint32(msg.Data[4:8]))
				y := int32(binary.LittleEndian.Uint32(msg.Data[8:]))

				switch event {
				case 0:
					robotgo.Move(int(x), int(y))
				case 1:
					robotgo.Click()
				case 2:
					robotgo.Scroll(-int(x), int(y))
				}
				fmt.Printf("x: %d, y: %d\n", x, y)
			case 20:
				if pGamepad != nil {
					parseGamepadData(*pGamepad, msg.Data)
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

	for _, track := range s.GetVideoTracks() {
		//videoTrack := track.(mediadevices.VideoTrack)
		fmt.Printf("Track (ID: %s) %+v %T\n", track.ID(), track, track)

		switch v := track.(type) {
		case *mediadevices.VideoTrack:
			v.Transform(mark)
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
	}

	// Set the remote SessionDescription
	err = peerConnection.SetRemoteDescription(offer)
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

	// Output the answer in base64 so we can paste it in browser
	_, err = fmt.Fprintf(w, signal.Encode(*peerConnection.LocalDescription()))
	if err != nil {
		return
	}
}

func main() {
	robotgo.MouseSleep = 0

	gamepad, err := uinput.CreateGamepad("/dev/uinput", []byte("testpad"), 0x045E, 0x02EA)
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

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)
	http.HandleFunc("/rtcOffer", rtcOffer)
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
