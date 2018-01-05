package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/lucapette/deloominator/pkg/query"
)

type Variable struct {
	Name           string `json:"name"`
	Value          string `json:"value"`
	IsControllable bool   `json:"isControllable"`
}

type QueryEvaluatorPayload struct {
	Query     string     `json:"query"`
	Variables []Variable `json:"variables"`
}

// QueryEvaluator analyzes queries typed in the UI application
func QueryEvaluator(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	q, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	payload := QueryEvaluatorPayload{}
	err = json.Unmarshal(q, &payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	vars := make([]query.Variable, len(payload.Variables))

	for i, v := range payload.Variables {
		vars[i].Name = v.Name
		vars[i].Value = v.Value
	}

	evaledVars := query.NewEvaler(vars).ExtractVariables(payload.Query)
	variables := make([]Variable, len(evaledVars))
	for i, v := range evaledVars {
		variables[i].Name = v.Name
		variables[i].Value = v.Value
		variables[i].IsControllable = v.Name == "date" || v.Name == "timestamp"
	}
	json, err := json.Marshal(QueryEvaluatorPayload{Query: payload.Query, Variables: variables})
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
