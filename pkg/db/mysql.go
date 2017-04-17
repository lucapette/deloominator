package db

import "fmt"

type MySQL struct {
}

func (my *MySQL) TablesQuery() string {
	return `SHOW TABLES`
}

func (my *MySQL) ExtractCellInfo(data interface{}) Cell {
	var value string
	switch t := data.(type) {
	case []uint8:
		value = string(data.([]uint8))
	default:
		value = fmt.Sprint(t)
	}
	return Cell{Value: value}
}

func NewMySQLDialect() *MySQL {
	return &MySQL{}
}
