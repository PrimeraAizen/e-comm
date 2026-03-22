package domain

import "time"

// UserProductView represents a user viewing a product
type UserProductView struct {
	UserID    string    `json:"user_id"`
	ProductID string    `json:"product_id"`
	ViewedAt  time.Time `json:"viewed_at"`
}

// UserProductLike represents a user liking a product
type UserProductLike struct {
	UserID    string    `json:"user_id"`
	ProductID string    `json:"product_id"`
	LikedAt   time.Time `json:"liked_at"`
}

// UserProductPurchase represents a user purchasing a product
type UserProductPurchase struct {
	UserID          string    `json:"user_id"`
	ProductID       string    `json:"product_id"`
	Quantity        int       `json:"quantity"`
	PriceAtPurchase float64   `json:"price_at_purchase"`
	PurchasedAt     time.Time `json:"purchased_at"`
}

// UserInteractionSummary provides an overview of user's interactions
type UserInteractionSummary struct {
	UserID            string               `json:"user_id"`
	ViewedProducts    []ProductInteraction `json:"viewed_products"`
	LikedProducts     []ProductInteraction `json:"liked_products"`
	PurchasedProducts []ProductInteraction `json:"purchased_products"`
	TotalViews        int64                `json:"total_views"`
	TotalLikes        int64                `json:"total_likes"`
	TotalPurchases    int64                `json:"total_purchases"`
}

// ProductInteraction represents a single product interaction with details
type ProductInteraction struct {
	ProductID    string    `json:"product_id"`
	ProductName  string    `json:"product_name"`
	CategoryID   string    `json:"category_id"`
	Price        float64   `json:"price"`
	InteractedAt time.Time `json:"interacted_at"`
}
