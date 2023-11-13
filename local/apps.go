package local

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (l *Local) getApps(w http.ResponseWriter, r *http.Request) {
	apps := l.store.Apps()
	marshal, err := json.Marshal(apps)
	if err != nil {
		panic(err)
	}

	_, err = fmt.Fprintf(w, string(marshal))
}

func (l *Local) handleApps(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		l.getApps(w, r)
		return
	}

	w.WriteHeader(404)
	return
}
