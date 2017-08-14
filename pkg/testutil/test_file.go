package testutil

import (
	"bytes"
	"html/template"
	"io"
	"io/ioutil"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// TestFile represents a test file
type TestFile struct {
	t    *testing.T
	name string
	dir  string
}

// NewFixture returns a TestFile from a given fixture
func NewFixture(t *testing.T, fixture string) *TestFile {
	return &TestFile{t: t, name: fixture, dir: "fixtures"}
}

// NewGoldenFile returns a TestFile from a given golden name file
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

// Load returns the content of a TestFile
func (tf *TestFile) Load() string {
	content, err := ioutil.ReadFile(tf.path())
	if err != nil {
		tf.t.Fatalf("could not read file %s: %v", tf.name, err)
	}

	return string(content)
}

// Parse parses the content of a TestFile
func (tf *TestFile) Parse(w io.Writer, data string) {
	tmpl := template.Must(template.New(tf.name).Parse(tf.Load()))

	err := tmpl.Execute(w, data)
	if err != nil {
		tf.t.Fatalf("could not execute template %s: %v", tf.name, err)
	}
}

func (tf *TestFile) ParseOrUpdate(update bool, data, actual string) string {
	out := &bytes.Buffer{}

	tf.Parse(out, data)

	if update {
		tf.Write(strings.Replace(actual, data, "{{.}}", -1))
	}

	return strings.TrimSuffix(out.String(), "\n")
}
