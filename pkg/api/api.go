package api

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/lucapette/deloominator/pkg/app"
	"github.com/lucapette/deloominator/pkg/db"
	"github.com/rs/cors"
	"goji.io"
	"goji.io/pat"
)

func debugHandler(inner http.Handler) http.Handler {
	mw := func(w http.ResponseWriter, r *http.Request) {
		buf := bytes.NewBuffer(make([]byte, 0))
		reader := io.TeeReader(r.Body, buf)

		b, err := ioutil.ReadAll(reader)
		if err != nil {
			log.Fatal(err)
		}

		entry := log.WithFields(log.Fields{
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

		log.WithFields(log.Fields{
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
		log.Println(err)
	}

	w.Header().Set("Content-Type", "text/html")
	_, err = w.Write(asset)

	if err != nil {
		log.Println(err)
	}
}

func assetsHandler(w http.ResponseWriter, r *http.Request) {
	name := pat.Param(r, "name")
	ext := pat.Param(r, "ext")

	asset, err := Asset(strings.Join([]string{"ui", "dist", fmt.Sprintf("%s.%s", name, ext)}, "/"))

	if err != nil {
		log.Println(err)
	}

	switch ext {
	case "js":
		w.Header().Set("Content-Type", "text/javascript")
	case "js.map":
		w.Header().Set("Content-Type", "application/json")
	}

	_, err = w.Write(asset)

	if err != nil {
		log.Println(err)
	}
}

func Start(app *app.App, dataSources db.DataSources) {
	router := goji.NewMux()

	if app.Opts.Debug {
		router.Use(debugHandler)
	}

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:8080"},
	})

	router.Use(logHandler)
	router.Use(c.Handler)

	router.HandleFunc(pat.Post("/graphql"), GraphQLHandler(dataSources))
	router.HandleFunc(pat.Get("/:name.:ext"), assetsHandler)
	router.HandleFunc(pat.Get("/*"), uiHandler)

	go func() {
		err := http.ListenAndServe(":"+strconv.Itoa(app.Opts.Port), router)
		if err != nil {
			log.Fatal("Cant start server:", err)
		}
	}()
}
