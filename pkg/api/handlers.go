package api

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
	"goji.io/pat"
)

func debugHandler(inner http.Handler) http.Handler {
	mw := func(w http.ResponseWriter, r *http.Request) {
		buf := bytes.NewBuffer(make([]byte, 0))
		reader := io.TeeReader(r.Body, buf)

		b, err := ioutil.ReadAll(reader)
		if err != nil {
			logrus.Fatalf("cannot read request: %v", err)
		}

		entry := logrus.WithFields(logrus.Fields{
			"method": r.Method,
			"body":   string(b),
		})

		for k, v := range r.Header {
			entry = entry.WithField(k, v)
		}

		entry.Info("incoming request")

		r.Body = ioutil.NopCloser(buf)
		inner.ServeHTTP(w, r)
	}
	return http.HandlerFunc(mw)
}

func logHandler(inner http.Handler) http.Handler {
	mw := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		proxyWriter := wrapWriter(w)

		inner.ServeHTTP(proxyWriter, r)

		logrus.WithFields(logrus.Fields{
			"spent":  time.Now().Sub(start),
			"path":   r.URL.Path,
			"method": r.Method,
			"status": proxyWriter.status(),
		}).Info("request completed")
	}
	return http.HandlerFunc(mw)
}

func uiHandler(w http.ResponseWriter, request *http.Request) {
	asset, err := Asset("ui/dist/index.html")

	if err != nil {
		logrus.Printf("cannot load index: %v", err)
	}

	w.Header().Set("Content-Type", "text/html")
	_, err = w.Write(asset)

	if err != nil {
		logrus.Printf("cannot write: %v", err)
	}
}

func assetsHandler(w http.ResponseWriter, r *http.Request) {
	name := pat.Param(r, "name")
	ext := pat.Param(r, "ext")

	asset, err := Asset(strings.Join([]string{"ui", "dist", fmt.Sprintf("%s.%s", name, ext)}, "/"))
	if err != nil {
		logrus.Println(err)
	}

	switch ext {
	case "js":
		w.Header().Set("Content-Type", "text/javascript")
	case "js.map":
		w.Header().Set("Content-Type", "application/json")
	}

	if _, err = w.Write(asset); err != nil {
		logrus.Printf("cannot write: %v", err)
	}
}
