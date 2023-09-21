package gamepad

import (
	"github.com/openstadia/go-vigem"
	"math"
)
import X360 "github.com/openstadia/go-vigem/x360"

type vGamepad struct {
	client *vigem.ClientImpl
	x360   *X360.Gamepad
}

func CreateGamepad() (Gamepad, error) {
	client := vigem.NewClient()
	x360 := X360.NewGamepad(client)
	x360.Connect()
	return &vGamepad{client: client, x360: x360}, nil
}

var buttonsMap = map[int]uint16{
	0: X360.XUSB_GAMEPAD_A,
	1: X360.XUSB_GAMEPAD_B,
	2: X360.XUSB_GAMEPAD_X,
	3: X360.XUSB_GAMEPAD_Y,

	4: X360.XUSB_GAMEPAD_LEFT_SHOULDER,
	5: X360.XUSB_GAMEPAD_RIGHT_SHOULDER,
	6: 0,
	7: 0,

	8: X360.XUSB_GAMEPAD_BACK,
	9: X360.XUSB_GAMEPAD_START,

	10: X360.XUSB_GAMEPAD_LEFT_THUMB,
	11: X360.XUSB_GAMEPAD_RIGHT_THUMB,

	12: X360.XUSB_GAMEPAD_DPAD_UP,
	13: X360.XUSB_GAMEPAD_DPAD_DOWN,
	14: X360.XUSB_GAMEPAD_DPAD_LEFT,
	15: X360.XUSB_GAMEPAD_DPAD_RIGHT,

	16: X360.XUSB_GAMEPAD_GUIDE,
}

func (vg *vGamepad) ButtonDown(key int) error {
	button := buttonsMap[key]
	if button == 0 {
		return nil
	}

	vg.x360.PressButton(button)
	return nil
}

func (vg *vGamepad) ButtonUp(key int) error {
	button := buttonsMap[key]
	if button == 0 {
		return nil
	}

	vg.x360.ReleaseButton(button)
	return nil
}

func (vg *vGamepad) RightStick(x float32, y float32) error {
	xInt := scaleStick(x)
	yInt := -scaleStick(y)

	vg.x360.RightJoystick(xInt, yInt)
	return nil
}

func (vg *vGamepad) LeftStick(x float32, y float32) error {
	xScaled := scaleStick(x)
	yScaled := -scaleStick(y)

	vg.x360.LeftJoystick(xScaled, yScaled)
	return nil
}

func (vg *vGamepad) LeftTrigger(value float32) error {
	valueScaled := scaleTrigger(value)
	vg.x360.LeftTrigger(valueScaled)
	return nil
}

func (vg *vGamepad) RightTrigger(value float32) error {
	valueScaled := scaleTrigger(value)
	vg.x360.RightTrigger(valueScaled)
	return nil
}

func (vg *vGamepad) Update() {
	vg.x360.Update()
}

func (vg *vGamepad) Close() error {
	vg.x360.Disconnect()
	vg.x360.Release()

	vg.client.Release()
	return nil
}

func scaleStick(value float32) int16 {
	return int16(value * math.MaxInt16)
}

func scaleTrigger(value float32) uint8 {
	return uint8(value * math.MaxUint8)
}
