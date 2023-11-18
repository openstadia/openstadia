//go:build tray && windows

package runtime

import (
	"fmt"
	"github.com/getlantern/systray"
	s "github.com/openstadia/openstadia/store"
	"github.com/openstadia/openstadia/utils"
	"github.com/skratchdot/open-golang/open"
)

func Run(store s.Store) {
	systray.Run(func() {
		onReady(store)
	}, nil)
}

func onReady(store s.Store) {
	go func() {
		systray.SetTemplateIcon(Icon, Icon)
		systray.SetTitle("OpenStadia")
		systray.SetTooltip("OpenStadia")

		mUrl := systray.AddMenuItem("Open UI", "Open UI")
		mConfig := systray.AddMenuItem("Config Folder", "Open Config Folder")
		mQuit := systray.AddMenuItem("Quit", "Quit the whole app")

		for {
			select {
			case <-mUrl.ClickedCh:
				port := store.Local().Port
				open.Run(fmt.Sprintf("http://127.0.0.1:%s", port))
			case <-mConfig.ClickedCh:
				appConfigDir, _ := utils.GetConfigDir()
				open.Run(appConfigDir)
			case <-mQuit.ClickedCh:
				systray.Quit()
				fmt.Println("Quit now")
				return
			}
		}
	}()
}
