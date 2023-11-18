//go:build !tray

package runtime

import s "github.com/openstadia/openstadia/store"

func Run(store s.Store) {
	select {}
}
