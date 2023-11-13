package offer

import "github.com/pion/webrtc/v3"

type Offer struct {
	webrtc.SessionDescription
	AppId int   `json:"app_id"`
	Codec Codec `json:"codec"`
}
