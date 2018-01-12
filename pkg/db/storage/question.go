package storage

import "time"

// Question stores q&a
type Question struct {
	ID         int       `json:"id"`
	Title      string    `json:"title,omitempty"`
	Query      string    `json:"query,omitempty"`
	DataSource string    `json:"dataSource,omitempty"`
	Variables  string    `json:"variables,omitempty"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

// InsertQuestion stores a question into the storage
func (s *Storage) InsertQuestion(q *Question) (*Question, error) {
	if err := s.orm.Create(q).Error; err != nil {
		return nil, err
	}

	return q, nil
}

// FindQuestion returns a single question using its id
func (s *Storage) FindQuestion(id int) (*Question, error) {
	question := &Question{ID: id}

	if err := s.orm.Find(question).Error; err != nil {
		return nil, err
	}

	return question, nil
}

// AllQuestions returns all the questions stored in the given storage ordered by created_at and title
func (s *Storage) AllQuestions() (questions []*Question, err error) {
	if err := s.orm.Order("created_at, title").Find(&questions).Error; err != nil {
		return nil, err
	}

	return questions, nil
}
