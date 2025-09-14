package repository

import postgres "e-comm/pkg/adapter"

type Example interface {
	ExampleMethod() error
}

type ExampleRepository struct {
	pg *postgres.Postgres
}

func NewExampleRepository(pg *postgres.Postgres) *ExampleRepository {
	return &ExampleRepository{
		pg: pg,
	}
}

func (e *ExampleRepository) ExampleMethod() error {
	return nil
}
