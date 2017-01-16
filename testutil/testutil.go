package testutil

import (
	"bytes"
	"database/sql"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/lucapette/deluminator/app"
	"github.com/lucapette/deluminator/db"
)

type DBTemplate struct {
	Name string
}

func loadFixture(t *testing.T, fixture string) string {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatalf("problems recovering caller information")
	}

	content, err := ioutil.ReadFile(filepath.Join(filepath.Dir(filename), "fixtures", fixture))
	if err != nil {
		t.Fatal(err)
	}

	return string(content)
}

func randName() string {
	return fmt.Sprintf("%s_%s", app.Name, strconv.Itoa(int(time.Now().UnixNano()+int64(os.Getpid()))))
}

func setupPostgres(t *testing.T) (*db.DSN, func()) {
	randName := randName()

	dsn, err := db.NewDSN(fmt.Sprintf("postgres://localhost/%s?sslmode=disable", randName))
	if err != nil {
		t.Fatal(err)
	}

	tmpfile, err := ioutil.TempFile("", "db_test")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		err = os.Remove(tmpfile.Name())
		if err != nil {
			t.Fatal(err)
		}
	}()

	tmpl := template.Must(template.New("ddlPsql").Parse(loadFixture(t, "postgres.sql")))

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
	randName := randName()

	dsn, err := db.NewDSN(fmt.Sprintf("mysql://root:root@localhost/%s", randName))
	if err != nil {
		t.Fatal(err)
	}

	tmpl := template.Must(template.New("ddl").Parse(loadFixture(t, "mysql.sql")))

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

func SetupDB(driver db.DriverType, t *testing.T) (dsn *db.DSN, cleanup func()) {
	switch driver {
	case db.PostgresDriver:
		dsn, cleanup = setupPostgres(t)
	case db.MySQLDriver:
		dsn, cleanup = setupMysql(t)
	}

	return dsn, cleanup
}

func InitApp(t *testing.T, vars map[string]string) *app.App {
	for k, v := range vars {
		os.Setenv(fmt.Sprintf("%s_%s", strings.ToUpper(app.Name), k), v)
	}

	return app.NewApp()
}
