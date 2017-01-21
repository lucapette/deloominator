package db

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

type QueryResult struct {
	Rows []Row
}
