package graphql

import (
	"encoding/json"

	"github.com/Sirupsen/logrus"

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
		logrus.Infof("%v", string(variables))
		question.Variables = string(variables)

		q, err := s.InsertQuestion(&question)
		if err != nil {
			return nil, err
		}

		return convertQuestion(q)
	}
}
