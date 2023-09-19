//go:build !windows

package keyboard

import "github.com/go-vgo/robotgo"

type KeyboardImpl struct {
}

func Create() (*KeyboardImpl, error) {
	return &KeyboardImpl{}, nil
}

func (k *KeyboardImpl) KeyDown(key string) {
	err := robotgo.KeyDown(key)
	if err != nil {
		panic(err)
	}
}

func (k *KeyboardImpl) KeyUp(key string) {
	err := robotgo.KeyUp(key)
	if err != nil {
		panic(err)
	}
}
