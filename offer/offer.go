package offer

import "github.com/pion/webrtc/v3"

type Offer struct {
	webrtc.SessionDescription
	Codec Codec
	App   App
}
