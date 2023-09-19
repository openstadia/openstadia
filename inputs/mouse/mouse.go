package mouse

import (
	"github.com/go-vgo/robotgo"
)

type Button uint8

const (
	Left Button = iota
	Center
	Right
)

type Mouse interface {
	Move(x int, y int)
	MoveFloat(x float32, y float32)

	Scroll(x int, y int)

	MouseDown(button Button)
	MouseUp(button Button)
}

type MouseImpl struct {
}

func (m *MouseImpl) MoveFloat(x float32, y float32) {
	width, height := robotgo.GetScreenSize()

	xScaled := int(x * float32(width))
	yScaled := int(y * float32(height))

	m.Move(xScaled, yScaled)
}

func (m *MouseImpl) Move(x int, y int) {
	robotgo.Move(x, y)
}

func (m *MouseImpl) Scroll(x int, y int) {
	robotgo.Scroll(x, y)
}

func Create() (*MouseImpl, error) {
	robotgo.MouseSleep = 0

	return &MouseImpl{}, nil
}
