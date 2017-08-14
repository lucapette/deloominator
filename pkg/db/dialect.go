package db

import (
	"net/url"
	"strconv"
	"time"
)

const timeFormat = "2006-01-02 15:04:05"

type Dialect interface {
	TablesQuery() string
	ExtractCellInfo(interface{}) (Cell, Type)
	ConnectionString() string
	DBName() string
	IsUnknown(error) bool
	DriverName() string
}

func NewDialect(u *url.URL) (Dialect, error) {
	switch u.Scheme {
	case "postgres":
		return NewPostgresDialect(u)
	case "mysql":
		return NewMySQLDialect(u)
	}

	return nil, nil
}

func inferType(value string) Type {
	if _, err := strconv.ParseFloat(value, 32); err == nil {
		return Number
	}

	if _, err := time.Parse(timeFormat, value); err == nil {
		return Time
	}

	return Text
}
