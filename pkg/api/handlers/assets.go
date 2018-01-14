package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/gobuffalo/packr"
	"goji.io/pat"
)

// UI handles the requests for the UI application
func UI(port int) func(w http.ResponseWriter, r *http.Request) {
	box := packr.NewBox("../../../ui/dist")
	index := strings.Replace(box.String("index.html"), "%DELOOMINATOR_PORT%", strconv.Itoa(port), 1)

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")

		if _, err := w.Write([]byte(index)); err != nil {
			logrus.WithFields(logrus.Fields{
				"path":   r.URL.Path,
				"method": r.Method,
			}).Printf("cannot write asset: %v", err)
		}
	}
}

// Static handles the requests for the static assets
func Static() func(w http.ResponseWriter, r *http.Request) {
	box := packr.NewBox("../../../ui/dist")
	return func(w http.ResponseWriter, r *http.Request) {
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
}
