package handlers

import (
	"fmt"
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/gobuffalo/packr"
	"goji.io/pat"
)

var box = packr.NewBox("../../../ui/dist")

// UI handles the requests for the UI application
func UI(w http.ResponseWriter, r *http.Request) {
	asset := box.Bytes("index.html")

	w.Header().Set("Content-Type", "text/html")

	if _, err := w.Write(asset); err != nil {
		logrus.WithFields(logrus.Fields{
			"path":   r.URL.Path,
			"method": r.Method,
		}).Printf("cannot write asset: %v", err)
	}
}

// Static handles the requests for the static assets
func Static(w http.ResponseWriter, r *http.Request) {
	name := pat.Param(r, "name")
	ext := pat.Param(r, "ext")

	asset := box.Bytes(fmt.Sprintf("%s.%s", name, ext))

	switch ext {
	case "js":
		w.Header().Set("Content-Type", "text/javascript")
	case "js.map":
		w.Header().Set("Content-Type", "application/json")
	}

	if _, err := w.Write(asset); err != nil {
		logrus.WithFields(logrus.Fields{
			"path":   r.URL.Path,
			"method": r.Method,
		}).Printf("cannot write asset: %v", err)
	}
}
