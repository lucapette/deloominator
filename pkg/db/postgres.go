package db

import (
	"fmt"
	"net/url"
	"time"
)

type Postgres struct {
	u *url.URL
}

func (pg *Postgres) TablesQuery() string {
	return `SELECT tablename FROM pg_tables WHERE schemaname = 'public' ORDER BY 1`
}

func (pg *Postgres) ExtractCellInfo(data interface{}) (cell Cell, colType Type) {
	switch data.(type) {
	case int64:
		cell = Cell{Value: fmt.Sprint(data)}
		colType = Number
	case string:
		cell = Cell{Value: data.(string)}
		colType = Text
	case []uint8:
		value := string(data.([]uint8))
		cell = Cell{Value: value}
		colType = inferType(value)
	case time.Time:
		cell = Cell{Value: data.(time.Time).Format(TimeFormat)}
		colType = Time
	default:
		cell = Cell{Value: fmt.Sprint(data)}
		colType = UnknownType
	}
	return cell, colType
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
