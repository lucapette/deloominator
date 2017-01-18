package db

import (
	"github.com/jmoiron/sqlx"

	_ "github.com/lib/pq"
)

type Postgres struct {
	dsn *DSN
	*Executor
}

func (pg *Postgres) Tables() (tables []string, err error) {
	rows, err := pg.db.Query(`SELECT tablename FROM pg_tables WHERE schemaname = 'public' ORDER BY 1`)
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

func NewPostgres(dsn *DSN) (*Postgres, error) {
	db, err := sqlx.Open(dsn.Driver, dsn.Format())
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return &Postgres{dsn: dsn, Executor: &Executor{db: db}}, nil
}

func (pg *Postgres) DSN() *DSN {
	return pg.dsn
}

func (pg *Postgres) Close() error {
	return pg.db.Close()
}
