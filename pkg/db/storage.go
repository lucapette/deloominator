package db

import (
	"github.com/Sirupsen/logrus"
	"github.com/mattes/migrate"
	_ "github.com/mattes/migrate/database/mysql" // migration drivers
	_ "github.com/mattes/migrate/database/postgres"
	bindata "github.com/mattes/migrate/source/go-bindata"
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

// AutoUpgrade runs embedded migration
func (s *Storage) AutoUpgrade() error {
	resource := bindata.Resource(AssetNames(), func(name string) ([]byte, error) {
		return Asset(name)
	})

	driver, err := bindata.WithInstance(resource)
	if err != nil {
		return err
	}

	m, err := migrate.NewWithSourceInstance("go-bindata", driver, s.ds.dialect.ConnectionString())
	if err != nil {
		return err
	}
	err = m.Up()
	if err == migrate.ErrNoChange {
		return nil
	}

	return err
}

// Close closes deloominator internal storage
func (s *Storage) Close() {
	if err := s.ds.Close(); err != nil {
		logrus.Fatalf("could not close storage: %v", err)
	}
}
