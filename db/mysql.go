package db

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type MySQL struct {
	dsn *DSN
	*Executor
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

func NewMySQL(dsn *DSN) (*MySQL, error) {
	db, err := sqlx.Open(dsn.Driver, dsn.Format())
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return &MySQL{dsn: dsn, Executor: &Executor{db: db}}, nil
}

func (my *MySQL) DSN() *DSN {
	return my.dsn
}

func (my *MySQL) Close() error {
	return my.db.Close()
}
