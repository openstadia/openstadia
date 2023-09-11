package keyboard

type Keyboard interface {
}

type KeyboardImpl struct {
}

func Create() (Keyboard, error) {
	return &KeyboardImpl{}, nil
}
