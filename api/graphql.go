package api

import (
	"encoding/json"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/graphql-go/graphql"
)

var schema graphql.Schema

func GraphQLHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query")

	params := graphql.Params{Schema: schema, RequestString: query}
	res := graphql.Do(params)

	if res.HasErrors() {
		log.Infof("failed to execute graphql operation, errors: %+v", res.Errors)
	}

	rJSON, err := json.Marshal(res)
	if err != nil {
		log.Info(err)
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
