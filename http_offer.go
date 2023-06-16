package main

import (
	"encoding/json"
	"fmt"
	"github.com/openstadia/openstadia/types"
	"github.com/pion/webrtc/v3"
	"io"
	"log"
	"net/http"
)

func rtcOfferHandler(config *types.Openstadia, w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(404)
		return
	}

	if pTrack != nil {
		w.WriteHeader(404)
		return
	}

	body, err := io.ReadAll(r.Body)
	offer := webrtc.SessionDescription{}

	err = json.Unmarshal(body, &offer)
	if err != nil {
		panic(err)
	}

	answer := rtcOffer(config, offer)

	marshal, err := json.Marshal(answer)
	if err != nil {
		panic(err)
	}

	_, err = fmt.Fprintf(w, string(marshal))
	if err != nil {
		return
	}
}

func serveHttp(config *types.Openstadia) {
	http.HandleFunc("/rtcOffer", func(w http.ResponseWriter, r *http.Request) {
		rtcOfferHandler(config, w, r)
	})

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
