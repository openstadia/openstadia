package display

type vDisplay struct {
}

func Create(width uint, height uint) Display {
	return &vDisplay{}
}

func (v *vDisplay) Start() {
}

func (v *vDisplay) Stop() {
}

func (v *vDisplay) AppEnv() []string {
	return nil
}
