package gamepad

import (
	"errors"
	"github.com/openstadia/go-uinput"
)

type vGamepad struct {
	*uinput.UinputDevice
}

func CreateGamepad() (Gamepad, error) {
	path := "/dev/uinput"

	info := uinput.DeviceInfo{
		Name:    "Xbox One Wireless Controller",
		Vendor:  0x045E,
		Product: 0x02EA,
		Version: 0x0301,
	}

	device, err := createGamepadUinputDevice(path, info)
	if err != nil {
		return nil, err
	}

	return &vGamepad{device}, nil
}

func (vg *vGamepad) ButtonDown(key int) error {
	// TODO Fix DPAD rewrite opposite direction
	if isHatButton(key) {
		direction := hatMap[key]
		return vg.sendHatEvent(direction, Press)
	}

	button, ok := buttonsMap[key]
	if !ok {
		return nil
	}

	return vg.SendKeyEvent(uint16(button), uinput.BtnStatePressed)
}

func (vg *vGamepad) ButtonUp(key int) error {
	if isHatButton(key) {
		direction := hatMap[key]
		return vg.sendHatEvent(direction, Release)
	}

	button, ok := buttonsMap[key]
	if !ok {
		return nil
	}

	return vg.SendKeyEvent(uint16(button), uinput.BtnStateReleased)
}

func (vg *vGamepad) LeftStick(x float32, y float32) error {
	values := map[uint16]float32{}
	values[uinput.AbsX] = x
	values[uinput.AbsY] = y

	return vg.sendStickEvent(values)
}

func (vg *vGamepad) RightStick(x float32, y float32) error {
	values := map[uint16]float32{}
	values[uinput.AbsRX] = x
	values[uinput.AbsRY] = y

	return vg.sendStickEvent(values)
}

func (vg *vGamepad) LeftTrigger(value float32) error {
	return vg.sendTriggerEvent(uinput.AbsZ, value)
}

func (vg *vGamepad) RightTrigger(value float32) error {
	return vg.sendTriggerEvent(uinput.AbsRZ, value)
}

func (vg *vGamepad) Update() {
	err := vg.SendSyncEvent()
	if err != nil {
		panic(err)
	}
}

const MaximumAxisValue = 32767

// HatDirection specifies the direction of hat movement
type HatDirection int

const (
	HatUp HatDirection = iota + 1
	HatDown
	HatLeft
	HatRight
)

const (
	Press HatAction = iota + 1
	Release
)

var buttonsMap = map[int]int{
	0: uinput.ButtonSouth,
	1: uinput.ButtonEast,
	2: uinput.ButtonNorth,
	3: uinput.ButtonWest,

	4: uinput.ButtonBumperLeft,
	5: uinput.ButtonBumperRight,

	8: uinput.ButtonSelect,
	9: uinput.ButtonStart,

	10: uinput.ButtonThumbLeft,
	11: uinput.ButtonThumbRight,

	16: uinput.ButtonMode,
}

var hatMap = map[int]HatDirection{
	12: HatUp,
	13: HatDown,
	14: HatLeft,
	15: HatRight,
}

func (vg *vGamepad) sendStickAxisEvent(absCode uint16, value float32) error {
	valueScaled := denormalizeInput(value)
	err := vg.SendAbsEvent(absCode, valueScaled)

	return err
}

func (vg *vGamepad) sendStickEvent(values map[uint16]float32) error {
	for code, value := range values {
		err := vg.sendStickAxisEvent(code, value)
		if err != nil {
			return err
		}
	}

	return nil
}

func (vg *vGamepad) sendHatEvent(direction HatDirection, action HatAction) error {
	var event uint16
	var value int32

	switch direction {
	case HatUp:
		event = uinput.AbsHat0Y
		value = -MaximumAxisValue
	case HatDown:
		event = uinput.AbsHat0Y
		value = MaximumAxisValue
	case HatLeft:
		event = uinput.AbsHat0X
		value = -MaximumAxisValue
	case HatRight:
		event = uinput.AbsHat0X
		value = MaximumAxisValue
	default:
		{
			return errors.New("failed to parse input direction")
		}
	}

	if action == Release {
		value = 0
	}

	err := vg.SendAbsEvent(event, value)

	return err
}

func (vg *vGamepad) sendTriggerEvent(absCode uint16, value float32) error {
	valueScaled := scaleTrigger(value)
	err := vg.SendAbsEvent(absCode, valueScaled)

	return err
}

func (vg *vGamepad) Close() error {
	return vg.Close()
}

func createGamepadUinputDevice(path string, info uinput.DeviceInfo) (device *uinput.UinputDevice, err error) {
	keys := []uinput.KeyEvent{
		uinput.ButtonGamepad,

		uinput.ButtonSouth,
		uinput.ButtonEast,
		uinput.ButtonNorth,
		uinput.ButtonWest,

		uinput.ButtonBumperLeft,
		uinput.ButtonBumperRight,

		uinput.ButtonSelect,
		uinput.ButtonStart,

		uinput.ButtonMode,

		uinput.ButtonThumbLeft,
		uinput.ButtonThumbRight,
	}

	absEvents := []uinput.AbsEvent{
		uinput.AbsX,
		uinput.AbsY,
		uinput.AbsZ,
		uinput.AbsRX,
		uinput.AbsRY,
		uinput.AbsRZ,
		uinput.AbsHat0X,
		uinput.AbsHat0Y,
	}

	return uinput.CreateUinputDevice(path, info, keys, absEvents)
}

func isHatButton(button int) bool {
	_, ok := hatMap[button]
	return ok
}

func denormalizeInput(value float32) int32 {
	return int32(value * MaximumAxisValue)
}

func scaleTrigger(value float32) int32 {
	return int32((2*value - 1) * MaximumAxisValue)
}
