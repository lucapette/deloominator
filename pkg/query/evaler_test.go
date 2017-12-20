package query_test

import (
	"testing"

	"github.com/lucapette/deloominator/pkg/query"
)

var variables = query.Variables{
	"{timestamp}": "2017-10-13T23:59:59Z",
}

func TestEval(t *testing.T) {
	type args struct {
		query string
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
			},
			"select * from answers",
		},
		{
			"{timestamp}",
			args{
				"select * from answers where created_at < {timestamp}",
			},
			`select * from answers where created_at < '2017-10-13T23:59:59Z'`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			evaler := query.NewEvaler(variables)
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
