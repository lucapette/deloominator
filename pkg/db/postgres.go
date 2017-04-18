package db

import (
	"fmt"
	"net/url"
)

type Postgres struct {
	u *url.URL
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

func (pg *Postgres) ConnectionString() string {
	return pg.u.String()
}

func (pg *Postgres) DBName() string {
	return pg.u.Path[1:]
}

func NewPostgresDialect(source string) (*Postgres, error) {
	u, err := url.Parse(source)
	if err != nil {
		return nil, err
	}

	return &Postgres{u: u}, nil
}
