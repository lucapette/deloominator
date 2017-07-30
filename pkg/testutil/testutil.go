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

type TestFile struct {
	t    *testing.T
	name string
	dir  string
}

func NewFixture(t *testing.T, name string) *TestFile {
	return &TestFile{t: t, name: name, dir: "fixtures"}
}

func NewGoldenFile(t *testing.T, name string) *TestFile {
	return &TestFile{t: t, name: name, dir: "golden"}
}

func (tf *TestFile) path() string {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		tf.t.Fatal("problems recovering caller information")
	}

	return filepath.Join(filepath.Dir(filename), tf.dir, tf.name)
}

func (tf *TestFile) Write(content string) {
	err := ioutil.WriteFile(tf.path(), []byte(content), 0644)
	if err != nil {
		tf.t.Fatalf("could not write %s: %v", tf.name, err)
	}
}

func (tf *TestFile) Load() string {
	content, err := ioutil.ReadFile(tf.path())
	if err != nil {
		tf.t.Fatalf("could not read file %s: %v", tf.name, err)
	}

	return string(content)
}

func (tf *TestFile) Parse(w io.Writer, data string) {
	tmpl := template.Must(template.New(tf.name).Parse(tf.Load()))

	err := tmpl.Execute(w, data)
	if err != nil {
		tf.t.Fatalf("could not execute template %s: %v", tf.name, err)
	}
}

func randDBName() string {
	return fmt.Sprintf("%s_%s", config.BinaryName, strconv.Itoa(int(time.Now().UnixNano()+int64(os.Getpid()))))
}

func createPG(t *testing.T) (string, func()) {
	randName := randDBName()

	cfg, err := getTestConfig()
	if err != nil {
		t.Fatalf("could not get test config: %v", err)
	}

	db, err := sql.Open("postgres", fmt.Sprintf("postgres://%s:%s@localhost/?sslmode=disable", cfg.PGUser, cfg.PGPass))
	if err != nil {
		t.Fatalf("could not open postgres database: %v", err)
	}

	if _, err = db.Exec(fmt.Sprintf("CREATE DATABASE %s;", randName)); err != nil {
		t.Fatalf("could not create database %s: %v", randName, err)
	}

	return randName, func() {
		if _, err = db.Exec(fmt.Sprintf("DROP DATABASE %s;", randName)); err != nil {
			t.Fatalf("could not drop database %s: %v", randName, err)
		}

		if err := db.Close(); err != nil {
			t.Fatalf("could not close database %s: %v", randName, err)
		}
	}
}

func pgDSN(cfg *testConfig, name string) string {
	return fmt.Sprintf("postgres://%s:%s@localhost/%s?sslmode=disable", cfg.PGUser, cfg.PGPass, name)
}

func setupPG(t *testing.T) (string, func()) {
	name, cleanup := createPG(t)
	cfg, err := getTestConfig()
	if err != nil {
		t.Fatalf("could not get test config: %v", err)
	}

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

	NewFixture(t, "postgres.sql").Parse(tmpfile, name)

	if output, err := exec.Command("psql", "-f", tmpfile.Name()).CombinedOutput(); err != nil {
		t.Fatalf("%s\ncould not run psql: %v", output, err)
	}

	return pgDSN(cfg, name), cleanup
}

func createMySQL(t *testing.T) (string, func()) {
	randName := randDBName()
	cfg, err := getTestConfig()
	if err != nil {
		t.Fatalf("could not get test config: %v", err)
	}

	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@/?multiStatements=true", cfg.MysqlUser, cfg.MysqlPass))
	if err != nil {
		t.Fatalf("could not open mysql database %s: %v", randName, err)
	}

	if _, err := db.Exec(fmt.Sprintf("CREATE DATABASE %s;", randName)); err != nil {
		t.Fatalf("could not create mysql database %s: %v", randName, err)
	}

	return randName, func() {
		if _, err = db.Exec(fmt.Sprintf("DROP DATABASE %s;", randName)); err != nil {
			t.Fatalf("could not drop database %s: %v", randName, err)
		}

		if err := db.Close(); err != nil {
			t.Fatalf("could not close database %s: %v", randName, err)
		}
	}
}

func mysqlDSN(cfg *testConfig, name string) string {
	return fmt.Sprintf("mysql://%s:%s@/%s", cfg.MysqlUser, cfg.MysqlPass, name)
}

func setupMySQL(t *testing.T) (string, func()) {
	name, cleanup := createMySQL(t)
	cfg, err := getTestConfig()
	if err != nil {
		t.Fatalf("could not get test config: %v", err)
	}

	query := &bytes.Buffer{}
	NewFixture(t, "mysql.sql").Parse(query, name)

	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@/?multiStatements=true", cfg.MysqlUser, cfg.MysqlPass))
	if err != nil {
		t.Fatalf("could not open mysql database %s: %v", name, err)
	}

	if _, err := db.Exec(query.String()); err != nil {
		t.Fatalf("could not execute %s on mysql database %s: %v", query.String(), name, err)
	}

	return mysqlDSN(cfg, name), func() {
		if err = db.Close(); err != nil {
			t.Fatalf("could not close mysql database %s: %v", name, err)
		}
		cleanup()
	}
}

func LoadDataFromFixture(t *testing.T, ds *db.DataSource, fixture string) {
	query := NewFixture(t, fixture).Load()

	if _, err := ds.Exec(query); err != nil {
		t.Fatalf("could not execute query %s on %s: %v", query, ds.Name(), err)
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

func CreateDSN(t *testing.T) []string {
	cfg, err := getTestConfig()
	if err != nil {
		t.Fatalf("could not get test config: %v", err)
	}
	return []string{pgDSN(cfg, randDBName()), mysqlDSN(cfg, randDBName())}
}

func CreateDataSources(t *testing.T) ([]string, func()) {
	cfg, err := getTestConfig()
	if err != nil {
		t.Fatalf("could not get test config: %v", err)
	}
	pgName, cleanupPG := createPG(t)
	mysqlName, cleanupMySQL := createMySQL(t)

	return []string{pgDSN(cfg, pgName), mysqlDSN(cfg, mysqlName)}, func() {
		cleanupPG()
		cleanupMySQL()
	}
}

func SetupDataSources(t *testing.T) (db.DataSources, func()) {
	dsnPG, cleanupPG := setupPG(t)
	dsnMySQL, cleanupMySQL := setupMySQL(t)

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
