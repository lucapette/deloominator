package db

type DataSource struct {
	Name string
}

type DataSources map[string]*DataSource

func NewSources(dataSources []string) (sources DataSources, err error) {
	sources = make(DataSources, len(dataSources))

	for _, source := range dataSources {
		dsn, err := ParseDSN(source)
		if err != nil {
			return nil, err
		}

		sources[dsn.DBName] = &DataSource{
			Name: dsn.DBName,
		}
	}

	return sources, nil
}
