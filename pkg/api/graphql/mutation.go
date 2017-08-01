package graphql

import (
	gql "github.com/graphql-go/graphql"
	"github.com/lucapette/deloominator/pkg/db"
)

func saveQuestion(storage *db.Storage) func(p gql.ResolveParams) (interface{}, error) {
	return func(p gql.ResolveParams) (interface{}, error) {
		return nil, nil
	}
}

func mutation(storage *db.Storage) *gql.Object {
	questionType := gql.NewObject(gql.ObjectConfig{
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
		},
	})

	fields := gql.Fields{
		"saveQuestion": &gql.Field{
			Type: questionType,
			Args: gql.FieldConfigArgument{
				"title": &gql.ArgumentConfig{
					Type: gql.NewNonNull(gql.String),
				},
				"query": &gql.ArgumentConfig{
					Type: gql.NewNonNull(gql.String),
				},
			},
			Resolve: saveQuestion(storage),
		},
	}
	return gql.NewObject(gql.ObjectConfig{Name: "Mutation", Fields: fields})
}
