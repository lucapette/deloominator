package handlers

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/Sirupsen/logrus"
)

// Debug returns a debug handler
func Debug(inner http.Handler) http.Handler {
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
