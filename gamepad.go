package main

import (
	"encoding/binary"
	"fmt"
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

	12: uinput.ButtonDpadUp,
	13: uinput.ButtonDpadDown,
	14: uinput.ButtonDpadLeft,
	15: uinput.ButtonDpadRight,

	16: uinput.ButtonMode,
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

	//fmt.Printf("%+v %+v", axes, buttons)

	err := gamepad.LeftStickMove(axes[0], axes[1])
	if err != nil {
		panic(err)
	}

	err = gamepad.RightStickMove(axes[2], axes[3])
	if err != nil {
		panic(err)
	}

	for i := 0; i < 17; i++ {
		key := buttonsMap[i]
		if buttons[i] {
			fmt.Printf("%d\n", i)
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
}
