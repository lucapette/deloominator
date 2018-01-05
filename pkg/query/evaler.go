package query

import (
	"regexp"
	"time"
)

// Variable maps a single query variable
type Variable struct {
	Name  string
	Value string
}

type Evaler struct {
	vars map[string]string
}

func endOfDay() time.Time {
	now := time.Now().UTC()
	d := time.Duration(-now.Hour()) * time.Hour
	return now.Truncate(time.Hour).Add(d).Add(24*time.Hour - time.Second)
}

var defaults = map[string]string{
	"timestamp": endOfDay().Format(time.RFC3339),
	"date":      endOfDay().Format("2006-01-02"),
	"today":     endOfDay().Format("2006-01-02"),
	"yesterday": endOfDay().Add(-24*time.Hour - time.Nanosecond).Format("2006-01-02"),
}

func NewEvaler(vars []Variable) *Evaler {
	merged := make(map[string]string, len(defaults))

	for k, v := range defaults {
		merged[k] = v
	}

	for _, v := range vars {
		merged[v.Name] = v.Value
	}
	return &Evaler{vars: merged}
}

// Eval takes a query and evaluates it using vars
func (e *Evaler) Eval(query string) (evaled string) {
	evaled = query
	for k, v := range e.vars {
		reg := regexp.MustCompile("{" + k + "}")

		evaled = reg.ReplaceAllString(evaled, "'"+v+"'")
	}
	return evaled
}

// Eval takes a query and returns the used query variables
func (e *Evaler) ExtractVariables(query string) (variables []Variable) {
	variables = make([]Variable, 0)
	for k, v := range e.vars {
		reg := regexp.MustCompile("{" + k + "}")

		if reg.MatchString(query) {
			variables = append(variables, Variable{Name: k, Value: v})
		}
	}
	return variables
}
