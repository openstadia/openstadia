package mouse

type dummy struct {
}

func (d *dummy) Move(x int, y int) {
}

func (d *dummy) MoveFloat(x float32, y float32) {
}

func (d *dummy) Scroll(x int, y int) {
}

func (d *dummy) MouseDown(button Button) {
}

func (d *dummy) MouseUp(button Button) {
}

func (d *dummy) Update() {
}

func CreateDummy() (Mouse, error) {
	return &dummy{}, nil
}
