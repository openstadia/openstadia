package local

import (
	"encoding/json"
	"fmt"
	"github.com/openstadia/openstadia/rtc"
	s "github.com/openstadia/openstadia/store"
	"log"
	"net/http"
)

type Local struct {
	store *s.Store
	rtc   *rtc.Rtc
}

func New(store *s.Store, rtc *rtc.Rtc) *Local {
	return &Local{
		store: store,
		rtc:   rtc,
	}
}

func (l *Local) getConfig(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(404)
		return
	}

	config := l.store.Config()

	marshal, err := json.Marshal(config)
	if err != nil {
		panic(err)
	}

	_, err = fmt.Fprintf(w, string(marshal))
}

func (l *Local) ServeHttp() {
	http.HandleFunc("/api/offer", l.rtcOfferHandler)
	http.HandleFunc("/api/apps", l.handleApps)
	http.HandleFunc("/api/config", l.getConfig)
	http.HandleFunc("/api/hub", l.handleHub)

	fs := http.FileServer(http.Dir("./ui"))
	http.Handle("/assets/", fs)

	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./ui/favicon.ico")
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./ui/index.html")
	})

	local := l.store.Local()
	addr := fmt.Sprintf("%s:%s", local.Host, local.Port)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatal(err)
	}
}
