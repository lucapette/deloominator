package db

import (
	"github.com/Sirupsen/logrus"
	_ "github.com/go-sql-driver/mysql" // import db drivers
	_ "github.com/lib/pq"
)

// Storage is a struct for deloominator own database
type Storage struct {
	ds *DataSource
}

// NewStorage initializes deloominator own storage using source
func NewStorage(source string) (*Storage, error) {
	dataSource, err := NewDataSource(source)
	if err != nil {
		return nil, err
	}

	return &Storage{ds: dataSource}, nil
}

func (s *Storage) AutoUpgrade() error {
	names, err := s.ds.Tables()
	if err != nil {
		return err
	}

	found := false
	for _, name := range names {
		if name == "migrations" {
			found = true
			break
		}
	}

	if !found {
		_, err := s.ds.Query("CREATE TABLE migrations (version varchar(255));")
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Storage) Close() {
	if err := s.ds.Close(); err != nil {
		logrus.Fatalf("could not close storage: %v", err)
	}
}
