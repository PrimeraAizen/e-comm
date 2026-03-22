package repository

import (
	"context"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"

	"github.com/PrimeraAizen/e-comm/internal/domain"
	postgres "github.com/PrimeraAizen/e-comm/pkg/adapter"
)

type ProfileRepository interface {
	Create(ctx context.Context, profile *domain.Profile) error
	GetByUserID(ctx context.Context, userID string) (*domain.Profile, error)
	Update(ctx context.Context, profile *domain.Profile) error
	Delete(ctx context.Context, userID string) error
}

type profileRepository struct {
	db *postgres.Postgres
}

func NewProfileRepository(db *postgres.Postgres) ProfileRepository {
	return &profileRepository{db: db}
}

func (r *profileRepository) Create(ctx context.Context, profile *domain.Profile) error {
	now := time.Now()
	query, args, err := r.db.Builder.
		Insert("profiles").
		Columns("user_id", "first_name", "last_name", "middle_name", "date_of_birth", "gender", "phone", "address", "city", "country", "postal_code", "created_at", "updated_at").
		Values(profile.UserID, profile.FirstName, profile.LastName, profile.MiddleName, profile.DateOfBirth, profile.Gender, profile.Phone, profile.Address, profile.City, profile.Country, profile.PostalCode, now, now).
		Suffix("RETURNING id").
		ToSql()
	if err != nil {
		return err
	}
	return r.db.Pool.QueryRow(ctx, query, args...).Scan(&profile.ID)
}

func (r *profileRepository) GetByUserID(ctx context.Context, userID string) (*domain.Profile, error) {
	query, args, err := r.db.Builder.
		Select("id", "user_id", "first_name", "last_name", "middle_name", "date_of_birth", "gender", "phone", "address", "city", "country", "postal_code", "created_at", "updated_at").
		From("profiles").
		Where(squirrel.Eq{"user_id": userID}).
		ToSql()
	if err != nil {
		return nil, err
	}

	var p domain.Profile
	err = r.db.Pool.QueryRow(ctx, query, args...).Scan(
		&p.ID, &p.UserID, &p.FirstName, &p.LastName, &p.MiddleName,
		&p.DateOfBirth, &p.Gender, &p.Phone, &p.Address, &p.City,
		&p.Country, &p.PostalCode, &p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}

	return &p, nil
}

func (r *profileRepository) Update(ctx context.Context, profile *domain.Profile) error {
	now := time.Now()
	query, args, err := r.db.Builder.
		Update("profiles").
		Set("first_name", profile.FirstName).
		Set("last_name", profile.LastName).
		Set("middle_name", profile.MiddleName).
		Set("date_of_birth", profile.DateOfBirth).
		Set("gender", profile.Gender).
		Set("phone", profile.Phone).
		Set("address", profile.Address).
		Set("city", profile.City).
		Set("country", profile.Country).
		Set("postal_code", profile.PostalCode).
		Set("updated_at", now).
		Where(squirrel.Eq{"user_id": profile.UserID}).
		ToSql()
	if err != nil {
		return err
	}

	rows, err := r.db.Pool.Exec(ctx, query, args...)
	if err != nil {
		return err
	}
	if rows.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *profileRepository) Delete(ctx context.Context, userID string) error {
	query, args, err := r.db.Builder.
		Delete("profiles").
		Where(squirrel.Eq{"user_id": userID}).
		ToSql()
	if err != nil {
		return err
	}

	rows, err := r.db.Pool.Exec(ctx, query, args...)
	if err != nil {
		return err
	}
	if rows.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}
