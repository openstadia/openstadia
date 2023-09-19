package keyboard

type Keyboard interface {
	KeyDown(key string)
	KeyUp(key string)
}
