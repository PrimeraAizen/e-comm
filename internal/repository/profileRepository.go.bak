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

type ProfileRepository interface {
	Create(ctx context.Context, profile *domain.Profile) error
	GetByUserID(ctx context.Context, userID int) (*domain.Profile, error)
	Update(ctx context.Context, profile *domain.Profile) error
	Delete(ctx context.Context, userID int) error
}

type profileRepository struct {
	db *mongodb.MongoDB
}

func NewProfileRepository(db *mongodb.MongoDB) ProfileRepository {
	return &profileRepository{db: db}
}

// getNextID gets the next profile ID from the counter
func (r *profileRepository) getNextID(ctx context.Context) (int, error) {
	collection := r.db.Collection("counters")

	filter := bson.M{"_id": "profile_id"}
	update := bson.M{"$inc": bson.M{"seq": 1}}

	var result struct {
		Seq int `bson:"seq"`
	}

	opts := options.FindOneAndUpdate().
		SetReturnDocument(options.After).
		SetUpsert(true)

	err := collection.FindOneAndUpdate(
		ctx,
		filter,
		update,
		opts,
	).Decode(&result)

	if err != nil {
		return 0, fmt.Errorf("get next profile id: %w", err)
	}

	return result.Seq, nil
}

// Create creates a new profile
func (r *profileRepository) Create(ctx context.Context, profile *domain.Profile) error {
	collection := r.db.Collection("profiles")

	// Get next ID
	id, err := r.getNextID(ctx)
	if err != nil {
		return err
	}
	profile.ID = id

	// Set timestamps
	now := time.Now()
	profile.CreatedAt = now
	profile.UpdatedAt = now

	_, err = collection.InsertOne(ctx, profile)
	if err != nil {
		return fmt.Errorf("create profile: %w", err)
	}

	return nil
}

// GetByUserID gets profile by user ID
func (r *profileRepository) GetByUserID(ctx context.Context, userID int) (*domain.Profile, error) {
	collection := r.db.Collection("profiles")

	var profile domain.Profile
	err := collection.FindOne(ctx, bson.M{"user_id": userID}).Decode(&profile)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("get profile by user id: %w", err)
	}

	return &profile, nil
}

// Update updates a profile
func (r *profileRepository) Update(ctx context.Context, profile *domain.Profile) error {
	collection := r.db.Collection("profiles")

	profile.UpdatedAt = time.Now()

	update := bson.M{
		"$set": bson.M{
			"first_name":    profile.FirstName,
			"last_name":     profile.LastName,
			"middle_name":   profile.MiddleName,
			"date_of_birth": profile.DateOfBirth,
			"gender":        profile.Gender,
			"phone":         profile.Phone,
			"address":       profile.Address,
			"city":          profile.City,
			"country":       profile.Country,
			"postal_code":   profile.PostalCode,
			"updated_at":    profile.UpdatedAt,
		},
	}

	result, err := collection.UpdateOne(ctx, bson.M{"user_id": profile.UserID}, update)
	if err != nil {
		return fmt.Errorf("update profile: %w", err)
	}

	if result.MatchedCount == 0 {
		return domain.ErrNotFound
	}

	return nil
}

// Delete deletes a profile
func (r *profileRepository) Delete(ctx context.Context, userID int) error {
	collection := r.db.Collection("profiles")

	result, err := collection.DeleteOne(ctx, bson.M{"user_id": userID})
	if err != nil {
		return fmt.Errorf("delete profile: %w", err)
	}

	if result.DeletedCount == 0 {
		return domain.ErrNotFound
	}

	return nil
}
