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
}

func NewDialect(url *url.URL) (dialect Dialect, err error) {
	switch url.Scheme {
	case "postgres":
		return NewPostgresDialect(url.String())
	case "mysql":
		return NewMySQLDialect(url.String())
	}

	return dialect, err
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
