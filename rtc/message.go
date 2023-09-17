package rtc

import (
	"encoding/binary"
	"github.com/pion/webrtc/v3"
	"math"
)

type ReportId byte

const (
	MouseMove   ReportId = 0
	MouseScroll ReportId = 2

	GamepadId ReportId = 3

	KeyboardDown ReportId = 4
	KeyboardUp   ReportId = 5

	MouseDown ReportId = 6
	MouseUp   ReportId = 7
)

func handleMessage(r *Rtc, d *webrtc.DataChannel, msg webrtc.DataChannelMessage) {
	reportId := ReportId(msg.Data[0])
	payload := msg.Data[4:]

	switch reportId {
	case MouseMove:
		x := math.Float32frombits(binary.LittleEndian.Uint32(payload[0:]))
		y := math.Float32frombits(binary.LittleEndian.Uint32(payload[4:]))
		r.mouse.MoveFloat(x, y)
	case MouseScroll:
		x := binary.LittleEndian.Uint32(payload[0:])
		y := binary.LittleEndian.Uint32(payload[4:])
		r.mouse.Scroll(int(x), int(y))
	case GamepadId:
		if r.gamepad != nil {
			parseGamepadData(r.gamepad, payload)
		}
	case KeyboardDown:
		r.keyboard.KeyDown(string(payload))
	case KeyboardUp:
		r.keyboard.KeyUp(string(payload))
	case MouseDown:
		r.mouse.MouseDown()
	case MouseUp:
		r.mouse.MouseUp()
	}
}
