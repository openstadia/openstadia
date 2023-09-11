package display

type Display interface {
	Start()
	Stop()
	AppEnv() []string
}
