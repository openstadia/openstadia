package keyboard

import (
	"fmt"
	"github.com/lxn/win"
	"unsafe"
)

type KeyboardImpl struct {
}

func Create() (*KeyboardImpl, error) {
	return &KeyboardImpl{}, nil
}

func (k *KeyboardImpl) KeyDown(key string) {
	virtualKey, ok := codeVirtualKey[key]

	if !ok {
		return
	}

	err := SendKeydbInput(win.KEYBDINPUT{
		WVk:     virtualKey,
		DwFlags: 2,
	})
	if err != nil {
		panic(err)
	}
}

func (k *KeyboardImpl) KeyUp(key string) {
	virtualKey, ok := codeVirtualKey[key]

	if !ok {
		return
	}

	err := SendKeydbInput(win.KEYBDINPUT{
		WVk:     virtualKey,
		DwFlags: 2,
	})
	if err != nil {
		panic(err)
	}
}

func SendKeydbInput(input win.KEYBDINPUT) error {
	k := win.KEYBD_INPUT{
		Type: win.INPUT_KEYBOARD,
		Ki:   input,
	}

	numSent := win.SendInput(1, unsafe.Pointer(&k), int32(unsafe.Sizeof(k)))
	if numSent != 1 {
		return fmt.Errorf("failed to send input, unknown error")
	}

	return nil
}
