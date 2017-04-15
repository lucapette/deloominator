package db

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Dialect interface {
	TablesQuery() string
}

type DataSource struct {
	dialect Dialect
	db      *sqlx.DB
	dsn     *DSN
}

type DataSources map[string]*DataSource

func NewDataSources(sources []string) (dataSources DataSources, err error) {
	dataSources = make(DataSources, len(sources))

	for _, source := range sources {
		dsn, err := NewDSN(source)
		if err != nil {
			return nil, err
		}

		ds, err := NewDataSource(dsn)
		if err != nil {
			return nil, err
		}

		dataSources[ds.dsn.DBName] = ds
	}

	return dataSources, nil
}

func NewDataSource(dsn *DSN) (ds *DataSource, err error) {
	db, err := sqlx.Open(dsn.Driver, dsn.Format())
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	var dialect Dialect
	switch dsn.Driver {
	case "postgres":
		dialect = NewPostgresDialect()
	case "mysql":
		dialect = NewMySQLDialect()
	}

	return &DataSource{dsn: dsn, dialect: dialect, db: db}, nil
}

func (ds *DataSource) Tables() (QueryResult, error) {
	return ds.Query(ds.dialect.TablesQuery())
}

func (ds *DataSource) Name() string {
	return ds.dsn.DBName
}

func (ds *DataSource) Close() error {
	return ds.db.Close()
}

func (ds *DataSource) String() string {
	return ds.dsn.String()
}

func (ds *DataSource) Query(input string) (qr QueryResult, err error) {
	dbRows, err := ds.db.Queryx(input)
	if err != nil {
		return qr, err
	}

	dbCols, err := dbRows.Columns()
	if err != nil {
		return qr, err
	}

	qr.Columns = make([]Column, len(dbCols))
	for i, dbCol := range dbCols {
		qr.Columns[i].Name = dbCol
	}

	for dbRows.Next() {
		results, err := dbRows.SliceScan()
		if err != nil {
			return qr, err
		}

		cells := make([]Cell, len(results))
		for i, res := range results {
			var value string
			switch t := res.(type) {
			case []uint8: // this happens with MySQL result and we don't know why yet. To investigate.
				value = string(res.([]uint8))
			default:
				value = fmt.Sprint(t)
			}
			cells[i] = Cell{Value: value}
		}

		qr.Rows = append(qr.Rows, cells)
	}

	return qr, dbRows.Close()
}
