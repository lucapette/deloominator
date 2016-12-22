package api

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
)

func homeHandler(w http.ResponseWriter, request *http.Request) {
	asset, err := Asset("assets/index.html")

	if err != nil {
		log.Println(err)
	}

	w.Write(asset)
}

func Start(port int) {
	router := mux.NewRouter()

	router.HandleFunc("/", homeHandler)

	go func() {
		err := http.ListenAndServe(":"+strconv.Itoa(port), router)
		if err != nil {
			log.Println("Cant start server:", err)
			os.Exit(1)
		}
	}()
}
