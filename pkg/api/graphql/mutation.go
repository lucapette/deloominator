package graphql

import (
	gql "github.com/graphql-go/graphql"
	"github.com/lucapette/deloominator/pkg/db/storage"
)

func saveQuestion(s *storage.Storage) func(p gql.ResolveParams) (interface{}, error) {
	return func(p gql.ResolveParams) (interface{}, error) {
		title := p.Args["title"].(string)
		query := p.Args["query"].(string)
		dataSource := p.Args["dataSource"].(string)

		question, err := s.InsertQuestion(title, query, dataSource)
		if err != nil {
			return nil, err
		}
		return *question, nil
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
			},
			Resolve: saveQuestion(s),
		},
	}
	return gql.NewObject(gql.ObjectConfig{Name: "mutation", Fields: fields})
}
