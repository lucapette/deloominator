package graphql

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/Sirupsen/logrus"
	gql "github.com/graphql-go/graphql"
	ast "github.com/graphql-go/graphql/language/ast"
	"github.com/lucapette/deloominator/pkg/charts"
	"github.com/lucapette/deloominator/pkg/db"
	"github.com/lucapette/deloominator/pkg/db/storage"
	"github.com/lucapette/deloominator/pkg/query"
)

func resolveDataSources(dbDataSources db.DataSources) func(p gql.ResolveParams) (interface{}, error) {
	return func(p gql.ResolveParams) (interface{}, error) {
		var dataSources []*dataSource

		for _, ds := range dbDataSources {
			logrus.WithField("schema_name", ds.DBName()).Info("query metadata")

			start := time.Now()

			names, err := ds.Tables()
			if err != nil {
				return dataSources, err
			}

			ts := make([]*table, len(names))
			for i, name := range names {
				ts[i] = &table{Name: name}
			}

			logrus.WithFields(logrus.Fields{
				"schema_name": ds.DBName(),
				"n_tables":    len(names),
				"spent":       time.Since(start),
			}).Info("tables loaded")

			dataSources = append(dataSources, &dataSource{Name: ds.DBName(), Tables: ts})
		}

		return dataSources, nil
	}
}

type question struct {
	ID          int         `json:"id"`
	Title       string      `json:"title"`
	Description string      `json:"description"`
	CreatedAt   string      `json:"createdAt"`
	UpdatedAt   string      `json:"updatedAt"`
	DataSource  string      `json:"dataSource"`
	Query       string      `json:"query"`
	Variables   []variable  `json:"variables"`
	Results     interface{} `json:"results"`
}

func needsResults(p gql.ResolveParams) bool {
	for _, v := range p.Info.FieldASTs {
		for _, s := range v.SelectionSet.Selections {
			if f, ok := s.(*ast.Field); ok && f.Name.Value == "results" {
				return true
			}
		}
	}
	return false
}

func convertQuestion(in *storage.Question) (out question, err error) {
	out = question{
		ID:          in.ID,
		Title:       in.Title,
		Description: in.Description,
		DataSource:  in.DataSource,
		Query:       in.Query,
		CreatedAt:   in.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   in.UpdatedAt.Format(time.RFC3339),
	}

	if len(in.Variables) > 0 {
		if err := json.Unmarshal([]byte(in.Variables), &out.Variables); err != nil {
			return out, err
		}
	}
	return out, nil
}

func resolveQuestion(dataSources db.DataSources, s *storage.Storage) func(p gql.ResolveParams) (interface{}, error) {
	return func(p gql.ResolveParams) (interface{}, error) {
		id, err := strconv.Atoi(p.Args["id"].(string))
		if err != nil {
			return nil, err
		}

		in, err := s.FindQuestion(id)
		if err != nil {
			return nil, err
		}

		question, err := convertQuestion(in)
		if err != nil {
			return question, err
		}

		if needsResults(p) {
			queryResultResolver := resolveQuery(dataSources)
			question.Results, err = queryResultResolver(gql.ResolveParams{
				Args: map[string]interface{}{
					"dataSource": question.DataSource,
					"query":      question.Query,
					"variables":  question.Variables,
				},
			})

			if err != nil {
				return question, err
			}
		}
		return question, nil
	}
}

func resolveQuestions(dataSources db.DataSources, s *storage.Storage) func(p gql.ResolveParams) (interface{}, error) {
	return func(p gql.ResolveParams) (interface{}, error) {
		dbQuestions, err := s.AllQuestions()
		if err != nil {
			return nil, err
		}
		needsResults := needsResults(p)
		queryResultResolver := resolveQuery(dataSources)

		questions := make([]question, len(dbQuestions))
		for i, q := range dbQuestions {
			question, err := convertQuestion(q)
			if err != nil {
				return nil, err
			}
			if needsResults {
				question.Results, err = queryResultResolver(gql.ResolveParams{
					Args: map[string]interface{}{
						"dataSource": question.DataSource,
						"query":      question.Query,
						"variables":  question.Variables,
					},
				})
				if err != nil {
					return nil, err
				}
			}
			questions[i] = question
		}
		return questions, nil
	}
}

func resolveQuery(dataSources db.DataSources) func(p gql.ResolveParams) (interface{}, error) {
	return func(p gql.ResolveParams) (interface{}, error) {
		dataSource := p.Args["dataSource"].(string)
		q := p.Args["query"].(string)
		b, err := json.Marshal(p.Args["variables"])
		if err != nil {
			return nil, err
		}

		variables := make([]variable, 0)
		err = json.Unmarshal(b, &variables)
		if err != nil {
			return nil, err
		}
		vars := make([]query.Variable, len(variables))
		for i, v := range variables {
			vars[i].Name = v.Name
			vars[i].Value = v.Value
		}

		logrus.WithFields(logrus.Fields{
			"dataSource": dataSource,
			"query":      q,
			"variables":  vars,
			"args":       p.Args,
		}).Infof("Query requested")

		qr, err := dataSources[dataSource].Query(db.Input{Query: q, Variables: vars})
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
