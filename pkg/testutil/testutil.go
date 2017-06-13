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
	"github.com/lucapette/deloominator/pkg/config"
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
	return fmt.Sprintf("%s_%s", config.BinaryName, strconv.Itoa(int(time.Now().UnixNano()+int64(os.Getpid()))))
}

func SetupPG(t *testing.T) (string, func()) {
	randName := randName()

	cfg, err := getTestConfig()
	if err != nil {
		t.Fatalf("could not get test config: %v", err)
	}

	dsn := fmt.Sprintf("postgres://%s:%s@localhost/%s?sslmode=disable", cfg.PGUser, cfg.PGPass, randName)

	tmpfile, err := ioutil.TempFile("", "db_test")
	if err != nil {
		t.Fatalf("could not create temp file: %v", err)
	}
	defer func() {
		err = os.Remove(tmpfile.Name())
		if err != nil {
			t.Fatalf("could not remove temp file: %v", err)
		}
	}()

	ParseFixture(t, tmpfile, "postgres.sql", DBTemplate{Name: randName})

	if output, err := exec.Command("psql", "-f", tmpfile.Name()).CombinedOutput(); err != nil {
		t.Fatalf("%s\ncould not run psql: %v", output, err)
	}

	return dsn, func() {
		db, err := sql.Open("postgres",
			fmt.Sprintf("postgres://%s:%s@localhost/?sslmode=disable", cfg.PGUser, cfg.PGPass))
		if err != nil {
			t.Fatalf("could not open postgres database: %v", err)
		}

		if _, err = db.Exec(fmt.Sprintf("DROP DATABASE %s;", randName)); err != nil {
			t.Fatalf("could not drop database %s: %v", randName, err)
		}

		if err := db.Close(); err != nil {
			t.Fatalf("could not close database %s: %v", randName, err)
		}
	}
}

func SetupMySQL(t *testing.T) (string, func()) {
	randName := randName()

	cfg, err := getTestConfig()
	if err != nil {
		t.Fatalf("could not get test config: %v", err)
	}

	dsn := fmt.Sprintf("mysql://%s:%s@/%s", cfg.MysqlUser, cfg.MysqlPass, randName)

	query := bytes.Buffer{}
	ParseFixture(t, &query, "mysql.sql", DBTemplate{Name: randName})

	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@/?multiStatements=true", cfg.MysqlUser, cfg.MysqlPass))
	if err != nil {
		t.Fatalf("could not open mysql database %s: %v", randName, err)
	}

	if _, err := db.Exec(query.String()); err != nil {
		t.Fatalf("could not execute %s on mysql database %s: %v", query.String(), randName, err)
	}

	return dsn, func() {
		if _, err := db.Exec(fmt.Sprintf("DROP DATABASE %s;", randName)); err != nil {
			t.Fatalf("could not drop mysql database %s: %v", randName, err)
		}

		if err = db.Close(); err != nil {
			t.Fatalf("could not close mysql database %s: %v", randName, err)
		}
	}
}

func LoadData(t *testing.T, ds *db.DataSource, table string, data db.QueryResult) {
	query := &bytes.Buffer{}

	columns := make([]string, len(data.Columns))
	for i, col := range data.Columns {
		columns[i] = col.Name
	}

	rows := make([]string, len(data.Rows))
	for i, r := range data.Rows {
		row := &bytes.Buffer{}

		cells := make([]string, len(r))
		for i, c := range r {
			cells[i] = fmt.Sprintf("'%s'", c.Value)
		}

		fmt.Fprintf(row, "(%s)", strings.Join(cells, ","))

		rows[i] = row.String()
	}

	fmt.Fprintf(query, "insert into %s (%s) values %s",
		table,
		strings.Join(columns, ","),
		strings.Join(rows, ","),
	)

	if _, err := ds.Query(query.String()); err != nil {
		t.Fatalf("could not execute query %s on db %s: %v", query.String(), ds.Name(), err)
	}
}

func SetupDataSources(t *testing.T) (db.DataSources, func()) {
	dsnPG, cleanupPG := SetupPG(t)
	dsnMySQL, cleanupMySQL := SetupMySQL(t)

	sources := []string{dsnMySQL, dsnPG}

	dataSources, err := db.NewDataSources(sources)
	if err != nil {
		t.Fatalf("could not create DataSources from %v: %v", sources, err)
	}

	return dataSources, func() {
		dataSources.Close()
		cleanupPG()
		cleanupMySQL()
	}
}

// Diff returns a string of differences between expected and actual if any
func Diff(expected, actual interface{}) []string {
	return pretty.Diff(expected, actual)
}
