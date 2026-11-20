package dto

import (
	"github.com/PrimeraAizen/e-comm/internal/domain"
)

type CreateProductRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description"`
	CategoryID  *string `json:"category_id"`
	Price       float64 `json:"price" binding:"required,min=0"`
	Stock       int     `json:"stock" binding:"min=0"`
	ImageURL    string  `json:"image_url"`
}

type UpdateProductRequest struct {
	Name        *string  `json:"name"`
	Description *string  `json:"description"`
	CategoryID  *string  `json:"category_id"`
	Price       *float64 `json:"price"`
	Stock       *int     `json:"stock"`
	ImageURL    *string  `json:"image_url"`
	IsActive    *bool    `json:"is_active"`
}

type ProductListResponse struct {
	Products []*domain.ProductWithCategory `json:"products"`
	Total    int64                         `json:"total"`
	Page     int                           `json:"page"`
	Limit    int                           `json:"limit"`
}

type CreateCategoryRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description"`
	ParentID    *string `json:"parent_id"`
}

type UpdateCategoryRequest struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
	ParentID    *string `json:"parent_id"`
}

type PurchaseProductRequest struct {
	Quantity int `json:"quantity" binding:"required,min=1"`
}
