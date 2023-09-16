package gamepad

import "io"

type HatAction int

// Gamepad is a hybrid key / absolute change event output device.
// It used to enable a program to simulate gamepad input events.
type Gamepad interface {
	ButtonDown(key int) error
	ButtonUp(key int) error

	LeftStick(x float32, y float32) error
	RightStick(x float32, y float32) error

	LeftTrigger(value float32) error
	RightTrigger(value float32) error

	Update()

	io.Closer
}
