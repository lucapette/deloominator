package db

import (
	"fmt"

	"github.com/Sirupsen/logrus"

	"github.com/mattes/migrate"
	_ "github.com/mattes/migrate/database/mysql" // migration drivers
	_ "github.com/mattes/migrate/database/postgres"
	bindata "github.com/mattes/migrate/source/go-bindata"
)

// Storage is a struct for deloominator own database
type Storage struct {
	source string
}

// NewStorage initializes deloominator own storage using source
func NewStorage(source string) (*Storage, error) {
	return &Storage{source: source}, nil
}

// AutoUpgrade runs embedded migration
func (s *Storage) AutoUpgrade() error {
	dataSource, err := NewDataSource(s.source)
	if err != nil {
		return fmt.Errorf("could not create data source from %s: %v", s.source, err)
	}
	defer dataSource.Close()

	if err := dataSource.CreateDBIfNotExist(); err != nil {
		return fmt.Errorf("cannot create storage %s: %v", dataSource.Name(), err)
	}
	resource := bindata.Resource(AssetNames(), func(n string) ([]byte, error) { return Asset(n) })

	driver, err := bindata.WithInstance(resource)
	if err != nil {
		return fmt.Errorf("cannot create migration driver %s: %v", dataSource.Name(), err)
	}

	m, err := migrate.NewWithSourceInstance("go-bindata", driver, s.source)
	if err != nil {
		return fmt.Errorf("cannot create migration instance %s: %v", dataSource.Name(), err)

	}
	defer func() {
		sourceErr, dbErr := m.Close()
		if sourceErr != nil {
			logrus.Printf("cannot close migration instance %s: %v", dataSource.Name(), sourceErr)
		}
		if dbErr != nil {
			logrus.Printf("cannot close migration instance %s: %v", dataSource.Name(), sourceErr)
		}
	}()

	err = m.Up()
	if err == migrate.ErrNoChange {
		return nil
	}

	return err
}
