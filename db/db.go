package db

type Loader interface {
	Tables() ([]string, error)
	DSN() *DSN
}

type Loaders map[string]Loader

func NewLoaders(dataSources []string) (loaders Loaders, err error) {
	loaders = make(Loaders, len(dataSources))

	for _, source := range dataSources {
		ds, err := NewDSN(source)
		if err != nil {
			return nil, err
		}

		pg, err := NewPGLoader(ds)
		if err != nil {
			return nil, err
		}

		loaders[pg.DSN().DBName] = pg
	}

	return loaders, nil
}
