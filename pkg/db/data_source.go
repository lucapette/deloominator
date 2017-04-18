package db

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

type Dialect interface {
	TablesQuery() string
	ExtractCellInfo(interface{}) Cell
}

type DataSource struct {
	dialect Dialect
	db      *sql.DB
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

		dataSources[ds.Name()] = ds
	}

	return dataSources, nil
}

func NewDataSource(dsn *DSN) (ds *DataSource, err error) {
	db, err := sql.Open(dsn.Driver, dsn.Format())
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
	// We use a prepare statement here so we can force MySQL binary protocol and
	// get real types back. See: https://github.com/go-sql-driver/mysql/issues/407#issuecomment-172583652
	statement, err := ds.db.Prepare(input)
	if err != nil {
		return qr, err
	}

	dbRows, err := statement.Query()
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
		columns := make([]interface{}, len(dbCols))
		columnPointers := make([]interface{}, len(dbCols))

		for i := 0; i < len(dbCols); i++ {
			columnPointers[i] = &columns[i]
		}

		if err := dbRows.Scan(columnPointers...); err != nil {
			return qr, err
		}

		cells := make([]Cell, len(columns))
		for i, res := range columns {
			cells[i] = ds.dialect.ExtractCellInfo(res)
		}

		qr.Rows = append(qr.Rows, cells)
	}

	return qr, dbRows.Close()
}
