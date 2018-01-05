package graphql

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	gql "github.com/graphql-go/graphql"
	"github.com/lucapette/deloominator/pkg/db"
	"github.com/lucapette/deloominator/pkg/db/storage"
)

// Handler is an HTTP handler for GraphQL queries
func Handler(dataSources db.DataSources, storage *storage.Storage) func(w http.ResponseWriter, r *http.Request) {
	schema := createSchema(dataSources, storage)

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/json")

		query, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		payload := payload{}

		err = json.Unmarshal(query, &payload)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		res := gql.Do(gql.Params{
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
		}
	}
}
