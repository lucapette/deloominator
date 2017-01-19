package db

import "io"

type Type int

const (
	Number = iota
	Text   = iota
	Date   = iota
)

type Column struct {
	Name  string
	Value string
	Type  Type
}

type Row []Column

type Rows []Row

type Inquirer interface {
	Tables() (Rows, error)
	Query(string) (Rows, error)
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
