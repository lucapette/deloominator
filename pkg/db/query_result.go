package db

type Type int

const (
	UnknownType Type = iota
	Number
	Text
	Time
)

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

func (t Type) isUnknown() bool {
	return t == UnknownType
}

func (t Type) String() string {
	switch t {
	case Text:
		return "Text"
	case Number:
		return "Number"
	case Time:
		return "Time"
	}

	return "Unknown"
}
