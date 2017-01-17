package db

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

type MySQL struct {
	dsn *DSN
	db  *sql.DB
}

func (my *MySQL) Tables() (tables []string, err error) {
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

func NewMySQL(dsn *DSN) (my *MySQL, err error) {
	db, err := sql.Open(dsn.Driver, dsn.Format())
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return &MySQL{dsn: dsn, db: db}, nil
}

func (my *MySQL) DSN() *DSN {
	return my.dsn
}

func (my *MySQL) Close() error {
	return my.db.Close()
}

func (my *MySQL) Query(query string) (rows Rows, err error) {
	return rows, err
}
