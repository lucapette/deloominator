package db

import (
	"database/sql"
	"net/url"

	log "github.com/Sirupsen/logrus"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

type DriverType int

const (
	PostgresDriver DriverType = iota
	MySQLDriver
)

type DataSource struct {
	dialect Dialect
	db      *sql.DB
}

type DataSources map[string]*DataSource

func NewDataSources(sources []string) (dataSources DataSources, err error) {
	dataSources = make(DataSources, len(sources))

	for _, source := range sources {
		ds, err := NewDataSource(source)
		if err != nil {
			return nil, err
		}

		dataSources[ds.Name()] = ds
	}

	return dataSources, nil
}

func (dataSources DataSources) Shutdown() {
	for _, ds := range dataSources {
		err := ds.Close()
		if err != nil {
			log.Fatal(err.Error())
		}
	}
}

func parseDriver(source string) (string, error) {
	url, err := url.Parse(source)
	if err != nil {
		return "", err
	}

	return url.Scheme, nil
}

func NewDataSource(source string) (ds *DataSource, err error) {
	driver, err := parseDriver(source)
	if err != nil {
		return nil, err
	}

	var dialect Dialect
	switch driver {
	case "postgres":
		dialect, err = NewPostgresDialect(source)
	case "mysql":
		dialect, err = NewMySQLDialect(source)
	}
	if err != nil {
		return nil, err
	}

	db, err := sql.Open(driver, dialect.ConnectionString())
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return &DataSource{dialect: dialect, db: db}, nil
}

func (ds *DataSource) Tables() (QueryResult, error) {
	return ds.Query(ds.dialect.TablesQuery())
}

func (ds *DataSource) Name() string {
	return ds.dialect.DBName()
}

func (ds *DataSource) Close() error {
	return ds.db.Close()
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
			value, colType := ds.dialect.ExtractCellInfo(res)

			cells[i] = value

			if qr.Columns[i].Type.isUnknown() {
				qr.Columns[i].Type = colType
			}
		}

		qr.Rows = append(qr.Rows, cells)
	}

	return qr, dbRows.Close()
}
