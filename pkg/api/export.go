package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"goji.io/pat"

	"github.com/Sirupsen/logrus"
	"github.com/lucapette/deloominator/pkg/db"
)

type ExportPayload struct {
	Source string `json:"source"`
	Query  string `json:"query"`
}

func exportHandler(dataSources db.DataSources) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		format := pat.Param(r, "format")

		logrus.WithFields(logrus.Fields{
			"path":   r.URL.Path,
			"method": r.Method,
			"format": format,
		}).Print("export started")

		query, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		payload := ExportPayload{}

		err = json.Unmarshal(query, &payload)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		qr, err := dataSources[payload.Source].Query(payload.Query)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		data := &bytes.Buffer{}

		line := make([]string, len(qr.Columns))
		for i, col := range qr.Columns {
			line[i] = col.Name
		}
		fmt.Fprintln(data, strings.Join(line, ","))

		for _, row := range qr.Rows {
			for i, cell := range row {
				if len(cell.Value) == 0 {
					line[i] = ""
				} else {
					line[i] = cell.Value
				}
			}
			fmt.Fprintln(data, strings.Join(line, ","))
		}

		_, err = w.Write(data.Bytes())
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	}
}
