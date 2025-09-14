package service

import "e-comm/internal/repository"

type Example interface {
	ExampleMethod() error
}

type ExampleServiceDeps struct {
	repo repository.Example
}

func NewExampleService(repo *repository.ExampleRepository) *ExampleServiceDeps {
	return &ExampleServiceDeps{
		repo: repo,
	}
}

func (e *ExampleServiceDeps) ExampleMethod() error {
	return e.repo.ExampleMethod()
}
