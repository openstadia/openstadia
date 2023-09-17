package mouse

import (
	"github.com/go-vgo/robotgo"
)

type Mouse interface {
	Move(x int, y int)
	MoveFloat(x float32, y float32)

	Scroll(x int, y int)

	MouseDown()
	MouseUp()
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

func (m *MouseImpl) MouseDown() {
	err := robotgo.MouseDown()
	if err != nil {
		panic(err)
	}
}

func (m *MouseImpl) MouseUp() {
	err := robotgo.MouseUp()
	if err != nil {
		panic(err)
	}
}

func (m *MouseImpl) Scroll(x int, y int) {
	robotgo.Scroll(x, y)
}

func Create() (*MouseImpl, error) {
	robotgo.MouseSleep = 0

	return &MouseImpl{}, nil
}
