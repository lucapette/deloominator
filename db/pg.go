package db

import (
	"database/sql"
	"fmt"

	log "github.com/Sirupsen/logrus"
	_ "github.com/lib/pq"
)

type PGLoader struct {
	dsn *DataSource
}

func (pg *PGLoader) Tables() (tables []string, err error) {
	db, err := sql.Open(pg.DSN().Driver, fmt.Sprint(pg.DSN()))
	if err != nil {
		return tables, err
	}

	err = db.Ping()
	if err != nil {
		return tables, err
	}

	rows, err := db.Query(`SELECT tablename FROM pg_tables where schemaname = 'public'`)

	if err != nil {
		return tables, err
	}

	defer rows.Close()

	for rows.Next() {
		var name string
		err = rows.Scan(&name)
		if err != nil {
			return tables, err
		}

		log.Info(name)
	}

	return tables, err
}

func NewPGLoader(DSN *DataSource) (ds *PGLoader, err error) {
	return &PGLoader{dsn: DSN}, nil
}

func (pg *PGLoader) DSN() *DataSource {
	return pg.dsn
}
