package db

import (
	"fmt"
	"strings"

	mysql "github.com/go-sql-driver/mysql"
)

type MySQL struct {
	cfg *mysql.Config
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
