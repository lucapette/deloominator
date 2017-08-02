package storage

import (
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/lucapette/deloominator/pkg/db"

	"github.com/mattes/migrate"
	_ "github.com/mattes/migrate/database/mysql" // migration drivers
	_ "github.com/mattes/migrate/database/postgres"
	bindata "github.com/mattes/migrate/source/go-bindata"
)

// Storage is a struct for deloominator own database
type Storage struct {
	source     string
	dataSource *db.DataSource
}

// NewStorage initializes deloominator own storage using source
func NewStorage(source string) (*Storage, error) {
	dataSource, err := db.NewDataSource(source)
	if err != nil {
		return nil, fmt.Errorf("could not create data source from %s: %v", source, err)
	}
	return &Storage{source: source, dataSource: dataSource}, nil
}

// AutoUpgrade runs embedded migration
func (s *Storage) AutoUpgrade() error {
	if err := s.dataSource.CreateDBIfNotExist(); err != nil {
		return fmt.Errorf("cannot create storage %s: %v", s.dataSource.Name(), err)
	}
	resource := bindata.Resource(db.AssetNames(), func(n string) ([]byte, error) { return db.Asset(n) })

	driver, err := bindata.WithInstance(resource)
	if err != nil {
		return fmt.Errorf("cannot create migration driver %s: %v", s.dataSource.Name(), err)
	}

	m, err := migrate.NewWithSourceInstance("go-bindata", driver, s.source)
	if err != nil {
		return fmt.Errorf("cannot create migration instance %s: %v", s.dataSource.Name(), err)

	}
	defer func() {
		sourceErr, dbErr := m.Close()
		if sourceErr != nil {
			logrus.Printf("cannot close migration instance %s: %v", s.dataSource.Name(), sourceErr)
		}
		if dbErr != nil {
			logrus.Printf("cannot close migration instance %s: %v", s.dataSource.Name(), sourceErr)
		}
	}()

	err = m.Up()
	if err == migrate.ErrNoChange {
		return nil
	}

	return err
}

func (s *Storage) Close() {
	s.dataSource.Close()
}
