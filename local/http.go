package local

import (
	"encoding/json"
	"fmt"
	o "github.com/openstadia/openstadia/offer"
	"github.com/openstadia/openstadia/rtc"
	s "github.com/openstadia/openstadia/store"
	"io"
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

func (l *Local) rtcOfferHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(404)
		return
	}

	if l.rtc.IsBusy() {
		w.WriteHeader(404)
		return
	}

	body, err := io.ReadAll(r.Body)
	offer := o.Offer{}

	err = json.Unmarshal(body, &offer)
	if err != nil {
		panic(err)
	}

	answer := l.rtc.Offer(offer)

	marshal, err := json.Marshal(answer)
	if err != nil {
		panic(err)
	}

	_, err = fmt.Fprintf(w, string(marshal))
	if err != nil {
		return
	}
}

func (l *Local) getApps(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(404)
		return
	}

	apps := make([]string, 0)
	for _, app := range l.store.Apps() {
		apps = append(apps, app.Name)
	}

	marshal, err := json.Marshal(apps)
	if err != nil {
		panic(err)
	}

	_, err = fmt.Fprintf(w, string(marshal))
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
	http.HandleFunc("/api/apps", l.getApps)
	http.HandleFunc("/api/config", l.getConfig)
	http.HandleFunc("/api/hub", l.handleHub)

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	local := l.store.Local()
	addr := fmt.Sprintf("%s:%s", local.Host, local.Port)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatal(err)
	}
}
