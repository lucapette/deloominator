package api

import (
	"fmt"

	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/graphql-go/graphql"
	"github.com/lucapette/deloominator/pkg/charts"
	"github.com/lucapette/deloominator/pkg/db"
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

type results struct {
	Columns   []column `json:"columns"`
	Rows      []row    `json:"rows"`
	ChartName string   `json:"chartName"`
}

type Table struct {
	Name string `json:"name"`
}

type DataSource struct {
	Name   string   `json:"name"`
	Tables []*Table `json:"tables"`
}

type GraphqlPayload struct {
	Query         string                 `json:"query"`
	OperationName string                 `json:"operationName,omitempty"`
	Variables     map[string]interface{} `json:"variables,omitempty"`
}

func resolveDataSources(dbDataSources db.DataSources) func(p graphql.ResolveParams) (interface{}, error) {
	return func(p graphql.ResolveParams) (interface{}, error) {
		var dataSources []*DataSource

		for _, ds := range dbDataSources {
			log.WithField("schema_name", ds.Name()).Info("query metadata")

			start := time.Now()

			qr, err := ds.Tables()
			if err != nil {
				return dataSources, err
			}

			ts := make([]*Table, len(qr.Rows))
			for i, t := range qr.Rows {
				ts[i] = &Table{Name: t[0].Value}
			}

			log.WithFields(log.Fields{
				"schema_name": ds.Name(),
				"n_tables":    len(qr.Rows),
				"spent":       time.Now().Sub(start),
			}).Info("tables loaded")

			dataSources = append(dataSources, &DataSource{Name: ds.Name(), Tables: ts})
		}

		return dataSources, nil
	}
}

func resolveQuery(dataSources db.DataSources) func(p graphql.ResolveParams) (interface{}, error) {
	return func(p graphql.ResolveParams) (interface{}, error) {
		source := p.Args["source"].(string)
		input := p.Args["input"].(string)

		log.WithFields(log.Fields{
			"source": source,
			"input":  input,
		}).Infof("Query requested")

		qr, err := dataSources[source].Query(input)

		if err != nil {
			return queryError{Message: err.Error()}, nil
		}

		columns := make([]column, len(qr.Columns))

		for i, col := range qr.Columns {
			columns[i].Name = col.Name
			columns[i].Type = col.Type.String()
		}

		rows := make([]row, len(qr.Rows))

		for i, r := range qr.Rows {
			rows[i].Cells = make([]cell, len(qr.Columns))

			for j, c := range r {
				rows[i].Cells[j].Value = c.Value
			}
		}

		detectedChart := charts.Detect(convertToChartTypes(qr.Columns))

		return results{
			ChartName: detectedChart.String(),
			Columns:   columns,
			Rows:      rows,
		}, nil
	}
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

func queries(dataSources db.DataSources) *graphql.Object {
	queryErrorType := graphql.NewObject(graphql.ObjectConfig{
		Name:        "queryError",
		Description: "An error represents an error message from the data source",
		Fields: graphql.Fields{
			"message": &graphql.Field{
				Description: "Error message from the server",
				Type:        graphql.NewNonNull(graphql.String),
			},
		},
		IsTypeOf: func(p graphql.IsTypeOfParams) bool {
			_, ok := p.Value.(queryError)
			return ok
		},
	})

	cellType := graphql.NewObject(graphql.ObjectConfig{
		Name:        "Cell",
		Description: "A cell represents a single piece of returnted data",
		Fields: graphql.Fields{
			"value": &graphql.Field{
				Description: "Value of the cell",
				Type:        graphql.String,
			},
		},
	})
	rowType := graphql.NewObject(graphql.ObjectConfig{
		Name:        "Row",
		Description: "A row holds the representation of a set of cells of the raw data returned by the data source",
		Fields: graphql.Fields{
			"cells": &graphql.Field{
				Description: "Name of the column",
				Type:        graphql.NewList(cellType),
			},
		},
	})

	columnType := graphql.NewObject(graphql.ObjectConfig{
		Name:        "Column",
		Description: "A column holds the representation of one column of the raw data returned by a data source",
		Fields: graphql.Fields{
			"name": &graphql.Field{
				Description: "Name of the column",
				Type:        graphql.String,
			},
			"type": &graphql.Field{
				Description: "Type of the column",
				Type:        graphql.String,
			},
		},
	})

	resultsType := graphql.NewObject(graphql.ObjectConfig{
		Name:        "results",
		Description: "results represents a collection of raw data returned by a data source",
		Fields: graphql.Fields{
			"chartName": &graphql.Field{
				Description: "Detected chart name",
				Type:        graphql.String,
			},
			"columns": &graphql.Field{
				Description: "Columns of the returned results",
				Type:        graphql.NewList(columnType),
			},
			"rows": &graphql.Field{
				Description: "Rows of the returned results",
				Type:        graphql.NewList(rowType),
			},
		},
		IsTypeOf: func(p graphql.IsTypeOfParams) bool {
			_, ok := p.Value.(results)
			return ok
		},
	})

	queryResultType := graphql.NewUnion(graphql.UnionConfig{
		Name:        "QueryResult",
		Description: "QueryResult represents all the possible outcomes of a Query",
		Types:       []*graphql.Object{resultsType, queryErrorType},
		ResolveType: func(p graphql.ResolveTypeParams) *graphql.Object {
			if _, ok := p.Value.(results); ok {
				return resultsType
			}
			if _, ok := p.Value.(queryError); ok {
				return queryErrorType
			}
			return nil
		},
	})

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
		"dataSources": &graphql.Field{
			Type:    graphql.NewList(dataSourceType),
			Resolve: resolveDataSources(dataSources),
		},
		"query": &graphql.Field{
			Type: queryResultType,
			Args: graphql.FieldConfigArgument{
				"source": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"input": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
			Resolve: resolveQuery(dataSources),
		},
	}

	return graphql.NewObject(graphql.ObjectConfig{Name: "query", Fields: fields})
}

func createSchema(dataSources db.DataSources) (schema graphql.Schema) {
	queries := queries(dataSources)
	schemaConfig := graphql.SchemaConfig{Query: queries}

	schema, err := graphql.NewSchema(schemaConfig)
	if err != nil {
		log.Fatalf("failed to create new schema, error: %v", err)
	}
	return schema
}
