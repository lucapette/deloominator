package db

type Cell struct {
	Value string
}

type Column struct {
	Name string
	Type Type
}

type Row []Cell

type QueryResult struct {
	Rows    []Row
	Columns []Column
}
