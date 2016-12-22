package api

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func HomeHandler(w http.ResponseWriter, request *http.Request) {
	w.Write([]byte("it works!"))
}

func Start(port int) error {
	router := mux.NewRouter()

	router.HandleFunc("/", HomeHandler)

	err := http.ListenAndServe(":"+strconv.Itoa(port), router)

	return err
}
