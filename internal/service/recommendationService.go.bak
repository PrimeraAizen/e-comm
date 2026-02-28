package service

import (
	"context"
	"fmt"
	"math"
	"sort"
	"time"

	"github.com/PrimeraAizen/e-comm/internal/domain"
	"github.com/PrimeraAizen/e-comm/internal/repository"
)

type RecommendationService interface {
	GetRecommendations(ctx context.Context, userID int, limit int) (*domain.RecommendationResponse, error)
	GetSimilarUsers(ctx context.Context, userID int, limit int) ([]domain.UserSimilarity, error)
}

type recommendationService struct {
	interactionRepo repository.InteractionRepository
	productRepo     repository.ProductRepository
}

func NewRecommendationService(
	interactionRepo repository.InteractionRepository,
	productRepo repository.ProductRepository,
) RecommendationService {
	return &recommendationService{
		interactionRepo: interactionRepo,
		productRepo:     productRepo,
	}
}

// GetRecommendations generates product recommendations using collaborative filtering
func (s *recommendationService) GetRecommendations(ctx context.Context, userID int, limit int) (*domain.RecommendationResponse, error) {
	if limit <= 0 || limit > 50 {
		limit = 10 // Default limit
	}

	// Get all interactions
	allLikes, err := s.interactionRepo.GetAllUserLikes(ctx)
	if err != nil {
		return nil, fmt.Errorf("get all likes: %w", err)
	}

	allViews, err := s.interactionRepo.GetAllUserViews(ctx)
	if err != nil {
		return nil, fmt.Errorf("get all views: %w", err)
	}

	allPurchases, err := s.interactionRepo.GetAllUserPurchases(ctx)
	if err != nil {
		return nil, fmt.Errorf("get all purchases: %w", err)
	}

	// Create sets for current user's interactions
	userLikedProducts := make(map[int]bool)
	userViewedProducts := make(map[int]bool)
	userPurchasedProducts := make(map[int]bool)

	for _, like := range allLikes {
		if like.UserID == userID {
			userLikedProducts[like.ProductID] = true
		}
	}
	for _, view := range allViews {
		if view.UserID == userID {
			userViewedProducts[view.ProductID] = true
		}
	}
	for _, purchase := range allPurchases {
		if purchase.UserID == userID {
			userPurchasedProducts[purchase.ProductID] = true
		}
	}

	// If user has no interactions, return popular products
	if len(userLikedProducts) == 0 && len(userViewedProducts) == 0 && len(userPurchasedProducts) == 0 {
		return s.getPopularProducts(ctx, limit)
	}

	// Find similar users based on collaborative filtering
	similarUsers, err := s.GetSimilarUsers(ctx, userID, 10)
	if err != nil {
		return nil, fmt.Errorf("get similar users: %w", err)
	}

	// If no similar users, return popular products
	if len(similarUsers) == 0 {
		return s.getPopularProducts(ctx, limit)
	}

	// Aggregate recommendations from similar users
	productScores := make(map[int]float64)
	productDetails := make(map[int]*domain.Product)

	// Score from similar users' purchases (strongest signal - weight 3.0)
	for _, simUser := range similarUsers {
		for _, purchase := range allPurchases {
			if purchase.UserID != simUser.UserID {
				continue
			}

			// Skip products the user already purchased
			if userPurchasedProducts[purchase.ProductID] {
				continue
			}

			// Get product details if not cached
			if productDetails[purchase.ProductID] == nil {
				product, err := s.productRepo.GetByID(ctx, purchase.ProductID)
				if err != nil {
					continue
				}
				productDetails[purchase.ProductID] = product
			}

			// Weight by user similarity score and boost for purchases
			productScores[purchase.ProductID] += simUser.SimilarityScore * 3.0
		}
	}

	// Score from similar users' likes (medium signal - weight 1.5)
	for _, simUser := range similarUsers {
		for _, like := range allLikes {
			if like.UserID != simUser.UserID {
				continue
			}

			// Skip products the user already liked or purchased
			if userLikedProducts[like.ProductID] || userPurchasedProducts[like.ProductID] {
				continue
			}

			// Get product details if not cached
			if productDetails[like.ProductID] == nil {
				product, err := s.productRepo.GetByID(ctx, like.ProductID)
				if err != nil {
					continue
				}
				productDetails[like.ProductID] = product
			}

			// Weight by user similarity score
			productScores[like.ProductID] += simUser.SimilarityScore * 1.5
		}
	}

	// Convert to recommendation list
	recommendations := make([]domain.ProductRecommendation, 0, limit)
	for productID, score := range productScores {
		product := productDetails[productID]
		if product == nil {
			continue
		}

		categoryID := 0
		if product.CategoryID != nil {
			categoryID = *product.CategoryID
		}

		recommendations = append(recommendations, domain.ProductRecommendation{
			ProductID:   productID,
			ProductName: product.Name,
			CategoryID:  categoryID,
			Price:       product.Price,
			Score:       score,
			Reason:      "Users with similar interests liked this",
		})
	}

	// Sort by score descending
	sort.Slice(recommendations, func(i, j int) bool {
		return recommendations[i].Score > recommendations[j].Score
	})

	// Limit results
	if len(recommendations) > limit {
		recommendations = recommendations[:limit]
	}

	// If still no recommendations, fallback to popular products
	if len(recommendations) == 0 {
		return s.getPopularProducts(ctx, limit)
	}

	return &domain.RecommendationResponse{
		UserID:          userID,
		Recommendations: recommendations,
		Algorithm:       "collaborative_filtering",
		GeneratedAt:     time.Now().Format(time.RFC3339),
	}, nil
}

// GetSimilarUsers finds users with similar interaction patterns
func (s *recommendationService) GetSimilarUsers(ctx context.Context, userID int, limit int) ([]domain.UserSimilarity, error) {
	// Get all likes, views, and purchases
	allLikes, err := s.interactionRepo.GetAllUserLikes(ctx)
	if err != nil {
		return nil, fmt.Errorf("get all likes: %w", err)
	}

	allViews, err := s.interactionRepo.GetAllUserViews(ctx)
	if err != nil {
		return nil, fmt.Errorf("get all views: %w", err)
	}

	allPurchases, err := s.interactionRepo.GetAllUserPurchases(ctx)
	if err != nil {
		return nil, fmt.Errorf("get all purchases: %w", err)
	}

	// Create sets for current user and group by user for others
	userLikedProducts := make(map[int]bool)
	userViewedProducts := make(map[int]bool)
	userPurchasedProducts := make(map[int]bool)
	otherUsersLikes := make(map[int]map[int]bool)
	otherUsersViews := make(map[int]map[int]bool)
	otherUsersPurchases := make(map[int]map[int]bool)

	for _, like := range allLikes {
		if like.UserID == userID {
			userLikedProducts[like.ProductID] = true
		} else {
			if otherUsersLikes[like.UserID] == nil {
				otherUsersLikes[like.UserID] = make(map[int]bool)
			}
			otherUsersLikes[like.UserID][like.ProductID] = true
		}
	}

	for _, view := range allViews {
		if view.UserID == userID {
			userViewedProducts[view.ProductID] = true
		} else {
			if otherUsersViews[view.UserID] == nil {
				otherUsersViews[view.UserID] = make(map[int]bool)
			}
			otherUsersViews[view.UserID][view.ProductID] = true
		}
	}

	for _, purchase := range allPurchases {
		if purchase.UserID == userID {
			userPurchasedProducts[purchase.ProductID] = true
		} else {
			if otherUsersPurchases[purchase.UserID] == nil {
				otherUsersPurchases[purchase.UserID] = make(map[int]bool)
			}
			otherUsersPurchases[purchase.UserID][purchase.ProductID] = true
		}
	}

	// Collect all unique user IDs
	allUserIDs := make(map[int]bool)
	for userID := range otherUsersLikes {
		allUserIDs[userID] = true
	}
	for userID := range otherUsersViews {
		allUserIDs[userID] = true
	}
	for userID := range otherUsersPurchases {
		allUserIDs[userID] = true
	}

	// Calculate similarity with each user
	similarities := make([]domain.UserSimilarity, 0)

	for otherUserID := range allUserIDs {
		otherLikes := otherUsersLikes[otherUserID]
		otherViews := otherUsersViews[otherUserID]
		otherPurchases := otherUsersPurchases[otherUserID]

		// Calculate Jaccard similarity for purchases (strongest signal)
		commonPurchases := 0
		for productID := range userPurchasedProducts {
			if otherPurchases != nil && otherPurchases[productID] {
				commonPurchases++
			}
		}

		// Calculate Jaccard similarity for likes
		commonLikes := 0
		for productID := range userLikedProducts {
			if otherLikes != nil && otherLikes[productID] {
				commonLikes++
			}
		}

		// Calculate Jaccard similarity for views
		commonViews := 0
		for productID := range userViewedProducts {
			if otherViews != nil && otherViews[productID] {
				commonViews++
			}
		}

		// Need at least one common interaction
		if commonLikes == 0 && commonViews == 0 && commonPurchases == 0 {
			continue
		}

		// Jaccard similarity: |A ∩ B| / |A ∪ B|
		unionPurchases := len(userPurchasedProducts) + len(otherPurchases) - commonPurchases
		unionLikes := len(userLikedProducts) + len(otherLikes) - commonLikes
		unionViews := len(userViewedProducts) + len(otherViews) - commonViews

		purchaseSimilarity := 0.0
		if unionPurchases > 0 {
			purchaseSimilarity = float64(commonPurchases) / float64(unionPurchases)
		}

		likeSimilarity := 0.0
		if unionLikes > 0 {
			likeSimilarity = float64(commonLikes) / float64(unionLikes)
		}

		viewSimilarity := 0.0
		if unionViews > 0 {
			viewSimilarity = float64(commonViews) / float64(unionViews)
		}

		// Combined similarity (purchases weighted most heavily)
		// Purchases: 50%, Likes: 35%, Views: 15%
		similarity := (purchaseSimilarity * 0.5) + (likeSimilarity * 0.35) + (viewSimilarity * 0.15)

		// Apply minimum threshold
		if similarity < 0.1 {
			continue
		}

		similarities = append(similarities, domain.UserSimilarity{
			UserID:          otherUserID,
			SimilarityScore: similarity,
			CommonLikes:     commonLikes,
			CommonViews:     commonViews,
		})
	}

	// Sort by similarity descending
	sort.Slice(similarities, func(i, j int) bool {
		return similarities[i].SimilarityScore > similarities[j].SimilarityScore
	})

	// Limit results
	if len(similarities) > limit {
		similarities = similarities[:limit]
	}

	return similarities, nil
}

// getPopularProducts returns most liked products as fallback
func (s *recommendationService) getPopularProducts(ctx context.Context, limit int) (*domain.RecommendationResponse, error) {
	// Get all likes
	allLikes, err := s.interactionRepo.GetAllUserLikes(ctx)
	if err != nil {
		return nil, fmt.Errorf("get all likes: %w", err)
	}

	// Count likes per product
	likeCounts := make(map[int]int)
	for _, like := range allLikes {
		likeCounts[like.ProductID]++
	}

	// Create sorted list
	type productCount struct {
		productID int
		count     int
	}

	productCounts := make([]productCount, 0, len(likeCounts))
	for productID, count := range likeCounts {
		productCounts = append(productCounts, productCount{productID, count})
	}

	sort.Slice(productCounts, func(i, j int) bool {
		return productCounts[i].count > productCounts[j].count
	})

	// Limit and get product details
	if len(productCounts) > limit {
		productCounts = productCounts[:limit]
	}

	recommendations := make([]domain.ProductRecommendation, 0, len(productCounts))
	maxCount := 1
	if len(productCounts) > 0 {
		maxCount = productCounts[0].count
	}

	for _, pc := range productCounts {
		product, err := s.productRepo.GetByID(ctx, pc.productID)
		if err != nil {
			continue
		}

		// Normalize score to 0-1 range
		score := float64(pc.count) / float64(maxCount)

		categoryID := 0
		if product.CategoryID != nil {
			categoryID = *product.CategoryID
		}

		recommendations = append(recommendations, domain.ProductRecommendation{
			ProductID:   pc.productID,
			ProductName: product.Name,
			CategoryID:  categoryID,
			Price:       product.Price,
			Score:       score,
			Reason:      fmt.Sprintf("Popular choice - %d users liked this", pc.count),
		})
	}

	return &domain.RecommendationResponse{
		UserID:          0,
		Recommendations: recommendations,
		Algorithm:       "popularity_based",
		GeneratedAt:     time.Now().Format(time.RFC3339),
	}, nil
}

// Helper function to calculate cosine similarity (alternative to Jaccard)
func cosineSimilarity(a, b map[int]bool) float64 {
	if len(a) == 0 || len(b) == 0 {
		return 0.0
	}

	dotProduct := 0
	for key := range a {
		if b[key] {
			dotProduct++
		}
	}

	magnitudeA := math.Sqrt(float64(len(a)))
	magnitudeB := math.Sqrt(float64(len(b)))

	if magnitudeA == 0 || magnitudeB == 0 {
		return 0.0
	}

	return float64(dotProduct) / (magnitudeA * magnitudeB)
}
