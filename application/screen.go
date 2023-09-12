package application

type screenApp struct {
}

func NewScreen() Application {
	return &screenApp{}
}

func (s *screenApp) Start() error {
	return nil
}

func (s *screenApp) Stop() {

}
