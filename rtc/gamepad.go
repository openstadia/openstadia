package rtc

import (
	"encoding/binary"
	"github.com/openstadia/openstadia/inputs/gamepad"
	"math"
)

func parseGamepadData(gamepad gamepad.Gamepad, data []byte) {
	var axes [4]float32
	for i := 0; i < 4; i++ {
		axes[i] = math.Float32frombits(binary.LittleEndian.Uint32(data[4*i:]))
	}

	var buttons [17]bool
	buttonsData := binary.LittleEndian.Uint32(data[24:])
	for i := 0; i < 17; i++ {
		buttons[i] = (buttonsData & (1 << i)) != 0
	}

	var triggers [2]float32
	for i := 0; i < 2; i++ {
		triggers[i] = math.Float32frombits(binary.LittleEndian.Uint32(data[16+4*i:]))
	}

	err := gamepad.LeftStick(axes[0], axes[1])
	if err != nil {
		panic(err)
	}

	err = gamepad.RightStick(axes[2], axes[3])
	if err != nil {
		panic(err)
	}

	for i := 0; i < 16; i++ {
		pressButton(gamepad, i, buttons)
	}

	err = gamepad.LeftTrigger(triggers[0])
	if err != nil {
		panic(err)
	}

	err = gamepad.RightTrigger(triggers[1])
	if err != nil {
		panic(err)
	}

	gamepad.Update()
}

func pressButton(gamepad gamepad.Gamepad, index int, buttons [17]bool) {
	if buttons[index] {
		err := gamepad.ButtonDown(index)
		if err != nil {
			panic(err)
		}
	} else {
		err := gamepad.ButtonUp(index)
		if err != nil {
			panic(err)
		}
	}
}
