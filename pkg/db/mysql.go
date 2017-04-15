package db

type MySQL struct {
}

func (my *MySQL) TablesQuery() string {
	return `SHOW TABLES`
}

func NewMySQLDialect() *MySQL {
	return &MySQL{}
}
