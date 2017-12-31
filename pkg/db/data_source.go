package db

import (
	"database/sql"
	"fmt"
	"net/url"
	"sort"

	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/lucapette/deloominator/pkg/query"
)

type DataSource struct {
	Dialect
	*sql.DB
}

type DataSources map[string]*DataSource

type Input struct {
	Query     string
	Variables map[string]string
}

func NewDataSources(sources []string) (dataSources DataSources, err error) {
	dataSources = make(DataSources, len(sources))

	for _, source := range sources {
		ds, err := NewDataSource(source)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"source": source,
				"err":    err.Error(),
			}).Print("datasource not available")

			continue
		}
		if err := ds.Ping(); err != nil {
			logrus.WithFields(logrus.Fields{
				"source": source,
				"err":    err.Error(),
			}).Print("datasource not available")

			continue
		}

		dataSources[ds.DBName()] = ds
	}

	if len(dataSources) == 0 {
		return dataSources, fmt.Errorf("no datasource available")
	}

	return dataSources, nil
}

func (dataSources DataSources) Close() {
	for _, ds := range dataSources {
		ds.Close()
	}
}

func NewDataSource(source string) (ds *DataSource, err error) {
	url, err := url.Parse(source)
	if err != nil {
		return nil, err
	}
	dialect, err := NewDialect(url)
	if err != nil {
		return nil, err
	}

	db, err := sql.Open(url.Scheme, dialect.ConnectionString())
	if err != nil {
		return nil, err
	}

	return &DataSource{Dialect: dialect, DB: db}, nil
}

// Tables returns the names of the available tables in the data source
func (ds *DataSource) Tables() (names []string, err error) {
	queryResult, err := ds.Query(Input{Query: ds.TablesQuery()})
	if err != nil {
		return names, err
	}

	names = make([]string, len(queryResult.Rows))
	for i, r := range queryResult.Rows {
		names[i] = r[0].Value
	}
	sort.Strings(names)

	return names, err
}

func (ds *DataSource) CreateDBIfNotExist() error {
	if !ds.IsUnknown(ds.Ping()) {
		return nil
	}

	db, err := sql.Open(ds.DriverName(), strings.Replace(ds.ConnectionString(), ds.DBName(), "", 1))
	if err != nil {
		return err
	}
	defer func() {
		if err := db.Close(); err != nil {
			logrus.WithFields(logrus.Fields{
				"storage_name": ds.DBName(),
				"err":          err.Error(),
			}).Print("could not close")
		}
	}()

	logrus.WithFields(logrus.Fields{
		"storage_name": ds.DBName(),
	}).Printf("creating storage")

	_, err = db.Exec(fmt.Sprintf("CREATE DATABASE %s", ds.DBName()))
	if err == nil {
		logrus.WithFields(logrus.Fields{
			"storage_name": ds.DBName(),
		}).Printf("storage created successfully")
	}

	return err
}

// Close closes the connection to the data source
func (ds *DataSource) Close() {
	logrus.WithFields(logrus.Fields{
		"storage_name": ds.DBName(),
	}).Print("closing db connection")

	if err := ds.DB.Close(); err != nil {
		logrus.WithFields(logrus.Fields{
			"storage_name": ds.DBName(),
			"err":          err.Error(),
		}).Print("could not close")
	}
}

func (ds *DataSource) Query(input Input) (qr QueryResult, err error) {
	evaler := query.NewEvaler(input.Variables)
	query, err := evaler.Eval(input.Query)
	if err != nil {
		return qr, err
	}

	// We use a prepare statement here so we can force MySQL binary protocol and get real types back. See:
	// https://github.com/go-sql-driver/mysql/issues/407#issuecomment-172583652
	statement, err := ds.DB.Prepare(query)
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
			value, colType := ds.ExtractCellInfo(res)

			cells[i] = value

			if qr.Columns[i].Type.isUnknown() {
				qr.Columns[i].Type = colType
			}
		}

		qr.Rows = append(qr.Rows, cells)
	}

	return qr, dbRows.Close()
}
