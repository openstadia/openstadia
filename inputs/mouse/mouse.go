package mouse

import "github.com/go-vgo/robotgo"

type Mouse interface {
	Move(x int, y int)
	Click()
	Scroll(x int, y int)
}

type MouseImpl struct {
}

func (m *MouseImpl) Move(x int, y int) {
	robotgo.Move(x, y)
}

func (m *MouseImpl) Click() {
	robotgo.Click()
}

func (m *MouseImpl) Scroll(x int, y int) {
	robotgo.Scroll(x, y)
}

func Create() (Mouse, error) {
	robotgo.MouseSleep = 0

	return &MouseImpl{}, nil
}
