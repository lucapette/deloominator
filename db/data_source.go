package db

import "io"

type Inquirer interface {
	Tables() (QueryResult, error)
	Query(string) (QueryResult, error)
}

type DataSource interface {
	DSN() *DSN
	io.Closer
	Inquirer
}

type DataSources map[string]DataSource

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

		dataSources[ds.DSN().DBName] = ds
	}

	return dataSources, nil
}

func NewDataSource(dsn *DSN) (ds DataSource, err error) {
	switch dsn.Driver {
	case "mysql":
		ds, err = NewMySQL(dsn)
	case "postgres":
		ds, err = NewPostgres(dsn)
	}
	return ds, err
}
