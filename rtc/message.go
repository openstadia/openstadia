package rtc

import (
	"encoding/binary"
	"github.com/openstadia/openstadia/inputs/mouse"
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
		r.mouse.Update()
	case MouseScroll:
		x := binary.LittleEndian.Uint32(payload[0:])
		y := binary.LittleEndian.Uint32(payload[4:])
		r.mouse.Scroll(int(x), int(y))
		r.mouse.Update()
	case GamepadId:
		if r.gamepad != nil {
			parseGamepadData(r.gamepad, payload)
		}
	case KeyboardDown:
		r.keyboard.KeyDown(string(payload))
		r.keyboard.Update()
	case KeyboardUp:
		r.keyboard.KeyUp(string(payload))
		r.keyboard.Update()
	case MouseDown:
		button := mouse.Button(payload[0])
		r.mouse.MouseDown(button)
		r.mouse.Update()
	case MouseUp:
		button := mouse.Button(payload[0])
		r.mouse.MouseUp(button)
		r.mouse.Update()
	}
}
