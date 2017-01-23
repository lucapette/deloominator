package db

import (
	"bytes"
	"fmt"
	"regexp"
	"strconv"
)

// support for https://en.wikipedia.org/wiki/Data_source_name
type DSN struct {
	Driver         string
	Username, Pass string
	Host           string
	Port           int
	DBName         string
	Options        string
}

type DriverType int

const (
	PostgresDriver DriverType = iota
	MySQLDriver               = iota
)

var validDSN = regexp.MustCompile(`(?P<driver>[^:]+)://(?P<cred>(?P<username>[^:]+):(?P<pass>[^@]+)@)?(?P<host>[^:]+)(?P<opt_port>:(?P<port>[^/]+))?/(?P<db_name>[^?]+)?\??(?P<options>.+)?`)

func NewDSN(source string) (ds *DSN, err error) {
	fields := validDSN.FindStringSubmatch(source)

	if fields == nil {
		return ds, fmt.Errorf("%s is not a valid Data Source Name", source)
	}

	groupedMatches := make(map[string]string, len(fields))
	for i, match := range fields {
		groupedMatches[validDSN.SubexpNames()[i]] = match
	}

	var port int

	if len(groupedMatches["port"]) > 0 {
		port, err = strconv.Atoi(groupedMatches["port"])

		if err != nil {
			return ds, fmt.Errorf("%s is not a valid Data Source Name. %s is not a valid port", source, groupedMatches["port"])
		}
	}

	return &DSN{
		Driver:   groupedMatches["driver"],
		Username: groupedMatches["username"],
		Pass:     groupedMatches["pass"],
		Host:     groupedMatches["host"],
		Port:     port,
		DBName:   groupedMatches["db_name"],
		Options:  groupedMatches["options"],
	}, nil
}

func (ds *DSN) String() string {
	return ds.Format()
}

func (ds *DSN) Format() (output string) {
	switch ds.Driver {
	case "postgres":
		output = postgresFormat(ds)
	case "mysql":
		output = mysqlFormat(ds)
	}

	return output
}

func postgresFormat(ds *DSN) string {
	buf := bytes.NewBufferString(ds.Driver)
	buf.WriteString("://")

	if len(ds.Username) > 0 {
		buf.WriteString(ds.Username)
		buf.WriteString(":")
		buf.WriteString(ds.Pass)
		buf.WriteString("@")
	}

	buf.WriteString(ds.Host)

	if ds.Port > 0 {
		buf.WriteString(":")
		buf.WriteString(strconv.Itoa(ds.Port))
	}

	buf.WriteString("/")
	buf.WriteString(ds.DBName)

	if len(ds.Options) > 0 {
		buf.WriteString("?")
		buf.WriteString(ds.Options)
	}

	return buf.String()
}

func mysqlFormat(ds *DSN) string {
	var buf bytes.Buffer

	if len(ds.Username) > 0 {
		buf.WriteString(ds.Username)

		if len(ds.Pass) > 0 {
			buf.WriteString(":")
			buf.WriteString(ds.Pass)
		}

		buf.WriteString("@")
	}

	if ds.Port > 0 {
		buf.WriteString(ds.Host)
		buf.WriteString(":")
		buf.WriteString(strconv.Itoa(ds.Port))
	}

	buf.WriteString("/")
	buf.WriteString(ds.DBName)

	if len(ds.Options) > 0 {
		buf.WriteString("?")
		buf.WriteString(ds.Options)
	}

	return buf.String()
}
