package repository

import (
	"context"
	"errors"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/PrimeraAizen/e-comm/internal/domain"
	postgres "github.com/PrimeraAizen/e-comm/pkg/adapter"
)

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	GetByID(ctx context.Context, id string) (*domain.User, error)
	UpdatePasswordHash(ctx context.Context, id, passwordHash string) error
	UpdateStatus(ctx context.Context, id, status string) error
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
	query, args, err := repo.db.Builder.
		Insert("users").
		Columns("email", "password_hash").
		Values(&user.Email, &user.PasswordHash).
		Suffix("RETURNING id, status, created_at, updated_at").
		ToSql()

	if err != nil {
		return err
	}

	err = repo.db.Pool.QueryRow(ctx, query, args...).Scan(&user.ID, &user.Status, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return domain.ErrAlreadyExists
		}

		return err
	}

	return nil
}

func (repo *userRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	return repo.getOne(ctx, squirrel.And{squirrel.Eq{"email": email}, notDeleted()})
}

func (repo *userRepository) GetByID(ctx context.Context, id string) (*domain.User, error) {
	return repo.getOne(ctx, activeUser(id))
}

func (repo *userRepository) UpdatePasswordHash(ctx context.Context, id, passwordHash string) error {
	return repo.execUpdate(ctx, repo.updateUser(id).Set("password_hash", passwordHash))
}

func (repo *userRepository) UpdateStatus(ctx context.Context, id, status string) error {
	return repo.execUpdate(ctx, repo.updateUser(id).Set("status", status))
}

func (repo *userRepository) UpdateLastLogin(ctx context.Context, id string) error {
	return repo.execUpdate(ctx, repo.db.Builder.
		Update("users").
		Set("last_login_at", squirrel.Expr("now()")).
		Where(activeUser(id)))
}

func notDeleted() squirrel.NotEq {
	return squirrel.NotEq{"status": domain.StatusUserInactive}
}

func activeUser(id string) squirrel.And {
	return squirrel.And{squirrel.Eq{"id": id}, notDeleted()}
}

func (repo *userRepository) getOne(ctx context.Context, where squirrel.Sqlizer) (*domain.User, error) {
	query, args, err := repo.db.Builder.
		Select("id", "email", "password_hash", "status", "created_at", "updated_at", "last_login_at").
		From("users").
		Where(where).
		ToSql()
	if err != nil {
		return nil, err
	}

	var u domain.User
	err = repo.db.Pool.QueryRow(ctx, query, args...).
		Scan(&u.ID, &u.Email, &u.PasswordHash, &u.Status, &u.CreatedAt, &u.UpdatedAt, &u.LastLoginAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (repo *userRepository) updateUser(id string) squirrel.UpdateBuilder {
	return repo.db.Builder.
		Update("users").
		Set("updated_at", squirrel.Expr("now()")).
		Where(activeUser(id))
}

func (repo *userRepository) execUpdate(ctx context.Context, b squirrel.UpdateBuilder) error {
	query, args, err := b.ToSql()
	if err != nil {
		return err
	}
	rows, err := repo.db.Pool.Exec(ctx, query, args...)
	if err != nil {
		return err
	}
	if rows.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}
