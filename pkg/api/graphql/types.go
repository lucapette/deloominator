package graphql

import (
	gql "github.com/graphql-go/graphql"
)

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
		"dataSource": &gql.Field{
			Type: gql.String,
		},
	},
})
