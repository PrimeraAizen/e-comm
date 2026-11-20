package domain

// ProductRecommendation represents a recommended product with a score
type ProductRecommendation struct {
	ProductID   string  `json:"product_id" bson:"product_id"`
	ProductName string  `json:"product_name" bson:"product_name"`
	CategoryID  string  `json:"category_id" bson:"category_id"`
	Price       float64 `json:"price" bson:"price"`
	Score       float64 `json:"score" bson:"score"`   // Similarity/relevance score
	Reason      string  `json:"reason" bson:"reason"` // Why recommended
}

// RecommendationResponse is the API response structure
type RecommendationResponse struct {
	UserID          string                  `json:"user_id"`
	Recommendations []ProductRecommendation `json:"recommendations"`
	Algorithm       string                  `json:"algorithm"` // e.g., "collaborative_filtering"
	GeneratedAt     string                  `json:"generated_at"`
}

// UserSimilarity represents similarity between two users
type UserSimilarity struct {
	UserID          string  `json:"user_id"`
	SimilarityScore float64 `json:"similarity_score"`
	CommonLikes     int     `json:"common_likes"`
	CommonViews     int     `json:"common_views"`
}
