package graphql

import (
	gql "github.com/graphql-go/graphql"
)

var settingsType = gql.NewObject(gql.ObjectConfig{
	Name:        "settings",
	Description: "Settings represents the current deloominator",
	Fields: gql.Fields{
		"isReadOnly": &gql.Field{
			Description: "Readonly server value",
			Type:        gql.Boolean,
		},
	},
})

var queryErrorType = gql.NewObject(gql.ObjectConfig{
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

var cellType = gql.NewObject(gql.ObjectConfig{
	Name:        "Cell",
	Description: "A cell represents a single piece of returnted data",
	Fields: gql.Fields{
		"value": &gql.Field{
			Description: "Value of the cell",
			Type:        gql.String,
		},
	},
})

var rowType = gql.NewObject(gql.ObjectConfig{
	Name:        "Row",
	Description: "A row holds the representation of a set of cells of the raw data returned by the data source",
	Fields: gql.Fields{
		"cells": &gql.Field{
			Description: "Name of the column",
			Type:        gql.NewList(cellType),
		},
	},
})

var columnType = gql.NewObject(gql.ObjectConfig{
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

var variableType = gql.NewObject(gql.ObjectConfig{
	Name:        "Variable",
	Description: "variable represents a single query variable",
	Fields: gql.Fields{
		"name": &gql.Field{
			Description: "variable name",
			Type:        gql.String,
		},
		"value": &gql.Field{
			Description: "variable value",
			Type:        gql.String,
		},
		"isControllable": &gql.Field{
			Description: "true if the variable has a controllable UI",
			Type:        gql.Boolean,
		},
	},
	IsTypeOf: func(p gql.IsTypeOfParams) bool {
		_, ok := p.Value.(variable)
		return ok
	},
})

var resultsType = gql.NewObject(gql.ObjectConfig{
	Name:        "results",
	Description: "results represents a collection of raw data returned by a data source",
	Fields: gql.Fields{
		"chartName": &gql.Field{
			Description: "Detected chart name",
			Type:        gql.String,
		},
		"variables": &gql.Field{
			Description: "query variables",
			Type:        gql.NewList(variableType),
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

var queryResultType = gql.NewUnion(gql.UnionConfig{
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

var tableType = gql.NewObject(gql.ObjectConfig{
	Name:        "Table",
	Description: "A table of a data source",
	Fields: gql.Fields{
		"name": &gql.Field{
			Type: gql.NewNonNull(gql.String),
		},
	},
})

var dataSourceType = gql.NewObject(gql.ObjectConfig{
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

var questionType = gql.NewObject(gql.ObjectConfig{
	Name: "Question",
	Fields: gql.Fields{
		"id": &gql.Field{
			Type: gql.String,
		},
		"title": &gql.Field{
			Type: gql.String,
		},
		"query": &gql.Field{
			Type: gql.String,
		},
		"createdAt": &gql.Field{
			Type: gql.String,
		},
		"updatedAt": &gql.Field{
			Type: gql.String,
		},
		"dataSource": &gql.Field{
			Type: gql.String,
		},
		"variables": &gql.Field{
			Type: gql.NewList(variableType),
		},
		"results": &gql.Field{
			Type: queryResultType,
		},
	},
})
