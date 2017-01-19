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

type table struct {
	Name string `json:"name"`
}

type dataSource struct {
	Name   string   `json:"name"`
	Tables []*table `json:"tables"`
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
	tableType := graphql.NewObject(graphql.ObjectConfig{
		Name:        "Table",
		Description: fmt.Sprintf("A table of a data source"),
		Fields: graphql.Fields{
			"name": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
			},
		},
	})

	dataSourceType := graphql.NewObject(graphql.ObjectConfig{
		Name:        "DataSource",
		Description: fmt.Sprintf("A DataSource represents a single source of data to analyze"),
		Fields: graphql.Fields{
			"name": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
			},
			"tables": &graphql.Field{
				Type: graphql.NewList(tableType),
			},
		},
	})

	fields := graphql.Fields{
		"DataSources": &graphql.Field{
			Type: graphql.NewList(dataSourceType),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				app := p.Context.Value("app").(*app.App)
				ds, err := getDataSources(app)
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

func getDataSources(app *app.App) (dataSources []*dataSource, err error) {
	for _, ds := range app.GetDataSources() {
		name := ds.DSN().DBName
		log.WithField("schema_name", name).Info("query metadata")

		start := time.Now()

		tables, err := ds.Tables()
		if err != nil {
			return dataSources, err
		}

		ts := make([]*table, len(tables))
		for i, t := range tables {
			ts[i] = &table{Name: t[0].Value}
		}

		log.WithFields(log.Fields{
			"schema_name": name,
			"n_tables":    len(tables),
			"spent":       time.Now().Sub(start),
		}).Info("tables loaded")

		dataSources = append(dataSources, &dataSource{Name: name, Tables: ts})
	}

	return dataSources, nil
}
