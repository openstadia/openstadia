package local

import (
	"encoding/json"
	"fmt"
	o "github.com/openstadia/openstadia/offer"
	"io"
	"net/http"
)

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

	fmt.Println(body)
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
