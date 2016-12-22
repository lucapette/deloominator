package api

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
)

func HomeHandler(w http.ResponseWriter, request *http.Request) {
	w.Write([]byte("it works!"))
}

func Start(port int) {
	router := mux.NewRouter()

	router.HandleFunc("/", HomeHandler)

	go func() {
		err := http.ListenAndServe(":"+strconv.Itoa(port), router)
		if err != nil {
			log.Println("Cant start server:", err)
			os.Exit(1)
		}
	}()
}
