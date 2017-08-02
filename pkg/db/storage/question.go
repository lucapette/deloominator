package storage

import (
	"fmt"

	"github.com/satori/go.uuid"
)

// Question stores q&a
type Question struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Query string `json:"query"`
}

func (s *Storage) FindQuestion(id string) (*Question, error) {
	q := fmt.Sprintf(`SELECT title, query FROM questions WHERE id='%s'`, id)

	row := s.dataSource.QueryRow(q)

	var title, query string
	err := row.Scan(&title, &query)
	if err != nil {
		return nil, err
	}
	return &Question{ID: id, Title: title, Query: query}, err
}

func (s *Storage) InsertQuestion(title, query string) (question *Question, err error) {
	id := uuid.NewV4().String()
	q := fmt.Sprintf(
		`INSERT INTO QUESTIONS (id, title, query) VALUES('%s', '%s', '%s')`,
		id,
		title,
		query,
	)
	_, err = s.dataSource.Exec(q)

	if err != nil {
		return nil, err
	}

	return &Question{ID: id, Title: title, Query: query}, err
}
