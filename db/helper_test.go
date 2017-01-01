package db_test

import (
	"bytes"
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

func setupDB(driver db.DriverType, t *testing.T) (*db.DSN, func()) {
	var dsn *db.DSN
	var cleanup func()

	switch driver {
	case db.Postgres:
		dsn, cleanup = setupPostgres(t)
	case db.MySQL:
		dsn, cleanup = setupMysql(t)
	}

	return dsn, cleanup
}

func setupPostgres(t *testing.T) (*db.DSN, func()) {
	randName := "deluminator_" + strconv.Itoa(int(time.Now().UnixNano()+int64(os.Getpid())))

	dsn, err := db.NewDSN(fmt.Sprintf("postgres://localhost/%s?sslmode=disable", randName))
	if err != nil {
		t.Fatal(err)
	}

	tmpfile, err := ioutil.TempFile("fixtures", "db_test")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		err = os.Remove(tmpfile.Name())
		if err != nil {
			t.Fatal(err)
		}
	}()

	source, err := ioutil.ReadFile("fixtures/postgres.sql")
	if err != nil {
		t.Fatal(err)
	}

	tmpl := template.Must(template.New("ddlPsql").Parse(string(source)))

	dbTemplate := DBTemplate{Name: randName}
	err = tmpl.Execute(tmpfile, dbTemplate)
	if err != nil {
		t.Fatal(err)
	}

	err = exec.Command("psql", "-f", tmpfile.Name()).Run()
	if err != nil {
		t.Fatal(err)
	}

	return dsn, func() {
		db, err := sql.Open("postgres", "postgres://localhost/?sslmode=disable")
		if err != nil {
			t.Fatal(err)
		}

		_, err = db.Exec(fmt.Sprintf("DROP DATABASE %s;", randName))
		if err != nil {
			t.Fatal(err)
		}

		err = db.Close()
		if err != nil {
			t.Fatal(err)
		}
	}
}

func setupMysql(t *testing.T) (*db.DSN, func()) {
	randName := "deluminator_" + strconv.Itoa(int(time.Now().UnixNano()+int64(os.Getpid())))

	dsn, err := db.NewDSN(fmt.Sprintf("mysql://root:root@localhost/%s", randName))
	if err != nil {
		t.Fatal(err)
	}

	source, err := ioutil.ReadFile("fixtures/mysql.sql")
	if err != nil {
		t.Fatal(err)
	}

	tmpl := template.Must(template.New("ddl").Parse(string(source)))

	var query bytes.Buffer
	dbTemplate := DBTemplate{Name: randName}
	err = tmpl.Execute(&query, dbTemplate)
	if err != nil {
		t.Fatal(err)
	}

	db, err := sql.Open("mysql", "root:root@/?multiStatements=true")
	if err != nil {
		t.Fatal(err)
	}

	_, err = db.Exec(query.String())
	if err != nil {
		t.Fatal(err)
	}

	return dsn, func() {
		_, err = db.Exec(fmt.Sprintf("DROP DATABASE %s;", randName))
		if err != nil {
			t.Fatal(err)
		}

		err = db.Close()
		if err != nil {
			t.Fatal(err)
		}
	}
}
