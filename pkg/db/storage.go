package db

import (
	"net/url"

	_ "github.com/go-sql-driver/mysql" // import db drivers
	_ "github.com/lib/pq"
)

// Storage is a struct for deloominator own database
type Storage struct {
	dialect Dialect
}

// NewStorage initializes deloominator own storage using source
func NewStorage(source string) (*Storage, error) {
	url, err := url.Parse(source)
	if err != nil {
		return nil, err
	}
	dialect, err := NewDialect(url)
	if err != nil {
		return nil, err
	}
	return &Storage{dialect: dialect}, nil
}

func (s *Storage) AutoUpgrade() error {
	return nil
}
