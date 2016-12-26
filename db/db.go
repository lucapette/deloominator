package db

import (
	"bytes"
	"fmt"
	"regexp"
	"strconv"

	"database/sql"

	log "github.com/Sirupsen/logrus"
	_ "github.com/lib/pq"
)

type DataSource struct {
	Driver         string
	Username, Pass string
	Host           string
	Port           int
	DBName         string
	Options        string
}

// support for https://en.wikipedia.org/wiki/Data_source_name
var validDSN = regexp.MustCompile(`(?P<driver>[^:]+)://(?P<cred>(?P<username>[^:]+):(?P<pass>[^@]+)@)?(?P<host>[^:]+)(?P<opt_port>:(?P<port>[^/]+))?/(?P<db_name>[^?]+)\??(?P<options>.+)?`)

type DataSources map[string]*DataSource

func NewDataSource(source string) (ds *DataSource, err error) {
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

	ds = &DataSource{
		Driver:   groupedMatches["driver"],
		Username: groupedMatches["username"],
		Pass:     groupedMatches["pass"],
		Host:     groupedMatches["host"],
		Port:     port,
		DBName:   groupedMatches["db_name"],
		Options:  groupedMatches["options"],
	}

	return ds, nil
}

func NewDataSources(dataSources []string) (sources DataSources, err error) {
	sources = make(DataSources, len(dataSources))

	for _, source := range dataSources {
		ds, err := NewDataSource(source)
		if err != nil {
			return nil, err
		}

		sources[ds.DBName] = ds
	}

	return sources, nil
}

func (ds DataSource) String() string {
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

func (ds *DataSource) Tables() (tables []string, err error) {
	db, err := sql.Open(ds.Driver, fmt.Sprint(ds))
	if err != nil {
		return tables, err
	}

	err = db.Ping()
	if err != nil {
		return tables, err
	}

	rows, err := db.Query(`SELECT tablename FROM pg_tables where schemaname = 'public'`)

	if err != nil {
		return tables, err
	}

	defer rows.Close()

	for rows.Next() {
		var name string
		err = rows.Scan(&name)
		if err != nil {
			return tables, err
		}

		log.Info(name)
	}

	return tables, err
}
