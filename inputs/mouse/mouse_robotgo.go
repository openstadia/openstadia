//go:build !windows

package mouse

import "github.com/go-vgo/robotgo"

var buttonMap = map[Button]string{
	Left:   "left",
	Center: "center",
	Right:  "right",
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
