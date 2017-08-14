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

// Question stores q&a
type Question struct {
	ID         int    `json:"id"`
	Title      string `json:"title,omitempty"`
	Query      string `json:"query,omitempty"`
	DataSource string `json:"dataSource,omitempty"`
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

func (s *Storage) InsertQuestion(title, query, dataSource string) (*Question, error) {
	question := &Question{Title: title, Query: query, DataSource: dataSource}

	if err := s.orm.Create(question).Error; err != nil {
		return nil, err
	}

	return question, nil
}

func (s *Storage) FindQuestion(id int) (*Question, error) {
	question := &Question{ID: id}

	if err := s.orm.Find(question).Error; err != nil {
		return nil, err
	}

	return question, nil
}

func (s *Storage) AllQuestions() (questions []*Question, err error) {
	if err := s.orm.Order("title").Find(&questions).Error; err != nil {
		return nil, err
	}

	return questions, nil
}

// Close closes the underlining DB connection
func (s *Storage) Close() {
	s.dataSource.Close()
	if err := s.orm.Close(); err != nil {
		logrus.Warnf("could not close orm: %v", err)
	}
}

func (s *Storage) String() string {
	return fmt.Sprintf("%s-%s", s.dataSource.DBName(), s.dataSource.DriverName())
}
