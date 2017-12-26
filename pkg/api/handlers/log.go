package handlers

import (
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/felixge/httpsnoop"
)

// Log returns an handler that logs requests
func Log(inner http.Handler) http.Handler {
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
