package query

import (
	"regexp"
	"time"
)

// Variables is map of "query variables"
type Variables map[string]string

type Evaler struct {
	Variables
}

func endOfDay() time.Time {
	now := time.Now().UTC()
	d := time.Duration(-now.Hour()) * time.Hour
	return now.Truncate(time.Hour).Add(d).Add(24*time.Hour - time.Second)
}

var defaults = Variables{
	"{timestamp}": endOfDay().Format(time.RFC3339),
	"{date}":      endOfDay().Format(time.RFC3339),
	"{today}":     endOfDay().Format(time.RFC3339),
	"{yesterday}": endOfDay().Add(-24*time.Hour - time.Nanosecond).Format(time.RFC3339),
}

func NewEvaler(vars Variables) *Evaler {
	merged := make(Variables, len(defaults))
	for k, v := range defaults {
		merged[k] = v
	}

	for k, v := range vars {
		merged[k] = v
	}
	return &Evaler{Variables: merged}
}

// Eval takes a query and evaluates it using vars
func (e *Evaler) Eval(query string) (string, error) {
	for key, value := range e.Variables {
		reg := regexp.MustCompile(key)

		query = reg.ReplaceAllString(query, "'"+value+"'")
	}
	return query, nil
}

// ExtractVariables extracts variables from a query
func ExtractVariables(query string) (variables Variables) {
	variables = make(Variables, 1)
	for key, value := range defaults {
		reg := regexp.MustCompile(key)

		if match := reg.FindString(query); len(match) > 0 {
			variables[key] = value
		}

	}
	return variables
}
