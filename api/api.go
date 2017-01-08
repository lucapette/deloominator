package api

import (
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/lucapette/deluminator/db"
	"goji.io"
	"goji.io/pat"
)

type Config struct {
	Port    int
	Loaders db.Loaders
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

func homeHandler(w http.ResponseWriter, request *http.Request) {
	asset, err := Asset("assets/index.html")

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
	kind := pat.Param(r, "kind")
	name := pat.Param(r, "name")

	asset, err := Asset(strings.Join([]string{"assets", kind, name}, "/"))

	if err != nil {
		log.Println(err)
	}

	w.Header().Set("Content-Type", "text/"+kind)
	_, err = w.Write(asset)

	if err != nil {
		log.Println(err)
	}
}

func Start(config *Config) {
	router := goji.NewMux()

	router.Use(logHandler)

	router.HandleFunc(pat.Get("/"), homeHandler)
	router.HandleFunc(pat.New("/graphql"), GraphQLHandler)
	router.HandleFunc(pat.Get("/assets/:kind/:name"), assetsHandler)

	for _, loader := range config.Loaders {
		log.WithField("schema_name", loader.DSN().DBName).
			Info("query metadata")

		start := time.Now()

		tables, err := loader.Tables()
		if err != nil {
			log.Println("Cant start server:", err)
			os.Exit(1)
		}

		log.WithFields(log.Fields{
			"schema_name": loader.DSN().DBName,
			"n_tables":    len(tables),
			"spent":       time.Now().Sub(start),
		}).Info("tables loaded")
	}

	go func() {
		err := http.ListenAndServe(":"+strconv.Itoa(config.Port), router)
		if err != nil {
			log.Println("Cant start server:", err)
			os.Exit(1)
		}
	}()
}
