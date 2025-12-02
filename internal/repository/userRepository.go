package repository

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/PrimeraAizen/e-comm/internal/domain"
	mongodb "github.com/PrimeraAizen/e-comm/pkg/adapter/mongodb"
)

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	GetByID(ctx context.Context, id int) (*domain.User, error)
	Update(ctx context.Context, user *domain.User) error
	UpdateLastLogin(ctx context.Context, id int) error
}

type userRepository struct {
	db *mongodb.MongoDB
}

func NewUserRepository(db *mongodb.MongoDB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *domain.User) error {
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	user.Status = "active"

	collection := r.db.Collection("users")

	// Get the next ID
	nextID, err := r.getNextID(ctx)
	if err != nil {
		return fmt.Errorf("get next ID: %w", err)
	}
	user.ID = nextID

	_, err = collection.InsertOne(ctx, user)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return fmt.Errorf("user with this email already exists: %w", err)
		}
		return fmt.Errorf("create user: %w", err)
	}

	return nil
}

// getNextID gets the next auto-increment ID for users
func (r *userRepository) getNextID(ctx context.Context) (int, error) {
	collection := r.db.Collection("users")

	// Find the maximum ID
	opts := options.Find().SetSort(bson.M{"_id": -1}).SetLimit(1)
	cursor, err := collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return 0, err
	}
	defer cursor.Close(ctx)

	if cursor.Next(ctx) {
		var result domain.User
		if err := cursor.Decode(&result); err != nil {
			return 0, err
		}
		return result.ID + 1, nil
	}

	// If no users exist, start from 1
	return 1, nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	collection := r.db.Collection("users")

	var user domain.User
	err := collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("get user by email: %w", err)
	}

	return &user, nil
}

func (r *userRepository) GetByID(ctx context.Context, id int) (*domain.User, error) {
	collection := r.db.Collection("users")

	var user domain.User
	err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("get user by id: %w", err)
	}

	return &user, nil
}

func (r *userRepository) Update(ctx context.Context, user *domain.User) error {
	collection := r.db.Collection("users")

	user.UpdatedAt = time.Now()

	update := bson.M{
		"$set": bson.M{
			"email":         user.Email,
			"password_hash": user.Password,
			"status":        user.Status,
			"updated_at":    user.UpdatedAt,
		},
	}

	result, err := collection.UpdateOne(ctx, bson.M{"_id": user.ID}, update)
	if err != nil {
		return fmt.Errorf("update user: %w", err)
	}

	if result.MatchedCount == 0 {
		return domain.ErrNotFound
	}

	return nil
}

func (r *userRepository) UpdateLastLogin(ctx context.Context, id int) error {
	collection := r.db.Collection("users")

	now := time.Now()
	update := bson.M{
		"$set": bson.M{
			"last_login_at": now,
		},
	}

	result, err := collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	if err != nil {
		return fmt.Errorf("update last login: %w", err)
	}

	if result.MatchedCount == 0 {
		return domain.ErrNotFound
	}

	return nil
}
