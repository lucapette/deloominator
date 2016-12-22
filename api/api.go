package api

import (
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
)

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
	w.Write(asset)
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
	w.Write(asset)
}

func Start(port int) {
	router := mux.NewRouter()

	router.HandleFunc("/", LogHandler(homeHandler))
	router.HandleFunc("/assets/{kind}/{name}", LogHandler(assetsHandler))

	go func() {
		err := http.ListenAndServe(":"+strconv.Itoa(port), router)
		if err != nil {
			log.Println("Cant start server:", err)
			os.Exit(1)
		}
	}()
}
