package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/graphql-go/graphql"
	"github.com/lucapette/deluminator/app"
)

type dataSource struct {
	Name string `json:"name"`
}

var schema graphql.Schema

func GraphQLHandler(app *app.App) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/json")

		query, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		}

		res := graphql.Do(graphql.Params{
			Schema:        schema,
			RequestString: string(query),
			Context:       context.WithValue(context.Background(), "app", app),
		})

		if res.HasErrors() {
			w.WriteHeader(http.StatusBadRequest)
		}

		rJSON, err := json.Marshal(res)
		if err != nil {
			w.Write([]byte(err.Error()))
		}

		w.Write(rJSON)
	}
}

func init() {
	dataSourceType := graphql.NewObject(graphql.ObjectConfig{
		Name:        "DataSource",
		Description: fmt.Sprintf("A DataSource represents a single source of data to analyze"),
		Fields: graphql.Fields{
			"name": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
			},
		},
	})

	fields := graphql.Fields{
		"DataSources": &graphql.Field{
			Type: graphql.NewList(dataSourceType),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				app := p.Context.Value("app").(*app.App)
				ds, err := getLoaders(app)
				if err != nil {
					return nil, err
				}
				return ds, nil
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

func getLoaders(app *app.App) (ds []*dataSource, err error) {
	for _, loader := range app.Sources() {
		name := loader.DSN().DBName
		log.WithField("schema_name", name).Info("query metadata")

		start := time.Now()

		tables, err := loader.Tables()
		if err != nil {
			return ds, err
		}

		log.WithFields(log.Fields{
			"schema_name": name,
			"n_tables":    len(tables),
			"spent":       time.Now().Sub(start),
		}).Info("tables loaded")

		ds = append(ds, &dataSource{Name: name})
	}

	return ds, nil
}
