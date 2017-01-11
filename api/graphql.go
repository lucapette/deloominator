package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/graphql-go/graphql"
)

var schema graphql.Schema

func GraphQLHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/json")

	query, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	params := graphql.Params{
		Schema:        schema,
		RequestString: string(query),
	}
	res := graphql.Do(params)

	if res.HasErrors() {
		w.WriteHeader(http.StatusBadRequest)
	}

	rJSON, err := json.Marshal(res)
	if err != nil {
		w.Write([]byte(err.Error()))
	}

	w.Write(rJSON)
}

func init() {
	fields := graphql.Fields{
		"hello": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return "world", nil
			},
		},
		"ciao": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return "mondo", nil
			},
		},
	}

	rootQuery := graphql.ObjectConfig{Name: "Query", Fields: fields}
	schemaConfig := graphql.SchemaConfig{Query: graphql.NewObject(rootQuery)}

	var err error
	schema, err = graphql.NewSchema(schemaConfig)

	if err != nil {
		log.Fatalf("failed to create new schema, error: %v", err)
	}
}
