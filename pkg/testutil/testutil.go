package testutil

import (
	"bytes"
	"database/sql"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/kr/pretty"
	"github.com/lucapette/deloominator/pkg/app"
	"github.com/lucapette/deloominator/pkg/db"
)

type DBTemplate struct {
	Name string
}

func fixturePath(t *testing.T, fixture string) string {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatalf("problems recovering caller information")
	}

	return filepath.Join(filepath.Dir(filename), "fixtures", fixture)
}

func WriteFixture(t *testing.T, fixture, content string) {
	err := ioutil.WriteFile(fixturePath(t, fixture), []byte(content), 0644)
	if err != nil {
		t.Fatal(err)
	}
}

func LoadFixture(t *testing.T, fixture string) string {
	content, err := ioutil.ReadFile(fixturePath(t, fixture))
	if err != nil {
		t.Fatal(err)
	}

	return string(content)
}

func ParseFixture(t *testing.T, w io.Writer, fixture string, data interface{}) {
	tmpl := template.Must(template.New(fixture).Parse(LoadFixture(t, fixture)))

	err := tmpl.Execute(w, data)
	if err != nil {
		t.Fatal(err)
	}
}

func randName() string {
	return fmt.Sprintf("%s_%s", app.Name, strconv.Itoa(int(time.Now().UnixNano()+int64(os.Getpid()))))
}

func setupPostgres(t *testing.T) (string, func()) {
	randName := randName()

	dsn := fmt.Sprintf("postgres://localhost/%s?sslmode=disable", randName)

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

	ParseFixture(t, tmpfile, "postgres.sql", DBTemplate{Name: randName})

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

func setupMysql(t *testing.T) (string, func()) {
	randName := randName()

	dsn := fmt.Sprintf("mysql://root:root@/%s", randName)

	var query bytes.Buffer
	ParseFixture(t, &query, "mysql.sql", DBTemplate{Name: randName})

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

func SetupDB(t *testing.T, driver db.DriverType) (dsn string, cleanup func()) {
	switch driver {
	case db.PostgresDriver:
		dsn, cleanup = setupPostgres(t)
	case db.MySQLDriver:
		dsn, cleanup = setupMysql(t)
	}

	return dsn, cleanup
}

func LoadData(t *testing.T, ds *db.DataSource, table string, result db.QueryResult) {
	query := bytes.NewBufferString(fmt.Sprintf("insert into %s (", table))

	columns := make([]string, len(result.Columns))
	for i, col := range result.Columns {
		columns[i] = col.Name
	}
	query.WriteString(strings.Join(columns, ","))

	query.WriteString(") values ")

	rows := make([]string, len(result.Rows))
	for i, r := range result.Rows {
		row := bytes.NewBufferString("(")

		cells := make([]string, len(r))
		for i, c := range r {
			cells[i] = fmt.Sprintf("'%s'", c.Value)
		}

		row.WriteString(strings.Join(cells, ","))

		row.WriteString(")")

		rows[i] = row.String()
	}
	query.WriteString(strings.Join(rows, ","))

	_, err := ds.Query(query.String())
	if err != nil {
		t.Fatal(err)
	}
}

func InitApp(t *testing.T, vars map[string]string) *app.App {
	for k, v := range vars {
		err := os.Setenv(fmt.Sprintf("%s_%s", strings.ToUpper(app.Name), k), v)

		if err != nil {
			t.Fatal(err)
		}
	}

	return app.NewApp()
}

func Diff(expected, actual interface{}) []string {
	return pretty.Diff(expected, actual)
}
