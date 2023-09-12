package gamepad

import (
	"errors"
	"fmt"
	"github.com/openstadia/openstadia/inputs/gamepad/uinput"
	"os"
)

type vGamepad struct {
	name       []byte
	deviceFile *os.File
}

var buttonsMap = map[int]int{
	0: 0x130,
	1: 0x131,
	2: 0x133,
	3: 0x134,

	4: uinput.ButtonBumperLeft,
	5: uinput.ButtonBumperRight,
	6: uinput.ButtonTriggerLeft,
	7: uinput.ButtonTriggerRight,

	8: 0x13a,
	9: 0x13b,

	10: 0x13d,
	11: 0x13e,

	16: uinput.ButtonMode,
}

// CreateGamepad will create a new gamepad using the given uinput
// device path of the uinput device.
func CreateGamepad(path string, name []byte, vendor uint16, product uint16) (Gamepad, error) { // TODO: Consider moving this to a generic function that works for all devices
	err := uinput.ValidateDevicePath(path)
	if err != nil {
		return nil, err
	}
	err = uinput.ValidateUinputName(name)
	if err != nil {
		return nil, err
	}

	fd, err := createVGamepadDevice(path, name, vendor, product)
	if err != nil {
		return nil, err
	}

	return &vGamepad{name: name, deviceFile: fd}, nil
}

func (vg *vGamepad) ButtonPress(key int) error {
	key := buttonsMap[index]

	err := vg.ButtonDown(key)
	if err != nil {
		return err
	}
	err = vg.ButtonUp(key)
	if err != nil {
		return err
	}
	return nil
}

func (vg *vGamepad) ButtonDown(key int) error {
	return uinput.SendBtnEvent(vg.deviceFile, []int{key}, uinput.BtnStatePressed)
}

func (vg *vGamepad) ButtonUp(key int) error {
	return uinput.SendBtnEvent(vg.deviceFile, []int{key}, uinput.BtnStateReleased)
}

func (vg *vGamepad) LeftStickMoveX(value float32) error {
	return vg.sendStickAxisEvent(uinput.AbsX, value)
}

func (vg *vGamepad) LeftStickMoveY(value float32) error {
	return vg.sendStickAxisEvent(uinput.AbsY, value)
}

func (vg *vGamepad) RightStickMoveX(value float32) error {
	return vg.sendStickAxisEvent(uinput.AbsRX, value)
}

func (vg *vGamepad) RightStickMoveY(value float32) error {
	return vg.sendStickAxisEvent(uinput.AbsRY, value)
}

func (vg *vGamepad) RightStickMove(x, y float32) error {
	values := map[uint16]float32{}
	values[uinput.AbsRX] = x
	values[uinput.AbsRY] = y

	return vg.sendStickEvent(values)
}

func (vg *vGamepad) LeftStickMove(x, y float32) error {
	values := map[uint16]float32{}
	values[uinput.AbsX] = x
	values[uinput.AbsY] = y

	return vg.sendStickEvent(values)
}

func (vg *vGamepad) HatPress(direction HatDirection) error {
	return vg.sendHatEvent(direction, Press)
}

func (vg *vGamepad) HatRelease(direction HatDirection) error {
	return vg.sendHatEvent(direction, Release)
}

func (vg *vGamepad) sendStickAxisEvent(absCode uint16, value float32) error {
	ev := uinput.InputEvent{
		Type:  uinput.EvAbs,
		Code:  absCode,
		Value: denormalizeInput(value),
	}

	buf, err := uinput.InputEventToBuffer(ev)
	if err != nil {
		return fmt.Errorf("writing abs stick event failed: %v", err)
	}

	_, err = vg.deviceFile.Write(buf)
	if err != nil {
		return fmt.Errorf("failed to write abs stick event to device file: %v", err)
	}

	return uinput.SyncEvents(vg.deviceFile)
}

func (vg *vGamepad) sendStickEvent(values map[uint16]float32) error {
	for code, value := range values {
		ev := uinput.InputEvent{
			Type:  uinput.EvAbs,
			Code:  code,
			Value: denormalizeInput(value),
		}

		buf, err := uinput.InputEventToBuffer(ev)
		if err != nil {
			return fmt.Errorf("writing abs stick event failed: %v", err)
		}

		_, err = vg.deviceFile.Write(buf)
		if err != nil {
			return fmt.Errorf("failed to write abs stick event to device file: %v", err)
		}
	}

	return uinput.SyncEvents(vg.deviceFile)
}

func (vg *vGamepad) sendHatEvent(direction HatDirection, action HatAction) error {
	var event uint16
	var value int32

	switch direction {
	case HatUp:
		{
			event = uinput.AbsHat0Y
			value = -MaximumAxisValue
		}
	case HatDown:
		{
			event = uinput.AbsHat0Y
			value = MaximumAxisValue
		}
	case HatLeft:
		{
			event = uinput.AbsHat0X
			value = -MaximumAxisValue
		}
	case HatRight:
		{
			event = uinput.AbsHat0X
			value = MaximumAxisValue
		}
	default:
		{
			return errors.New("failed to parse input direction")
		}
	}

	if action == Release {
		value = 0
	}

	ev := uinput.InputEvent{
		Type:  uinput.EvAbs,
		Code:  event,
		Value: value,
	}

	buf, err := uinput.InputEventToBuffer(ev)
	if err != nil {
		return fmt.Errorf("writing abs stick event failed: %v", err)
	}

	_, err = vg.deviceFile.Write(buf)
	if err != nil {
		return fmt.Errorf("failed to write abs stick event to device file: %v", err)
	}

	return uinput.SyncEvents(vg.deviceFile)
}

func (vg *vGamepad) Close() error {
	return uinput.CloseDevice(vg.deviceFile)
}

func createVGamepadDevice(path string, name []byte, vendor uint16, product uint16) (fd *os.File, err error) {
	// This array is needed to register the event keys for the gamepad device.
	keys := []uint16{
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

	// absEvents is for the absolute events for the gamepad device.
	absEvents := []uint16{
		uinput.AbsX,
		uinput.AbsY,
		uinput.AbsZ,
		uinput.AbsRX,
		uinput.AbsRY,
		uinput.AbsRZ,
		uinput.AbsHat0X,
		uinput.AbsHat0Y,
	}

	deviceFile, err := uinput.CreateDeviceFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to create virtual gamepad device: %v", err)
	}

	// register button events
	err = uinput.RegisterDevice(deviceFile, uintptr(uinput.EvKey))
	if err != nil {
		_ = deviceFile.Close()
		return nil, fmt.Errorf("failed to register virtual gamepad device: %v", err)
	}

	for _, code := range keys {
		err = uinput.Ioctl(deviceFile, uinput.UiSetKeyBit, uintptr(code))
		if err != nil {
			_ = deviceFile.Close()
			return nil, fmt.Errorf("failed to register key number %d: %v", code, err)
		}
	}

	// register absolute events
	err = uinput.RegisterDevice(deviceFile, uintptr(uinput.EvAbs))
	if err != nil {
		_ = deviceFile.Close()
		return nil, fmt.Errorf("failed to register absolute event input device: %v", err)
	}

	for _, event := range absEvents {
		err = uinput.Ioctl(deviceFile, uinput.UiSetAbsBit, uintptr(event))
		if err != nil {
			_ = deviceFile.Close()
			return nil, fmt.Errorf("failed to register absolute event %v: %v", event, err)
		}
	}

	return uinput.CreateUsbDevice(deviceFile,
		uinput.UinputUserDev{
			Name: uinput.ToUinputName(name),
			ID: uinput.InputID{
				Bustype: uinput.BusUsb,
				Vendor:  vendor,
				Product: product,
				Version: 0x301}})
}

// Takes in a normalized value (-1.0:1.0) and return an event value
func denormalizeInput(value float32) int32 {
	return int32(value * MaximumAxisValue)
}
