package service

import (
	"context"
	"fmt"

	"github.com/PrimeraAizen/e-comm/internal/domain"
	"github.com/PrimeraAizen/e-comm/internal/repository"
)

type InteractionService interface {
	// View interactions
	RecordProductView(ctx context.Context, userID, productID int) error
	GetUserViewHistory(ctx context.Context, userID int, limit int) ([]domain.ProductInteraction, error)

	// Like interactions
	LikeProduct(ctx context.Context, userID, productID int) error
	UnlikeProduct(ctx context.Context, userID, productID int) error
	GetUserLikedProducts(ctx context.Context, userID int, limit int) ([]domain.ProductInteraction, error)
	IsProductLiked(ctx context.Context, userID, productID int) (bool, error)

	// Purchase interactions
	PurchaseProduct(ctx context.Context, userID, productID int, quantity int) error
	GetUserPurchaseHistory(ctx context.Context, userID int, limit int) ([]domain.ProductInteraction, error)
	HasPurchasedProduct(ctx context.Context, userID, productID int) (bool, error)

	// Summary
	GetUserInteractionSummary(ctx context.Context, userID int) (*domain.UserInteractionSummary, error)
}

type interactionService struct {
	interactionRepo repository.InteractionRepository
	productRepo     repository.ProductRepository
}

func NewInteractionService(
	interactionRepo repository.InteractionRepository,
	productRepo repository.ProductRepository,
) InteractionService {
	return &interactionService{
		interactionRepo: interactionRepo,
		productRepo:     productRepo,
	}
}

// RecordProductView records a user viewing a product
func (s *interactionService) RecordProductView(ctx context.Context, userID, productID int) error {
	// Verify product exists
	_, err := s.productRepo.GetByID(ctx, productID)
	if err != nil {
		if err == domain.ErrNotFound {
			return fmt.Errorf("product not found")
		}
		return fmt.Errorf("verify product: %w", err)
	}

	// Record the view
	if err := s.interactionRepo.RecordView(ctx, userID, productID); err != nil {
		return fmt.Errorf("record view: %w", err)
	}

	return nil
}

// GetUserViewHistory retrieves the user's view history
func (s *interactionService) GetUserViewHistory(ctx context.Context, userID int, limit int) ([]domain.ProductInteraction, error) {
	if limit <= 0 || limit > 100 {
		limit = 50 // Default limit
	}

	views, err := s.interactionRepo.GetUserViews(ctx, userID, limit)
	if err != nil {
		return nil, fmt.Errorf("get view history: %w", err)
	}

	return views, nil
}

// LikeProduct records a user liking a product
func (s *interactionService) LikeProduct(ctx context.Context, userID, productID int) error {
	// Verify product exists
	_, err := s.productRepo.GetByID(ctx, productID)
	if err != nil {
		if err == domain.ErrNotFound {
			return fmt.Errorf("product not found")
		}
		return fmt.Errorf("verify product: %w", err)
	}

	// Record the like
	if err := s.interactionRepo.RecordLike(ctx, userID, productID); err != nil {
		return fmt.Errorf("record like: %w", err)
	}

	return nil
}

// UnlikeProduct removes a user's like from a product
func (s *interactionService) UnlikeProduct(ctx context.Context, userID, productID int) error {
	if err := s.interactionRepo.RemoveLike(ctx, userID, productID); err != nil {
		if err == domain.ErrNotFound {
			return fmt.Errorf("like not found")
		}
		return fmt.Errorf("remove like: %w", err)
	}

	return nil
}

// GetUserLikedProducts retrieves products the user has liked
func (s *interactionService) GetUserLikedProducts(ctx context.Context, userID int, limit int) ([]domain.ProductInteraction, error) {
	if limit <= 0 || limit > 100 {
		limit = 50 // Default limit
	}

	likes, err := s.interactionRepo.GetUserLikes(ctx, userID, limit)
	if err != nil {
		return nil, fmt.Errorf("get liked products: %w", err)
	}

	return likes, nil
}

// IsProductLiked checks if the user has liked a product
func (s *interactionService) IsProductLiked(ctx context.Context, userID, productID int) (bool, error) {
	liked, err := s.interactionRepo.HasLiked(ctx, userID, productID)
	if err != nil {
		return false, fmt.Errorf("check if liked: %w", err)
	}

	return liked, nil
}

// GetUserInteractionSummary gets a summary of all user interactions
func (s *interactionService) GetUserInteractionSummary(ctx context.Context, userID int) (*domain.UserInteractionSummary, error) {
	summary, err := s.interactionRepo.GetUserInteractionSummary(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get interaction summary: %w", err)
	}

	return summary, nil
}

// PurchaseProduct records a user purchasing a product
func (s *interactionService) PurchaseProduct(ctx context.Context, userID, productID int, quantity int) error {
	if quantity <= 0 {
		return fmt.Errorf("quantity must be greater than 0")
	}

	// Verify product exists and get current price
	product, err := s.productRepo.GetByID(ctx, productID)
	if err != nil {
		if err == domain.ErrNotFound {
			return fmt.Errorf("product not found")
		}
		return fmt.Errorf("verify product: %w", err)
	}

	// Check stock availability
	if product.Stock < quantity {
		return fmt.Errorf("insufficient stock: requested %d, available %d", quantity, product.Stock)
	}

	// Record the purchase
	if err := s.interactionRepo.RecordPurchase(ctx, userID, productID, quantity, product.Price); err != nil {
		return fmt.Errorf("record purchase: %w", err)
	}

	// Update stock (reduce by purchased quantity)
	product.Stock -= quantity
	if err := s.productRepo.Update(ctx, product); err != nil {
		return fmt.Errorf("update product stock: %w", err)
	}

	return nil
}

// GetUserPurchaseHistory retrieves the user's purchase history
func (s *interactionService) GetUserPurchaseHistory(ctx context.Context, userID int, limit int) ([]domain.ProductInteraction, error) {
	if limit <= 0 || limit > 100 {
		limit = 50 // Default limit
	}

	purchases, err := s.interactionRepo.GetUserPurchases(ctx, userID, limit)
	if err != nil {
		return nil, fmt.Errorf("get purchase history: %w", err)
	}

	return purchases, nil
}

// HasPurchasedProduct checks if the user has purchased a product
func (s *interactionService) HasPurchasedProduct(ctx context.Context, userID, productID int) (bool, error) {
	purchased, err := s.interactionRepo.HasPurchased(ctx, userID, productID)
	if err != nil {
		return false, fmt.Errorf("check if purchased: %w", err)
	}

	return purchased, nil
}
