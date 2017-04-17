package db

import "fmt"

type Postgres struct {
}

func (pg *Postgres) TablesQuery() string {
	return `SELECT tablename FROM pg_tables WHERE schemaname = 'public' ORDER BY 1`
}

func (pg *Postgres) ExtractCellInfo(data interface{}) Cell {
	var value string
	switch t := data.(type) {
	case []uint8:
		value = string(data.([]uint8))
	default:
		value = fmt.Sprint(t)
	}
	return Cell{Value: value}
}

func NewPostgresDialect() *Postgres {
	return &Postgres{}
}
