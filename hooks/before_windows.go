package hooks

import "golang.org/x/sys/windows"

func Before() {
	err := windows.SetDllDirectory("./lib")
	if err != nil {
		panic(err)
	}
}
