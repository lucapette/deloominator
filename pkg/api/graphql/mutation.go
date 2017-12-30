package graphql

import (
	gql "github.com/graphql-go/graphql"
	"github.com/lucapette/deloominator/pkg/db/storage"
)

func saveQuestion(s *storage.Storage) func(p gql.ResolveParams) (interface{}, error) {
	return func(p gql.ResolveParams) (interface{}, error) {
		question := storage.Question{
			Title:      p.Args["title"].(string),
			Query:      p.Args["query"].(string),
			DataSource: p.Args["dataSource"].(string),
		}

		if variables, ok := p.Args["variables"].(string); ok {
			question.Variables = variables
		}

		q, err := s.InsertQuestion(&question)
		if err != nil {
			return nil, err
		}

		return *q, nil
	}
}

func mutation(s *storage.Storage) *gql.Object {
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
				"dataSource": &gql.ArgumentConfig{
					Type: gql.NewNonNull(gql.String),
				},
				"variables": &gql.ArgumentConfig{
					Type: gql.String,
				},
			},
			Resolve: saveQuestion(s),
		},
	}
	return gql.NewObject(gql.ObjectConfig{Name: "mutation", Fields: fields})
}
