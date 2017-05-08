package charts

type DataType int

type DataTypes []DataType

//go:generate stringer -type=DataType -output=data_type_string.go
const (
	UnknownType DataType = iota
	Text
	Number
	Time
)
