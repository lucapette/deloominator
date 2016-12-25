package db

type DataSource struct {
	Name string
}

type DataSources map[string]DataSource

func NewSources(dataSources []string) (sources DataSources, err error) {
	sources = make(DataSources, len(dataSources))

	// for source := range sources {
	// }

	return sources, nil
}
