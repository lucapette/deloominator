package db

import (
	"strconv"
	"time"
)

type Dialect interface {
	TablesQuery() string
	ExtractCellInfo(interface{}) (Cell, Type)
	ConnectionString() string
	DBName() string
}

func inferType(value string) Type {
	_, err := strconv.ParseFloat(value, 32)
	if err == nil {
		return Number
	}

	_, err = time.Parse("2006-01-02 15:04:05", value)
	if err == nil {
		return Time
	}

	return Text
}
