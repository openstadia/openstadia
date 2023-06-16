package local

import (
	"encoding/json"
	"fmt"
	"github.com/openstadia/openstadia/config"
	o "github.com/openstadia/openstadia/offer"
	"github.com/openstadia/openstadia/rtc"
	"io"
	"log"
	"net/http"
)

type Local struct {
	config *config.Openstadia
	rtc    *rtc.Rtc
}

func New(config *config.Openstadia, rtc *rtc.Rtc) *Local {
	return &Local{
		config: config,
		rtc:    rtc,
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

func (l *Local) ServeHttp() {
	http.HandleFunc("/api/offer", l.rtcOfferHandler)

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	err := http.ListenAndServe(l.config.Local.Addr, nil)
	if err != nil {
		log.Fatal(err)
	}
}
