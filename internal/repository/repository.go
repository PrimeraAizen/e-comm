package repository

import (
	postgres "github.com/PrimeraAizen/e-comm/pkg/adapter"
)

type Repository struct {
	User UserRepository
	// Profile     ProfileRepository
	// Product     ProductRepository
	// Interaction InteractionRepository
}

func NewRepositories(db *postgres.Postgres) *Repository {
	return &Repository{
		User: NewUserRepository(db),
		// Profile:     NewProfileRepository(db),
		// Product:     NewProductRepository(db),
		// Interaction: NewInteractionRepository(db),
	}
}
