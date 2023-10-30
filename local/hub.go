package local

import (
	"encoding/json"
	"fmt"
	"github.com/openstadia/openstadia/config"
	"io"
	"net/http"
)

func (l *Local) getHub(w http.ResponseWriter, r *http.Request) {
	hub := l.store.Hub()

	marshal, err := json.Marshal(hub)
	if err != nil {
		panic(err)
	}

	_, err = fmt.Fprintf(w, string(marshal))
}

func (l *Local) postHub(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	hub := config.Hub{}

	err = json.Unmarshal(body, &hub)
	if err != nil {
		panic(err)
	}

	err = l.store.SetHub(&hub)
	if err != nil {
		panic(err)
	}
}

func (l *Local) handleHub(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		l.getHub(w, r)
		return
	}

	if r.Method == "POST" {
		l.postHub(w, r)
		return
	}

	w.WriteHeader(404)
	return
}
