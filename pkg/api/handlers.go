package api

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/felixge/httpsnoop"
	"github.com/gobuffalo/packr"
	"goji.io/pat"
)

var box = packr.NewBox("../../ui/dist")

func debugHandler(inner http.Handler) http.Handler {
	mw := func(w http.ResponseWriter, r *http.Request) {
		buf := bytes.NewBuffer(make([]byte, 0))
		reader := io.TeeReader(r.Body, buf)

		b, err := ioutil.ReadAll(reader)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"path":   r.URL.Path,
				"method": r.Method,
			}).Fatalf("cannot read request: %v", err)
		}

		entry := logrus.WithFields(logrus.Fields{
			"method": r.Method,
			"body":   string(b),
		})

		for k, v := range r.Header {
			entry = entry.WithField(k, v)
		}

		entry.Print("incoming request")

		r.Body = ioutil.NopCloser(buf)
		inner.ServeHTTP(w, r)
	}
	return http.HandlerFunc(mw)
}

func logHandler(inner http.Handler) http.Handler {
	mw := func(w http.ResponseWriter, r *http.Request) {
		m := httpsnoop.CaptureMetrics(inner, w, r)

		logrus.WithFields(logrus.Fields{
			"duration": m.Duration,
			"path":     r.URL.Path,
			"method":   r.Method,
			"status":   m.Code,
			"written":  m.Written,
		}).Info("request completed")
	}
	return http.HandlerFunc(mw)
}

func uiHandler(w http.ResponseWriter, r *http.Request) {
	asset := box.Bytes("index.html")

	w.Header().Set("Content-Type", "text/html")

	if _, err := w.Write(asset); err != nil {
		logrus.WithFields(logrus.Fields{
			"path":   r.URL.Path,
			"method": r.Method,
		}).Printf("cannot write asset: %v", err)
	}
}

func assetsHandler(w http.ResponseWriter, r *http.Request) {
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
