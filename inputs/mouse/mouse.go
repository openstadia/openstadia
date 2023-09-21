package mouse

type Button uint8

const (
	Left Button = iota
	Center
	Right
)

type Mouse interface {
	Move(x int, y int)
	MoveFloat(x float32, y float32)

	Scroll(x int, y int)

	MouseDown(button Button)
	MouseUp(button Button)

	Update()
}
