package uinput

import "io"

const MaximumAxisValue = 32767

// HatDirection specifies the direction of hat movement
type HatDirection int

const (
	HatUp HatDirection = iota + 1
	HatDown
	HatLeft
	HatRight
)

type HatAction int

const (
	Press HatAction = iota + 1
	Release
)

// Gamepad is a hybrid key / absolute change event output device.
// It used to enable a program to simulate gamepad input events.
type Gamepad interface {
	// ButtonPress will cause the button to be pressed and immediately released.
	ButtonPress(key int) error

	// ButtonDown will send a button-press event to an existing gamepad device.
	// The key can be any of the predefined keycodes from keycodes.go.
	// Note that the key will be "held down" until "KeyUp" is called.
	ButtonDown(key int) error

	// ButtonUp will send a button-release event to an existing gamepad device.
	// The key can be any of the predefined keycodes from keycodes.go.
	ButtonUp(key int) error

	// LeftStickMoveX performs a movement of the left stick along the x-axis
	LeftStickMoveX(value float32) error
	// LeftStickMoveY performs a movement of the left stick along the y-axis
	LeftStickMoveY(value float32) error

	// RightStickMoveX performs a movement of the right stick along the x-axis
	RightStickMoveX(value float32) error
	// RightStickMoveY performs a movement of the right stick along the y-axis
	RightStickMoveY(value float32) error

	// LeftStickMove moves the left stick along the x and y-axis
	LeftStickMove(x, y float32) error
	// RightStickMove moves the right stick along the x and y-axis
	RightStickMove(x, y float32) error

	// HatPress will issue a hat-press event in the given direction
	HatPress(direction HatDirection) error
	// HatRelease will issue a hat-release event in the given direction
	HatRelease(direction HatDirection) error

	io.Closer
}
