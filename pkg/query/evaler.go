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
	return now.Truncate(time.Hour).Add(d).Add(24*time.Hour - time.Nanosecond)
}

// NewEvaler returns an empty query variables map
func NewEvaler() *Evaler {
	return &Evaler{
		Variables: Variables{
			"{today}":     endOfDay().Format(time.RFC3339),
			"{yesterday}": endOfDay().Add(-24*time.Hour - time.Nanosecond).Format(time.RFC3339),
		},
	}
}

// Eval takes a query and evaluates it using vars
func (e *Evaler) Eval(query string) (string, error) {
	for key, value := range e.Variables {
		reg := regexp.MustCompile(key)

		query = reg.ReplaceAllString(query, "'"+value+"'")
	}
	return query, nil
}
