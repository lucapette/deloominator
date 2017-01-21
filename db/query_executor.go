package db

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

type Executor struct {
	db *sqlx.DB
}

func (ex *Executor) Query(input string) (qr QueryResult, err error) {
	dbRows, err := ex.db.Queryx(input)
	if err != nil {
		return qr, err
	}

	dbCols, err := dbRows.Columns()
	if err != nil {
		return qr, err
	}

	for dbRows.Next() {
		results, err := dbRows.SliceScan()
		if err != nil {
			return qr, err
		}

		cols := make(Row, len(results))
		for i, res := range results {
			var value string
			switch t := res.(type) {
			case []uint8: // this happens with MySQL result and we don't know why yet. To investigate.
				value = string(res.([]uint8))
			default:
				value = fmt.Sprint(t)
			}
			cols[i] = Column{
				Name:  dbCols[i],
				Value: value,
			}
		}

		qr.Rows = append(qr.Rows, cols)
	}

	return qr, dbRows.Close()
}
