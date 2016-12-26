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
	Port        int
	DataSources db.DataSources
}

func LogHandler(exe func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	f := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		exe(w, r)
		log.WithFields(log.Fields{
			"spent":  time.Now().Sub(start),
			"path":   r.URL.Path,
			"method": r.Method,
		}).Info("request completed")
	}
	return f
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

	router.HandleFunc(pat.Get("/"), LogHandler(homeHandler))
	router.HandleFunc(pat.Get("/assets/:kind/:name"), LogHandler(assetsHandler))

	for _, ds := range config.DataSources {
		log.WithField("ds", ds.DBName).
			Info("query metadata")

		tables, err := ds.Tables()
		if err != nil {
			log.Println("Cant start server:", err)
			os.Exit(1)
		}

		for table := range tables {
			log.Info(table)
		}
	}

	go func() {
		err := http.ListenAndServe(":"+strconv.Itoa(config.Port), router)
		if err != nil {
			log.Println("Cant start server:", err)
			os.Exit(1)
		}
	}()
}
