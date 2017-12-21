package graphql

import (
	"fmt"
	"strconv"
	"time"

	"github.com/Sirupsen/logrus"
	gql "github.com/graphql-go/graphql"
	"github.com/lucapette/deloominator/pkg/charts"
	"github.com/lucapette/deloominator/pkg/db"
	"github.com/lucapette/deloominator/pkg/db/storage"
)

func resolveDataSources(dbDataSources db.DataSources) func(p gql.ResolveParams) (interface{}, error) {
	return func(p gql.ResolveParams) (interface{}, error) {
		var dataSources []*DataSource

		for _, ds := range dbDataSources {
			logrus.WithField("schema_name", ds.DBName()).Info("query metadata")

			start := time.Now()

			names, err := ds.Tables()
			if err != nil {
				return dataSources, err
			}

			ts := make([]*Table, len(names))
			for i, name := range names {
				ts[i] = &Table{Name: name}
			}

			logrus.WithFields(logrus.Fields{
				"schema_name": ds.DBName(),
				"n_tables":    len(names),
				"spent":       time.Now().Sub(start),
			}).Info("tables loaded")

			dataSources = append(dataSources, &DataSource{Name: ds.DBName(), Tables: ts})
		}

		return dataSources, nil
	}
}

func resolveQuestion(s *storage.Storage) func(p gql.ResolveParams) (interface{}, error) {
	return func(p gql.ResolveParams) (interface{}, error) {
		id, err := strconv.Atoi(p.Args["id"].(string))
		if err != nil {
			return nil, err
		}

		return s.FindQuestion(id)
	}
}

func resolveQuestions(s *storage.Storage) func(p gql.ResolveParams) (interface{}, error) {
	return func(p gql.ResolveParams) (interface{}, error) {
		return s.AllQuestions()
	}
}

func resolveQuery(dataSources db.DataSources) func(p gql.ResolveParams) (interface{}, error) {
	return func(p gql.ResolveParams) (interface{}, error) {
		source := p.Args["source"].(string)
		query := p.Args["query"].(string)
		variables, ok := p.Args["variables"].(string)
		if !ok {
			variables = ""
		}

		logrus.WithFields(logrus.Fields{
			"source":    source,
			"query":     query,
			"variables": variables,
		}).Infof("Query requested")

		qr, err := dataSources[source].Query(db.Input{Query: query})

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
			Variables: variables,
		}, nil
	}
}

func resolveSettings(s *storage.Storage) func(p gql.ResolveParams) (interface{}, error) {
	return func(p gql.ResolveParams) (interface{}, error) {
		return settings{
			IsReadOnly: s == nil,
		}, nil
	}
}

func query(dataSources db.DataSources, s *storage.Storage) *gql.Object {
	settingsType := gql.NewObject(gql.ObjectConfig{
		Name:        "settings",
		Description: "Settings represents the current deloominator",
		Fields: gql.Fields{
			"isReadOnly": &gql.Field{
				Description: "Readonly server value",
				Type:        gql.Boolean,
			},
		},
	})

	queryErrorType := gql.NewObject(gql.ObjectConfig{
		Name:        "queryError",
		Description: "An error represents an error message from the data source",
		Fields: gql.Fields{
			"message": &gql.Field{
				Description: "Error message from the server",
				Type:        gql.NewNonNull(gql.String),
			},
		},
		IsTypeOf: func(p gql.IsTypeOfParams) bool {
			_, ok := p.Value.(queryError)
			return ok
		},
	})

	cellType := gql.NewObject(gql.ObjectConfig{
		Name:        "Cell",
		Description: "A cell represents a single piece of returnted data",
		Fields: gql.Fields{
			"value": &gql.Field{
				Description: "Value of the cell",
				Type:        gql.String,
			},
		},
	})
	rowType := gql.NewObject(gql.ObjectConfig{
		Name:        "Row",
		Description: "A row holds the representation of a set of cells of the raw data returned by the data source",
		Fields: gql.Fields{
			"cells": &gql.Field{
				Description: "Name of the column",
				Type:        gql.NewList(cellType),
			},
		},
	})

	columnType := gql.NewObject(gql.ObjectConfig{
		Name:        "Column",
		Description: "A column holds the representation of one column of the raw data returned by a data source",
		Fields: gql.Fields{
			"name": &gql.Field{
				Description: "Name of the column",
				Type:        gql.String,
			},
			"type": &gql.Field{
				Description: "Type of the column",
				Type:        gql.String,
			},
		},
	})

	resultsType := gql.NewObject(gql.ObjectConfig{
		Name:        "results",
		Description: "results represents a collection of raw data returned by a data source",
		Fields: gql.Fields{
			"chartName": &gql.Field{
				Description: "Detected chart name",
				Type:        gql.String,
			},
			"variables": &gql.Field{
				Description: "query variables",
				Type:        gql.String,
			},
			"columns": &gql.Field{
				Description: "Columns of the returned results",
				Type:        gql.NewList(columnType),
			},
			"rows": &gql.Field{
				Description: "Rows of the returned results",
				Type:        gql.NewList(rowType),
			},
		},
		IsTypeOf: func(p gql.IsTypeOfParams) bool {
			_, ok := p.Value.(results)
			return ok
		},
	})

	queryResultType := gql.NewUnion(gql.UnionConfig{
		Name:        "QueryResult",
		Description: "QueryResult represents all the possible outcomes of a Query",
		Types:       []*gql.Object{resultsType, queryErrorType},
		ResolveType: func(p gql.ResolveTypeParams) *gql.Object {
			if _, ok := p.Value.(results); ok {
				return resultsType
			}
			if _, ok := p.Value.(queryError); ok {
				return queryErrorType
			}
			return nil
		},
	})

	tableType := gql.NewObject(gql.ObjectConfig{
		Name:        "Table",
		Description: fmt.Sprintf("A table of a data source"),
		Fields: gql.Fields{
			"name": &gql.Field{
				Type: gql.NewNonNull(gql.String),
			},
		},
	})

	dataSourceType := gql.NewObject(gql.ObjectConfig{
		Name:        "DataSource",
		Description: "A DataSource represents a single source of data to analyze",
		Fields: gql.Fields{
			"name": &gql.Field{
				Type: gql.NewNonNull(gql.String),
			},
			"tables": &gql.Field{
				Type: gql.NewList(tableType),
			},
		},
	})

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
			Resolve: resolveQuestion(s),
		},
		"questions": &gql.Field{
			Type:    gql.NewList(questionType),
			Resolve: resolveQuestions(s),
		},
		"query": &gql.Field{
			Type: queryResultType,
			Args: gql.FieldConfigArgument{
				"source": &gql.ArgumentConfig{
					Type: gql.NewNonNull(gql.String),
				},
				"query": &gql.ArgumentConfig{
					Type: gql.NewNonNull(gql.String),
				},
				"variables": &gql.ArgumentConfig{
					Type: gql.String,
				},
			},
			Resolve: resolveQuery(dataSources),
		},
	}

	return gql.NewObject(gql.ObjectConfig{Name: "query", Fields: fields})
}
