package mouse

import (
	"fmt"
	"github.com/kbinani/screenshot"
	"github.com/lxn/win"
	"unsafe"
)

const MaxCoordinate = 65535
const WheelDelta = 120

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

type MouseImpl struct {
}

func Create() (*MouseImpl, error) {
	return &MouseImpl{}, nil
}

func (m *MouseImpl) Scroll(x int, y int) {
	err := SendMouseInput(win.MOUSEINPUT{
		MouseData: uint32(-y),
		DwFlags:   win.MOUSEEVENTF_WHEEL,
	})
	if err != nil {
		return
	}

	err = SendMouseInput(win.MOUSEINPUT{
		MouseData: uint32(x),
		DwFlags:   win.MOUSEEVENTF_HWHEEL,
	})
	if err != nil {
		return
	}
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

func (m *MouseImpl) Move(x int, y int) {
	// TODO Display select
	rect := screenshot.GetDisplayBounds(0)
	width := rect.Dx()
	height := rect.Dy()

	xScaled := float32(x) / float32(width)
	yScaled := float32(y) / float32(height)

	m.MoveFloat(xScaled, yScaled)
}

func (m *MouseImpl) MoveFloat(x float32, y float32) {
	xScaled := int32(x * MaxCoordinate)
	yScaled := int32(y * MaxCoordinate)

	err := SendMouseInput(win.MOUSEINPUT{
		Dx:      xScaled,
		Dy:      yScaled,
		DwFlags: win.MOUSEEVENTF_ABSOLUTE | win.MOUSEEVENTF_MOVE,
	})
	if err != nil {
		panic(err)
	}
}

func (m *MouseImpl) Update() {
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
