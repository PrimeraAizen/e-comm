package domain

import (
	"time"
)

// Product represents a product in the catalog
type Product struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CategoryID  *string   `json:"category_id,omitempty"`
	Price       float64   `json:"price"`
	Stock       int       `json:"stock"`
	ImageURL    string    `json:"image_url,omitempty"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Category represents a product category
type Category struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	ParentID    *string   `json:"parent_id,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ProductWithCategory includes category details
type ProductWithCategory struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	CategoryID   *string   `json:"category_id,omitempty"`
	Price        float64   `json:"price"`
	Stock        int       `json:"stock"`
	ImageURL     string    `json:"image_url,omitempty"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	CategoryName string    `json:"category_name,omitempty"`
}

// ProductFilter represents filtering options for products
type ProductFilter struct {
	CategoryID  *string
	MinPrice    *float64
	MaxPrice    *float64
	IsActive    *bool
	SearchQuery string
	Limit       int
	Offset      int
	SortBy      string // name, price, created_at
	SortOrder   string // asc, desc
}

// ProductStatistics represents aggregated product metrics
type ProductStatistics struct {
	ProductID     string  `json:"product_id"`
	ProductName   string  `json:"product_name"`
	ViewCount     int64   `json:"view_count"`
	LikeCount     int64   `json:"like_count"`
	PurchaseCount int64   `json:"purchase_count"`
	AverageRating float64 `json:"average_rating"`
	ReviewCount   int64   `json:"review_count"`
}
