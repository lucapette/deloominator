package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/lucapette/deloominator/pkg/query"
)

type queryPayload struct {
	Variables query.Variables `json:"variables"`
}

// QueryEvaluator analyzes queries typed in the UI application
func QueryEvaluator(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	q, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	payload := QueryPayload{}
	err = json.Unmarshal(q, &payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	variables := query.ExtractVariables(payload.Query)
	json, err := json.Marshal(queryPayload{Variables: variables})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = w.Write(json)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}
