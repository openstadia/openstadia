//go:build !tray

package runtime

import "github.com/openstadia/openstadia/config"

func Run(config *config.Openstadia) {
	select {}
}
