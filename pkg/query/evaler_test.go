package query_test

import (
	"testing"

	"github.com/lucapette/deloominator/pkg/query"
)

var variables = query.Variables{
	"{timestamp}": "2017-10-13T23:59:59Z",
	"{today}":     "2017-10-15",
	"{yesterday}": "2015-04-15",
}

func TestEval(t *testing.T) {
	tests := []struct {
		name  string
		query string
		want  string
	}{
		{
			"no variables",
			"select * from answers",
			"select * from answers",
		},
		{
			"{timestamp}",
			"select * from answers where created_at < {timestamp}",
			`select * from answers where created_at < '2017-10-13T23:59:59Z'`,
		},
		{
			"{today}",
			"select * from answers where created_at < {today}",
			`select * from answers where created_at < '2017-10-15'`,
		},
		{
			"{yesterday}",
			"select * from answers where created_at < {yesterday}",
			`select * from answers where created_at < '2015-04-15'`,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			evaler := query.NewEvaler(variables)
			got, err := evaler.Eval(test.query)
			if err != nil {
				t.Errorf("could not eval: %v", err)
			}

			if got != test.want {
				t.Errorf("Eval() = %v, want %v", got, test.want)
			}
		})
	}
}
