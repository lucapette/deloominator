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

type csv struct {
	Data string `json:"data"`
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

type results struct {
	Columns   []column `json:"columns"`
	Rows      []row    `json:"rows"`
	ChartName string   `json:"chartName"`
	Variables string   `json:"variables"`
}

type settings struct {
	IsReadOnly bool `json:"isReadOnly"`
}

type Table struct {
	Name string `json:"name"`
}

type DataSource struct {
	Name   string   `json:"name"`
	Tables []*Table `json:"tables"`
}

type Payload struct {
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

func createSchema(dataSources db.DataSources, storage *storage.Storage) (schema gql.Schema) {
	schemaConfig := gql.SchemaConfig{
		Query:    query(dataSources, storage),
		Mutation: mutation(storage),
	}

	schema, err := gql.NewSchema(schemaConfig)
	if err != nil {
		log.Fatalf("could not create gql schema: %v", err)
	}
	return schema
}
