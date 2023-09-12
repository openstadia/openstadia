package gamepad

type vGamepad struct {
	name []byte
}

// CreateGamepad will create a new gamepad using the given uinput
// device path of the uinput device.
func CreateGamepad(path string, name []byte, vendor uint16, product uint16) (Gamepad, error) {
	return &vGamepad{name: name}, nil
}

func (vg *vGamepad) ButtonPress(key int) error {
	return nil
}

func (vg *vGamepad) ButtonDown(key int) error {
	return nil
}

func (vg *vGamepad) ButtonUp(key int) error {
	return nil
}

func (vg *vGamepad) LeftStickMoveX(value float32) error {
	return nil
}

func (vg *vGamepad) LeftStickMoveY(value float32) error {
	return nil
}

func (vg *vGamepad) RightStickMoveX(value float32) error {
	return nil
}

func (vg *vGamepad) RightStickMoveY(value float32) error {
	return nil
}

func (vg *vGamepad) RightStickMove(x, y float32) error {
	return nil
}

func (vg *vGamepad) LeftStickMove(x, y float32) error {
	return nil
}

func (vg *vGamepad) HatPress(direction HatDirection) error {
	return nil
}

func (vg *vGamepad) HatRelease(direction HatDirection) error {
	return nil
}

func (vg *vGamepad) Close() error {
	return nil
}
