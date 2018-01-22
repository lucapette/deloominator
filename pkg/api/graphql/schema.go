package graphql

import (
	"log"

	gql "github.com/graphql-go/graphql"
	"github.com/lucapette/deloominator/pkg/charts"
	"github.com/lucapette/deloominator/pkg/db"
	"github.com/lucapette/deloominator/pkg/db/storage"
)

type queryError struct {
	Message string `json:"message"`
}

type cell struct {
	Value string `json:"value"`
}

type row struct {
	Cells []cell `json:"cells"`
}

type column struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type variable struct {
	Name           string `json:"name"`
	Value          string `json:"value"`
	IsControllable bool   `json:"isControllable"`
}

type results struct {
	Columns   []column   `json:"columns"`
	Rows      []row      `json:"rows"`
	ChartName string     `json:"chartName"`
	Variables []variable `json:"variables"`
}

type settings struct {
	IsReadOnly bool `json:"isReadOnly"`
}

type table struct {
	Name string `json:"name"`
}

type dataSource struct {
	Name   string   `json:"name"`
	Tables []*table `json:"tables"`
}

type payload struct {
	Query         string                 `json:"query"`
	OperationName string                 `json:"operationName,omitempty"`
	Variables     map[string]interface{} `json:"variables,omitempty"`
}

func convertToChartTypes(columns []db.Column) (types charts.DataTypes) {
	types = make(charts.DataTypes, len(columns))

	for i, col := range columns {
		switch col.Type {
		case db.Number:
			types[i] = charts.Number
		case db.Text:
			types[i] = charts.Text
		case db.Time:
			types[i] = charts.Time
		}
	}

	return types
}

func mutationRoot(s *storage.Storage) *gql.Object {
	fields := gql.Fields{
		"question": &gql.Field{
			Type: questionType,
			Args: gql.FieldConfigArgument{
				"title": &gql.ArgumentConfig{
					Type: gql.NewNonNull(gql.String),
				},
				"description": &gql.ArgumentConfig{
					Type: gql.String,
				},
				"query": &gql.ArgumentConfig{
					Type: gql.NewNonNull(gql.String),
				},
				"dataSource": &gql.ArgumentConfig{
					Type: gql.NewNonNull(gql.String),
				},
				"variables": &gql.ArgumentConfig{
					Type: gql.NewList(inputVariableType),
				},
			},
			Resolve: saveQuestion(s),
		},
	}
	return gql.NewObject(gql.ObjectConfig{Name: "Mutation", Fields: fields})
}

func queryRoot(dataSources db.DataSources, s *storage.Storage) *gql.Object {
	fields := gql.Fields{
		"settings": &gql.Field{
			Type:    settingsType,
			Resolve: resolveSettings(s),
		},
		"dataSources": &gql.Field{
			Type:    gql.NewList(dataSourceType),
			Resolve: resolveDataSources(dataSources),
		},
		"question": &gql.Field{
			Type: questionType,
			Args: gql.FieldConfigArgument{
				"id": &gql.ArgumentConfig{
					Type: gql.NewNonNull(gql.ID),
				},
			},
			Resolve: resolveQuestion(dataSources, s),
		},
		"questions": &gql.Field{
			Type:    gql.NewList(questionType),
			Resolve: resolveQuestions(dataSources, s),
		},
		"query": &gql.Field{
			Type: queryResultType,
			Args: gql.FieldConfigArgument{
				"dataSource": &gql.ArgumentConfig{
					Type: gql.NewNonNull(gql.String),
				},
				"query": &gql.ArgumentConfig{
					Type: gql.NewNonNull(gql.String),
				},
				"variables": &gql.ArgumentConfig{
					Type: gql.NewList(inputVariableType),
				},
			},
			Resolve: resolveQuery(dataSources),
		},
	}

	return gql.NewObject(gql.ObjectConfig{Name: "Query", Fields: fields})
}

func createSchema(dataSources db.DataSources, storage *storage.Storage) (schema gql.Schema) {
	schemaConfig := gql.SchemaConfig{
		Query:    queryRoot(dataSources, storage),
		Mutation: mutationRoot(storage),
	}

	schema, err := gql.NewSchema(schemaConfig)
	if err != nil {
		log.Fatalf("could not create gql schema: %v", err)
	}
	return schema
}
