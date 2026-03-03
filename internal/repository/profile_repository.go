package repository

import (
	"context"
	"time"

	"github.com/PrimeraAizen/e-comm/internal/domain"
	postgres "github.com/PrimeraAizen/e-comm/pkg/adapter"
)

type ProfileRepository interface {
	Create(ctx context.Context, profile *domain.Profile) error
	// GetByUserID(ctx context.Context, userID int) (*domain.Profile, error)
	// Update(ctx context.Context, profile *domain.Profile) error
	// Delete(ctx context.Context, userID int) error
}

type profileRepository struct {
	db postgres.Postgres
}

func (r *profileRepository) Create(ctx context.Context, profile *domain.Profile) error {
	now := time.Now()
	query, args, err := r.db.Builder.
		Insert("profiles").Columns("user_id", "first_name", "last_name", "middle_name", "date_of_birth", "gender", "phone", "address", "city", "country", "postal_code", "created_at", "updated_at").
		Values(&profile.UserID, &profile.FirstName, &profile.LastName, &profile.MiddleName, &profile.MiddleName, &profile.DateOfBirth, &profile.Gender, &profile.Phone, &profile.Address, &profile.City, &profile.Country, &profile.PostalCode, now, now).
		ToSql()
	if err != nil {
		return err
	}

	_, err = r.db.Pool.Exec(ctx, query, args...)
	if err != nil {
		return err
	}
	return nil
}
