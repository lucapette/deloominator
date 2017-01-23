package api

import "net/http"

// Heavily inspired by
// https://github.com/goji/glogrus2/blob/master/writer_proxy.go

// wrapWriter returns a proxy that wraps ResponseWriter
func wrapWriter(w http.ResponseWriter) *proxyWriter {
	bw := proxyWriter{ResponseWriter: w}
	return &bw
}

// proxyWriter holds the status code and a
// flag in addition to http.ResponseWriter
type proxyWriter struct {
	http.ResponseWriter
	wroteHeader bool
	code        int
}

// WriteHeader stores the status code and writes header
func (w *proxyWriter) WriteHeader(code int) {
	if !w.wroteHeader {
		w.code = code
		w.wroteHeader = true
		w.ResponseWriter.WriteHeader(code)
	}
}

// Write writes the bytes and calls MaybeWriteHeader
func (w *proxyWriter) Write(buf []byte) (int, error) {
	w.maybeWriteHeader()
	return w.ResponseWriter.Write(buf)
}

// maybeWriteHeader writes the header if it is not alredy set
func (w *proxyWriter) maybeWriteHeader() {
	if !w.wroteHeader {
		w.WriteHeader(http.StatusOK)
	}
}

// status returns the status
func (w *proxyWriter) status() int {
	return w.code
}
