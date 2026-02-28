package repository

import (
	"context"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"

	"github.com/PrimeraAizen/e-comm/internal/domain"
	postgres "github.com/PrimeraAizen/e-comm/pkg/adapter"
)

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	GetByID(ctx context.Context, id string) (*domain.User, error)
	// Update(ctx context.Context, user *domain.User) error
	UpdateLastLogin(ctx context.Context, id string) error
}

type userRepository struct {
	db *postgres.Postgres
}

func NewUserRepository(db *postgres.Postgres) *userRepository {
	return &userRepository{
		db: db,
	}
}

func (repo *userRepository) Create(ctx context.Context, user *domain.User) error {
	now := time.Now()
	query, args, err := repo.db.Builder.
		Insert("users").
		Columns("email", "password_hash", "created_at", "updated_at").
		Values(&user.Email, &user.Password, now, now).
		ToSql()
	if err != nil {
		return err
	}
	_, err = repo.db.Pool.Exec(ctx, query, args...)
	if err != nil {
		return err
	}
	return nil
}

func (repo *userRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	res := domain.User{}
	query, args, err := repo.db.Builder.Select("id", "email", "password_hash", "created_at", "updated_at", "last_login_at").From("users").Where(squirrel.Eq{
		"email": email,
	}).ToSql()
	if err != nil {
		return nil, err
	}
	err = repo.db.Pool.QueryRow(ctx, query, args...).Scan(&res.ID, &res.Email, &res.Password, &res.CreatedAt, &res.UpdatedAt, &res.LastLoginAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return &res, nil
}

func (repo *userRepository) GetByID(ctx context.Context, id string) (*domain.User, error) {
	res := domain.User{}
	query, args, err := repo.db.Builder.Select("id", "email", "password_hash", "created_at", "updated_at", "last_login_at").From("users").Where(squirrel.Eq{
		"id": id,
	}).ToSql()
	if err != nil {
		return nil, err
	}
	err = repo.db.Pool.QueryRow(ctx, query, args...).Scan(&res.ID, &res.Email, &res.Password, &res.CreatedAt, &res.UpdatedAt, &res.LastLoginAt)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (repo *userRepository) UpdateLastLogin(ctx context.Context, id string) error {
	now := time.Now()
	match := squirrel.Eq{
		"id": id,
	}
	query, args, err := repo.db.Builder.Update("users").Set("last_login_at", now).Where(match).ToSql()
	if err != nil {
		return err
	}
	rows, err := repo.db.Pool.Exec(ctx, query, args...)
	if err != nil {
		return err
	}
	if rows.RowsAffected() == 0 {
		return domain.ErrInvalidCredentials
	}
	return nil
}
