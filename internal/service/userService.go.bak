package service

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/PrimeraAizen/e-comm/internal/domain"
	"github.com/PrimeraAizen/e-comm/internal/repository"
)

type UserService interface {
	GetProfile(ctx context.Context, userID int) (*domain.User, *domain.Profile, error)
	UpdateProfile(ctx context.Context, userID int, profileData *domain.Profile) (*domain.Profile, error)
	ChangePassword(ctx context.Context, userID int, currentPassword, newPassword string) error
	DeleteAccount(ctx context.Context, userID int) error
}

type userService struct {
	userRepo    repository.UserRepository
	profileRepo repository.ProfileRepository
}

func NewUserService(userRepo repository.UserRepository, profileRepo repository.ProfileRepository) UserService {
	return &userService{
		userRepo:    userRepo,
		profileRepo: profileRepo,
	}
}

// GetProfile retrieves user and profile by ID
func (s *userService) GetProfile(ctx context.Context, userID int) (*domain.User, *domain.Profile, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, nil, fmt.Errorf("get user by id: %w", err)
	}

	profile, err := s.profileRepo.GetByUserID(ctx, userID)
	if err != nil {
		if err == domain.ErrNotFound {
			// Profile doesn't exist yet, return user with nil profile
			return user, nil, nil
		}
		return nil, nil, fmt.Errorf("get profile by user id: %w", err)
	}

	return user, profile, nil
}

// UpdateProfile updates user profile information (partial update supported)
func (s *userService) UpdateProfile(ctx context.Context, userID int, profileData *domain.Profile) (*domain.Profile, error) {
	// Get existing profile
	profile, err := s.profileRepo.GetByUserID(ctx, userID)
	if err != nil {
		if err == domain.ErrNotFound {
			// Create new profile
			profileData.UserID = userID
			if err := s.profileRepo.Create(ctx, profileData); err != nil {
				return nil, fmt.Errorf("create profile: %w", err)
			}
			return profileData, nil
		}
		return nil, fmt.Errorf("get profile: %w", err)
	}

	// Update only provided fields (partial update)
	if profileData.FirstName != "" {
		profile.FirstName = profileData.FirstName
	}
	if profileData.LastName != "" {
		profile.LastName = profileData.LastName
	}
	if profileData.MiddleName != nil {
		profile.MiddleName = profileData.MiddleName
	}
	if profileData.DateOfBirth != nil {
		profile.DateOfBirth = profileData.DateOfBirth
	}
	if profileData.Gender != nil {
		profile.Gender = profileData.Gender
	}
	if profileData.Phone != nil {
		profile.Phone = profileData.Phone
	}
	if profileData.Address != nil {
		profile.Address = profileData.Address
	}
	if profileData.City != nil {
		profile.City = profileData.City
	}
	if profileData.Country != nil {
		profile.Country = profileData.Country
	}
	if profileData.PostalCode != nil {
		profile.PostalCode = profileData.PostalCode
	}

	if err := s.profileRepo.Update(ctx, profile); err != nil {
		return nil, fmt.Errorf("update profile: %w", err)
	}

	return profile, nil
}

// ChangePassword changes user password
func (s *userService) ChangePassword(ctx context.Context, userID int, currentPassword, newPassword string) error {
	// Get user
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("get user by id: %w", err)
	}

	// Verify current password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(currentPassword)); err != nil {
		return fmt.Errorf("invalid current password")
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("hash password: %w", err)
	}

	// Update password
	user.Password = string(hashedPassword)
	user.UpdatedAt = time.Now()

	if err := s.userRepo.Update(ctx, user); err != nil {
		return fmt.Errorf("update user: %w", err)
	}

	return nil
}

// DeleteAccount marks user account as inactive (soft delete)
func (s *userService) DeleteAccount(ctx context.Context, userID int) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("get user by id: %w", err)
	}

	user.Status = "deleted"
	user.UpdatedAt = time.Now()

	if err := s.userRepo.Update(ctx, user); err != nil {
		return fmt.Errorf("update user: %w", err)
	}

	return nil
}
