package db_test

import (
	"database/sql"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"testing"
	"time"

	"github.com/lucapette/deluminator/db"
)

type DBTemplate struct {
	Name string
}

func checkErr(err error, t *testing.T) {
	if err != nil {
		t.Fatal(err)
	}
}

func setupDB(driver string, t *testing.T) (*db.DSN, func()) {
	randName := "deluminator_" + strconv.Itoa(int(time.Now().UnixNano()+int64(os.Getpid())))

	dsn, err := db.NewDSN(fmt.Sprintf("%s://localhost/%s?sslmode=disable", driver, randName))
	checkErr(err, t)

	tmpfile, err := ioutil.TempFile("fixtures", "db_test")
	checkErr(err, t)
	defer func() {
		checkErr(os.Remove(tmpfile.Name()), t)
	}()

	source, err := ioutil.ReadFile(fmt.Sprintf("fixtures/%s.sql", driver))
	checkErr(err, t)

	tmpl := template.Must(template.New("ddl").Parse(string(source)))

	dbTemplate := DBTemplate{Name: randName}
	err = tmpl.Execute(tmpfile, dbTemplate)
	checkErr(err, t)

	err = exec.Command("psql", "-f", tmpfile.Name()).Run()
	checkErr(err, t)

	return dsn, func() {
		db, err := sql.Open(driver, fmt.Sprint(fmt.Sprintf("%s://localhost/?sslmode=disable", driver)))
		checkErr(err, t)

		_, err = db.Exec(fmt.Sprintf("DROP DATABASE %s;", randName))
		checkErr(err, t)

		checkErr(db.Close(), t)
	}
}
