package db

import (
	"fmt"
	"regexp"
	"strconv"
)

// support for https://en.wikipedia.org/wiki/Data_source_name

type DSN struct {
	Driver   string
	Username string
	Pass     string
	Host     string
	Port     int
	DBName   string
	Options  string
}

var validDSN = regexp.MustCompile(`(?P<driver>.+)://(?P<cred>(?P<username>.+):(?P<pass>.+)@)?(?P<host>.+):(?P<port>.+)?/(?P<db_name>[^?]+)\??(?P<options>.+)?`)

func ParseDSN(source string) (DSN, error) {
	fields := validDSN.FindStringSubmatch(source)

	if fields == nil {
		return DSN{}, fmt.Errorf("%s is not a valid Data Source Name", source)
	}

	groupedMatches := make(map[string]string, len(fields))
	for i, match := range fields {
		groupedMatches[validDSN.SubexpNames()[i]] = match
	}

	port, err := strconv.Atoi(groupedMatches["port"])

	if err != nil {
		return DSN{}, fmt.Errorf("%s is not a valid Data Source Name. %s is not a valid port", source, groupedMatches["port"])
	}

	return DSN{
		Driver:   groupedMatches["driver"],
		Username: groupedMatches["username"],
		Pass:     groupedMatches["pass"],
		Host:     groupedMatches["host"],
		Port:     port,
		DBName:   groupedMatches["db_name"],
		Options:  groupedMatches["options"],
	}, nil
}
