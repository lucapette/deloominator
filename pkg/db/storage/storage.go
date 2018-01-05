package storage

import (
	"errors"
	"fmt"

	"github.com/Sirupsen/logrus"

	"github.com/jinzhu/gorm"
	"github.com/lucapette/deloominator/pkg/db"
)

// Storage is a struct for deloominator own database
type Storage struct {
	dataSource *db.DataSource
	orm        *gorm.DB
}

// NewStorage initializes deloominator own storage using source
func NewStorage(source string) (*Storage, error) {
	dataSource, err := db.NewDataSource(source)
	if err != nil {
		return nil, fmt.Errorf("could not create data source from %s: %v", source, err)
	}

	return &Storage{dataSource: dataSource}, nil
}

// AutoUpgrade migrates the DB if needed
func (s *Storage) AutoUpgrade() error {
	if err := s.dataSource.CreateDBIfNotExist(); err != nil {
		return fmt.Errorf("cannot create storage %s: %v", s.dataSource.DBName(), err)
	}

	orm, err := gorm.Open(s.dataSource.DriverName(), s.dataSource.ConnectionString())
	if err != nil {
		return fmt.Errorf("could not create orm: %v", err)
	}
	s.orm = orm

	errs := s.orm.AutoMigrate(&Question{}).GetErrors()
	if len(errs) > 0 {
		for _, e := range errs {
			logrus.Printf("%v", e)
		}
		return errors.New("could not migrate db")
	}

	return nil
}

// Close closes the underlining DB connection
func (s *Storage) Close() {
	if s.dataSource != nil {
		s.dataSource.Close()
	}
	if s.orm != nil {
		if err := s.orm.Close(); err != nil {
			logrus.Warnf("could not close orm: %v", err)
		}
	}
}

func (s *Storage) String() string {
	return fmt.Sprintf("%s-%s", s.dataSource.DBName(), s.dataSource.DriverName())
}
