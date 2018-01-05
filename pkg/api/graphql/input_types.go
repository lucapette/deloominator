package graphql

import "github.com/graphql-go/graphql"

var inputVariableType = graphql.NewInputObject(graphql.InputObjectConfig{
	Name:        "InputVariable",
	Description: "A query variable",

	Fields: graphql.InputObjectConfigFieldMap{
		"name": &graphql.InputObjectFieldConfig{
			Type: graphql.String,
		},
		"value": &graphql.InputObjectFieldConfig{
			Type: graphql.String,
		},
		"isControllable": &graphql.InputObjectFieldConfig{
			Type: graphql.Boolean,
		},
	},
})
