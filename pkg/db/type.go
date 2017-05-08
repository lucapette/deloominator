package db

type Type int

//go:generate stringer -type=Type
const (
	UnknownType Type = iota
	Number
	Text
	Time
)

func (t Type) isUnknown() bool {
	return t == UnknownType
}
