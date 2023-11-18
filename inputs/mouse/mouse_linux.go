//go:build !device

package mouse

import (
	"github.com/kbinani/screenshot"
	"github.com/openstadia/go-uinput"
)

const MaxCoordinate = 65535

var buttonMap = map[Button]uint16{
	Left:   uinput.BtnLeft,
	Center: uinput.BtnMiddle,
	Right:  uinput.BtnRight,
}

type MouseImpl struct {
	*uinput.Device
}

func Create() (*MouseImpl, error) {
	device, err := createMouseUinputDevice()
	if err != nil {
		return nil, err
	}

	return &MouseImpl{device}, nil
}

func (m *MouseImpl) Move(x int, y int) {
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

	err := m.SendAbsEvent(uinput.AbsX, xScaled)
	if err != nil {
		return
	}

	err = m.SendAbsEvent(uinput.AbsY, yScaled)
	if err != nil {
		return
	}
}

func (m *MouseImpl) Scroll(x int, y int) {
	err := m.SendRelEvent(uinput.RelHwheel, int32(x))
	if err != nil {
		return
	}

	err = m.SendRelEvent(uinput.RelWheel, int32(y))
	if err != nil {
		return
	}
}

func (m *MouseImpl) MouseDown(button Button) {
	key, ok := buttonMap[button]
	if !ok {
		return
	}

	err := m.SendKeyEvent(key, uinput.BtnStatePressed)
	if err != nil {
		return
	}
}

func (m *MouseImpl) MouseUp(button Button) {
	key, ok := buttonMap[button]
	if !ok {
		return
	}

	err := m.SendKeyEvent(key, uinput.BtnStateReleased)
	if err != nil {
		return
	}
}

func (m *MouseImpl) Update() {
	err := m.SendSyncEvent()
	if err != nil {
		return
	}
}

func createMouseUinputDevice() (device *uinput.Device, err error) {
	path := "/dev/uinput"

	info := uinput.DeviceInfo{
		Name:    "Basic Mouse",
		Vendor:  0x4711,
		Product: 0x0817,
		Version: 0x0001,
	}

	keyEvents := []uint16{
		uinput.BtnLeft,
		uinput.BtnRight,
		uinput.BtnMiddle,
	}

	absEvents := []uint16{
		uinput.AbsX,
		uinput.AbsY,
	}

	relEvents := []uint16{
		uinput.RelWheel,
		uinput.RelHwheel,
	}

	var absMax [64]int32
	var absMin [64]int32

	absMax[uinput.AbsX] = MaxCoordinate
	absMax[uinput.AbsY] = MaxCoordinate

	return uinput.CreateDevice(path, info, keyEvents, absEvents, relEvents, absMax, absMin)
}
