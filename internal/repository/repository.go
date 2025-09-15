package repository

import postgres "e-comm/pkg/adapter"

type Repository struct {
	Example Example
}

func NewRepositories(pg *postgres.Postgres) *Repository {
	return &Repository{
		Example: NewExampleRepository(pg),
	}
}
