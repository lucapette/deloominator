package db

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

type MyLoader struct {
	dsn *DSN
	db  *sql.DB
}

func (my *MyLoader) Tables() (tables []string, err error) {
	rows, err := my.db.Query(`SHOW TABLES;`)

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

func NewMyLoader(dsn *DSN) (my *MyLoader, err error) {
	db, err := sql.Open(dsn.Driver, dsn.Format())
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return &MyLoader{dsn: dsn, db: db}, nil
}

func (my *MyLoader) DSN() *DSN {
	return my.dsn
}

func (my *MyLoader) Close() error {
	return my.db.Close()
}
