package rtc

import (
	"encoding/binary"
	"github.com/pion/webrtc/v3"
	"math"
)

type ReportId byte

const (
	MouseMove   ReportId = 0
	MouseClick  ReportId = 1
	MouseScroll ReportId = 2
	GamepadId   ReportId = 3
)

func handleMessage(r *Rtc, d *webrtc.DataChannel, msg webrtc.DataChannelMessage) {
	reportId := ReportId(msg.Data[0])
	payload := msg.Data[4:]

	switch reportId {
	case MouseMove:
		x := math.Float32frombits(binary.LittleEndian.Uint32(payload[0:]))
		y := math.Float32frombits(binary.LittleEndian.Uint32(payload[4:]))
		r.mouse.MoveFloat(x, y)
	case MouseClick:
		r.mouse.Click()
	case MouseScroll:
		x := binary.LittleEndian.Uint32(payload[0:])
		y := binary.LittleEndian.Uint32(payload[4:])
		r.mouse.Scroll(int(x), int(y))
	case GamepadId:
		if r.gamepad != nil {
			parseGamepadData(r.gamepad, payload)
		}
	}
}
