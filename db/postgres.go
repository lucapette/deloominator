package db

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type Postgres struct {
	dsn *DSN
	db  *sql.DB
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
	db, err := sql.Open(dsn.Driver, dsn.Format())
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
