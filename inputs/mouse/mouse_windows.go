package mouse

import (
	"fmt"
	"github.com/lxn/win"
	"unsafe"
)

var buttonDownMap = map[Button]uint32{
	Left:   win.MOUSEEVENTF_LEFTDOWN,
	Center: win.MOUSEEVENTF_MIDDLEDOWN,
	Right:  win.MOUSEEVENTF_RIGHTDOWN,
}

var buttonUpMap = map[Button]uint32{
	Left:   win.MOUSEEVENTF_LEFTUP,
	Center: win.MOUSEEVENTF_MIDDLEUP,
	Right:  win.MOUSEEVENTF_RIGHTUP,
}

func (m *MouseImpl) MouseDown(button Button) {
	key, ok := buttonDownMap[button]

	if !ok {
		return
	}

	err := SendMouseInput(win.MOUSEINPUT{
		DwFlags: key,
	})
	if err != nil {
		panic(err)
	}
}

func (m *MouseImpl) MouseUp(button Button) {
	key, ok := buttonUpMap[button]

	if !ok {
		return
	}

	err := SendMouseInput(win.MOUSEINPUT{
		DwFlags: key,
	})
	if err != nil {
		panic(err)
	}
}

func SendMouseInput(input win.MOUSEINPUT) error {
	m := win.MOUSE_INPUT{
		Type: win.INPUT_MOUSE,
		Mi:   input,
	}

	numSent := win.SendInput(1, unsafe.Pointer(&m), int32(unsafe.Sizeof(m)))
	if numSent != 1 {
		return fmt.Errorf("failed to send input, unknown error")
	}

	return nil
}
