package main

import (
	"encoding/binary"
	"github.com/bendahl/uinput"
	"math"
)

var buttonsMap = map[int]int{
	0: 0x130,
	1: 0x131,
	2: 0x133,
	3: 0x134,

	4: uinput.ButtonBumperLeft,
	5: uinput.ButtonBumperRight,
	6: uinput.ButtonTriggerLeft,
	7: uinput.ButtonTriggerRight,

	8: 0x13a,
	9: 0x13b,

	10: 0x13d,
	11: 0x13e,

	16: uinput.ButtonMode,
}

var hatMap = map[int]uinput.HatDirection{
	12: uinput.HatUp,
	13: uinput.HatDown,
	14: uinput.HatLeft,
	15: uinput.HatRight,
}

func parseGamepadData(gamepad uinput.Gamepad, data []byte) {
	var axes [4]float32
	for i := 0; i < 4; i++ {
		axes[i] = math.Float32frombits(binary.LittleEndian.Uint32(data[4*i:]))
	}

	var buttons [17]bool
	buttonsData := binary.LittleEndian.Uint32(data[16:])
	for i := 0; i < 17; i++ {
		buttons[i] = (buttonsData & (1 << i)) != 0
	}

	err := gamepad.LeftStickMove(axes[0], axes[1])
	if err != nil {
		panic(err)
	}

	err = gamepad.RightStickMove(axes[2], axes[3])
	if err != nil {
		panic(err)
	}

	for i := 0; i < 12; i++ {
		pressButton(gamepad, i, buttons)
	}
	pressButton(gamepad, 16, buttons)

	pressHat(gamepad, 12, 13, buttons)
	pressHat(gamepad, 14, 15, buttons)
}

func pressHat(gamepad uinput.Gamepad, neg, pos int, buttons [17]bool) {
	if buttons[neg] {
		posHat := hatMap[neg]
		err := gamepad.HatPress(posHat)
		if err != nil {
			panic(err)
		}
	} else if buttons[pos] {
		negHat := hatMap[pos]
		err := gamepad.HatPress(negHat)
		if err != nil {
			panic(err)
		}
	} else {
		anyHat := hatMap[neg]
		err := gamepad.HatRelease(anyHat)
		if err != nil {
			panic(err)
		}
	}
}

func pressButton(gamepad uinput.Gamepad, index int, buttons [17]bool) {
	key := buttonsMap[index]
	if buttons[index] {
		err := gamepad.ButtonDown(key)
		if err != nil {
			panic(err)
		}
	} else {
		err := gamepad.ButtonUp(key)
		if err != nil {
			panic(err)
		}
	}
}
