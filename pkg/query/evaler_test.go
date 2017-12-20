package query_test

import (
	"testing"

	"github.com/lucapette/deloominator/pkg/query"
)

var variables = query.Variables{
	"{today}":     "2017-10-13T23:59:59Z",
	"{yesterday}": "2017-10-12T23:59:59Z",
}

func TestEval(t *testing.T) {
	type args struct {
		query string
		vars  query.Variables
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"no variables",
			args{
				"select * from answers",
				variables,
			},
			"select * from answers",
		},
		{
			"{today}",
			args{
				"select * from answers where created_at < {today}",
				variables,
			},
			`select * from answers where created_at < "2017-10-13T23:59:59Z"`,
		},
		{
			"{yesterday}",
			args{
				"select * from answers where created_at >= {yesterday}",
				variables,
			},
			`select * from answers where created_at >= "2017-10-12T23:59:59Z"`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			evaler := query.NewEvaler()
			evaler.Variables = variables
			got, err := evaler.Eval(tt.args.query)
			if err != nil {
				t.Errorf("could not eval: %v", err)
			}

			if got != tt.want {
				t.Errorf("Eval() = %v, want %v", got, tt.want)
			}
		})
	}
}
