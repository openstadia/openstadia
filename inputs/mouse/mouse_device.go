//go:build device

package mouse

import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
	gadget "github.com/openstadia/go-usb-gadget"
	o "github.com/openstadia/go-usb-gadget/option"
	"log"
	"time"
)

const MaxCoordinate = 32767

var buttonMap = map[Button]uint16{
	Left:   0,
	Right:  1,
	Center: 2,
}

type MouseImpl struct {
	buttons uint8
	x       uint16
	y       uint16
	wheel   int8

	writer *bufio.Writer
}

func Create() (*MouseImpl, error) {
	g := gadget.CreateGadget("my_hid")

	g.SetAttrs(&gadget.GadgetAttrs{
		BcdUSB:          o.Some[uint16](0x0200),
		BDeviceClass:    o.None[uint8](),
		BDeviceSubClass: o.None[uint8](),
		BDeviceProtocol: o.None[uint8](),
		BMaxPacketSize0: o.None[uint8](),
		IdVendor:        o.Some[uint16](0x1d6b),
		IdProduct:       o.Some[uint16](0x0104),
		BcdDevice:       o.Some[uint16](0x0100),
	})

	g.SetStrs(&gadget.GadgetStrs{
		SerialNumber: "fedcba9876543210",
		Manufacturer: "Tobias Girstmair",
		Product:      "iSticktoit.net USB Device",
	}, gadget.LangUsEng)

	c := gadget.CreateConfig(g, "c", 1)

	c.SetAttrs(&gadget.ConfigAttrs{
		BmAttributes: o.None[uint8](),
		BMaxPower:    o.Some[uint8](250),
	})

	c.SetStrs(&gadget.ConfigStrs{
		Configuration: "Config 1: ECM network",
	}, gadget.LangUsEng)

	hidFunction := gadget.CreateHidFunction(g, "usb0")
	hidFunction.SetAttrs(&gadget.HidFunctionAttrs{
		Subclass:     0,
		Protocol:     0,
		ReportLength: 6,
		ReportDesc:   ReportDesc,
	})

	b := gadget.CreateBinding(c, hidFunction, hidFunction.Name())
	fmt.Println(b)

	udcs := gadget.GetUdcs()
	if len(udcs) < 1 {
		return nil, errors.New("udc devices not found")
	}
	udc := udcs[0]

	fmt.Println(udc)

	g.Enable(udc)

	time.Sleep(time.Second)

	rw, _ := hidFunction.GetReadWriter()

	return &MouseImpl{
		buttons: 0,
		x:       0,
		y:       0,
		wheel:   0,
		writer:  rw.Writer,
	}, nil
}

func (m *MouseImpl) Move(x int, y int) {
	//TODO Fix this
	width := 1920
	height := 1080

	xScaled := float32(x) / float32(width)
	yScaled := float32(y) / float32(height)

	m.MoveFloat(xScaled, yScaled)
}

func (m *MouseImpl) MoveFloat(x float32, y float32) {
	xScaled := int32(x * MaxCoordinate)
	yScaled := int32(y * MaxCoordinate)

	m.x = uint16(xScaled)
	m.y = uint16(yScaled)
}

func (m *MouseImpl) Scroll(x int, y int) {
	m.wheel = int8(y)
}

func (m *MouseImpl) MouseDown(button Button) {
	key, ok := buttonMap[button]
	if !ok {
		return
	}

	m.buttons |= 1 << key
}

func (m *MouseImpl) MouseUp(button Button) {
	key, ok := buttonMap[button]
	if !ok {
		return
	}

	m.buttons &= ^(1 << key)
}

func (m *MouseImpl) Update() {
	report := make([]byte, 6)
	report[0] = m.buttons
	binary.LittleEndian.PutUint16(report[1:], m.x)
	binary.LittleEndian.PutUint16(report[3:], m.y)
	report[5] = byte(m.wheel)

	_, err := m.writer.Write(report)
	if err != nil {
		log.Fatal(err)
	}
	err = m.writer.Flush()
	if err != nil {
		log.Fatal(err)
	}
}
