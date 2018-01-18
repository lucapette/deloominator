package graphql

import (
	"encoding/json"

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

		variables, err := json.Marshal(p.Args["variables"])
		if err != nil {
			return nil, err
		}
		question.Variables = string(variables)

		if description, ok := p.Args["description"].(string); ok {
			question.Description = description
		}

		q, err := s.InsertQuestion(&question)
		if err != nil {
			return nil, err
		}

		return convertQuestion(q)
	}
}
