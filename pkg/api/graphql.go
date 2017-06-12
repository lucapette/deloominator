package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/graphql-go/graphql"
	"github.com/lucapette/deloominator/pkg/db"
)

// GraphQLHandler is an HTTP handler for GraphQL queries
func GraphQLHandler(dataSources db.DataSources) func(w http.ResponseWriter, r *http.Request) {
	schema := createSchema(dataSources)

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/json")

		query, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		payload := GraphqlPayload{}

		err = json.Unmarshal(query, &payload)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		res := graphql.Do(graphql.Params{
			Schema:         schema,
			RequestString:  payload.Query,
			OperationName:  payload.OperationName,
			VariableValues: payload.Variables,
		})

		if res.HasErrors() {
			w.WriteHeader(http.StatusBadRequest)
		}

		rJSON, err := json.Marshal(res)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		_, err = w.Write(rJSON)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}
}
