package storage

// Question stores q&a
type Question struct {
	ID         int    `json:"id"`
	Title      string `json:"title,omitempty"`
	Query      string `json:"query,omitempty"`
	DataSource string `json:"dataSource,omitempty"`
	Variables  string `json:"variables,omitempty"`
}

func (s *Storage) InsertQuestion(q *Question) (*Question, error) {
	if err := s.orm.Create(q).Error; err != nil {
		return nil, err
	}

	return q, nil
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
