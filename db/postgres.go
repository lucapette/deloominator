package db

import (
	"fmt"

	"github.com/jmoiron/sqlx"

	_ "github.com/lib/pq"
)

type Postgres struct {
	dsn *DSN
	db  *sqlx.DB
}

func (pg *Postgres) Tables() (tables []string, err error) {
	rows, err := pg.db.Query(`SELECT tablename FROM pg_tables where schemaname = 'public' order by 1`)
	if err != nil {
		return tables, err
	}

	for rows.Next() {
		var name string
		err = rows.Scan(&name)
		if err != nil {
			return tables, err
		}

		tables = append(tables, name)
	}

	return tables, err
}

func NewPostgres(dsn *DSN) (pg *Postgres, err error) {
	db, err := sqlx.Open(dsn.Driver, dsn.Format())
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return &Postgres{dsn: dsn, db: db}, nil
}

func (pg *Postgres) DSN() *DSN {
	return pg.dsn
}

func (pg *Postgres) Close() error {
	return pg.db.Close()
}

func (pg *Postgres) Query(query string) (rows Rows, err error) {
	dbRows, err := pg.db.Queryx(query)
	if err != nil {
		return rows, err
	}

	dbCols, err := dbRows.Columns()
	if err != nil {
		return rows, err
	}

	for dbRows.Next() {
		results, err := dbRows.SliceScan()
		if err != nil {
			return rows, err
		}

		cols := make(Row, len(results))
		for i, _ := range results {
			cols[i] = &Column{
				Name:  dbCols[i],
				Value: fmt.Sprint(results[i]),
			}
			// add a type switch here and do the converstion to the
			// deluminator type. For the Value we can keep string for now
		}

		rows = append(rows, &cols)
	}

	err = dbRows.Close()
	return rows, err
}
