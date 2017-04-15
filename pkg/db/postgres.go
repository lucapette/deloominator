package db

type Postgres struct {
}

func (pg *Postgres) TablesQuery() string {
	return `SELECT tablename FROM pg_tables WHERE schemaname = 'public' ORDER BY 1`
}

func NewPostgresDialect() *Postgres {
	return &Postgres{}
}
