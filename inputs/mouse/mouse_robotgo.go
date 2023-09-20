//go:build !windows

package mouse

import "github.com/go-vgo/robotgo"

var buttonMap = map[Button]string{
	Left:   "left",
	Center: "center",
	Right:  "right",
}

type MouseImpl struct {
}

func Create() (*MouseImpl, error) {
	robotgo.MouseSleep = 0

	return &MouseImpl{}, nil
}

func (m *MouseImpl) Scroll(x int, y int) {
	robotgo.Scroll(x, y)
}

func (m *MouseImpl) MouseDown(button Button) {
	key, ok := buttonMap[button]

	if !ok {
		return
	}

	err := robotgo.MouseDown(key)
	if err != nil {
		panic(err)
	}
}

func (m *MouseImpl) MouseUp(button Button) {
	key, ok := buttonMap[button]

	if !ok {
		return
	}

	err := robotgo.MouseUp(key)
	if err != nil {
		panic(err)
	}
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
