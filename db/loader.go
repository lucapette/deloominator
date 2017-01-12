package db

type Loader interface {
	Tables() ([]string, error)
	DSN() *DSN
	Close() error
}

type Loaders map[string]Loader

func NewLoaders(dataSources []string) (loaders Loaders, err error) {
	loaders = make(Loaders, len(dataSources))

	for _, source := range dataSources {
		ds, err := NewDSN(source)
		if err != nil {
			return nil, err
		}

		loader, err := NewLoader(ds)
		if err != nil {
			return nil, err
		}

		loaders[loader.DSN().DBName] = loader
	}

	return loaders, nil
}

func NewLoader(dsn *DSN) (loader Loader, err error) {
	switch dsn.Driver {
	case "mysql":
		loader, err = NewMyLoader(dsn)
	case "postgres":
		loader, err = NewPGLoader(dsn)
	}
	return loader, err
}
