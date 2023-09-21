package keyboard

import (
	"github.com/openstadia/go-uinput"
)

const KeyMax = 248

type KeyboardImpl struct {
	*uinput.Device
}

func Create() (*KeyboardImpl, error) {
	device, err := createKeyboardUinputDevice()
	if err != nil {
		return nil, err
	}

	return &KeyboardImpl{device}, nil
}

func (k *KeyboardImpl) KeyDown(key string) {
	keyCode, ok := codeKey[key]
	if !ok {
		return
	}

	err := k.SendKeyEvent(keyCode, uinput.BtnStatePressed)
	if err != nil {
		return
	}
}

func (k *KeyboardImpl) KeyUp(key string) {
	keyCode, ok := codeKey[key]
	if !ok {
		return
	}

	err := k.SendKeyEvent(keyCode, uinput.BtnStateReleased)
	if err != nil {
		return
	}
}

func (k *KeyboardImpl) Update() {
	err := k.SendSyncEvent()
	if err != nil {
		return
	}
}

func createKeyboardUinputDevice() (device *uinput.Device, err error) {
	path := "/dev/uinput"

	info := uinput.DeviceInfo{
		Name:    "Basic Keyboard",
		Vendor:  0x4711,
		Product: 0x0815,
		Version: 0x0001,
	}

	var keyEvents = make([]uint16, KeyMax)
	for i := 0; i < KeyMax; i++ {
		keyEvents[i] = uint16(i)
	}

	var absMax [64]int32
	var absMin [64]int32

	return uinput.CreateDevice(path, info, keyEvents, nil, nil, absMax, absMin)
}
