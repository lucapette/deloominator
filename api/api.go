package api

import (
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"github.com/lucapette/deluminator/db"
)

type Config struct {
	Port    int
	Sources db.DataSources
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

func assetsHandler(w http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	kind := vars["kind"]
	name := vars["name"]

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
	router := mux.NewRouter()

	router.HandleFunc("/", LogHandler(homeHandler))
	router.HandleFunc("/assets/{kind}/{name}", LogHandler(assetsHandler))

	go func() {
		err := http.ListenAndServe(":"+strconv.Itoa(config.Port), router)
		if err != nil {
			log.Println("Cant start server:", err)
			os.Exit(1)
		}
	}()
}
