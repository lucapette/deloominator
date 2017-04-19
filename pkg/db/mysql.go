package db

import (
	"fmt"
	"strings"

	"github.com/go-sql-driver/mysql"
)

type MySQL struct {
	cfg *mysql.Config
}

func (my *MySQL) TablesQuery() string {
	return `SHOW TABLES`
}

func (my *MySQL) ExtractCellInfo(data interface{}) (cell Cell, colType Type) {
	switch t := data.(type) {
	case int64:
		cell = Cell{Value: fmt.Sprint(data)}
		colType = Number
	case []uint8:
		value := string(data.([]uint8))
		cell = Cell{Value: value}
		colType = inferType(value)
	default:
		cell = Cell{Value: fmt.Sprint(t)}
		colType = UnknownType
	}
	return cell, colType
}

func (my *MySQL) ConnectionString() string {
	return my.cfg.FormatDSN()
}

func (my *MySQL) DBName() string {
	return my.cfg.DBName
}

func NewMySQLDialect(source string) (*MySQL, error) {
	config, err := mysql.ParseDSN(strings.Replace(source, "mysql://", "", 1))
	if err != nil {
		return nil, err
	}

	return &MySQL{cfg: config}, nil
}
