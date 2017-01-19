package db

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

type Executor struct {
	db *sqlx.DB
}

func (ex *Executor) Query(query string) (rows Rows, err error) {
	dbRows, err := ex.db.Queryx(query)
	if err != nil {
		return rows, err
	}

	dbCols, err := dbRows.Columns()
	if err != nil {
		return rows, err
	}

	for dbRows.Next() {
		results, err := dbRows.SliceScan()
		if err != nil {
			return rows, err
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

		rows = append(rows, cols)
	}

	return rows, dbRows.Close()
}
