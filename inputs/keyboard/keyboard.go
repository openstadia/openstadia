package keyboard

import "github.com/go-vgo/robotgo"

type Keyboard interface {
}

type KeyboardImpl struct {
}

func Create() (Keyboard, error) {
	return &KeyboardImpl{}, nil
}

func (k *KeyboardImpl) KeyDown(key string, args ...interface{}) {
	err := robotgo.KeyDown(key, args...)
	if err != nil {
		panic(err)
	}
}

func (k *KeyboardImpl) KeyUp(key string, args ...interface{}) {
	err := robotgo.KeyUp(key, args...)
	if err != nil {
		panic(err)
	}
}
