package service

import (
	"context"
	"fmt"

	"github.com/PrimeraAizen/e-comm/internal/domain"
	"github.com/PrimeraAizen/e-comm/internal/repository"
)

type ProductService interface {
	// Product operations
	CreateProduct(ctx context.Context, product *domain.Product) error
	GetProduct(ctx context.Context, id int) (*domain.Product, error)
	GetProductWithCategory(ctx context.Context, id int) (*domain.ProductWithCategory, error)
	UpdateProduct(ctx context.Context, product *domain.Product) error
	DeleteProduct(ctx context.Context, id int) error

	// Product listing and search
	ListProducts(ctx context.Context, filter domain.ProductFilter) ([]*domain.Product, int64, error)
	ListProductsWithCategories(ctx context.Context, filter domain.ProductFilter) ([]*domain.ProductWithCategory, int64, error)
	SearchProducts(ctx context.Context, query string, limit, offset int) ([]*domain.Product, int64, error)

	// Category operations
	CreateCategory(ctx context.Context, category *domain.Category) error
	GetCategory(ctx context.Context, id int) (*domain.Category, error)
	GetCategoryByName(ctx context.Context, name string) (*domain.Category, error)
	ListCategories(ctx context.Context) ([]*domain.Category, error)
	UpdateCategory(ctx context.Context, category *domain.Category) error
	DeleteCategory(ctx context.Context, id int) error

	// Product statistics
	GetProductStatistics(ctx context.Context, productID int) (*domain.ProductStatistics, error)
	RefreshStatistics(ctx context.Context) error

	// Stock management
	UpdateStock(ctx context.Context, productID int, quantity int) error
	CheckStock(ctx context.Context, productID int, quantity int) (bool, error)
}

type productService struct {
	productRepo repository.ProductRepository
}

func NewProductService(productRepo repository.ProductRepository) ProductService {
	return &productService{
		productRepo: productRepo,
	}
}

// CreateProduct creates a new product
func (s *productService) CreateProduct(ctx context.Context, product *domain.Product) error {
	// Validate product
	if err := s.validateProduct(product); err != nil {
		return err
	}

	// Check if category exists if provided
	if product.CategoryID != nil {
		_, err := s.productRepo.GetCategoryByID(ctx, *product.CategoryID)
		if err != nil {
			if err == domain.ErrNotFound {
				return fmt.Errorf("category not found")
			}
			return fmt.Errorf("check category: %w", err)
		}
	}

	// Set default values
	if product.Stock == 0 {
		product.Stock = 0
	}
	product.IsActive = true

	return s.productRepo.Create(ctx, product)
}

// GetProduct retrieves a product by ID
func (s *productService) GetProduct(ctx context.Context, id int) (*domain.Product, error) {
	return s.productRepo.GetByID(ctx, id)
}

// GetProductWithCategory retrieves a product with category information
func (s *productService) GetProductWithCategory(ctx context.Context, id int) (*domain.ProductWithCategory, error) {
	return s.productRepo.GetByIDWithCategory(ctx, id)
}

// UpdateProduct updates a product
func (s *productService) UpdateProduct(ctx context.Context, product *domain.Product) error {
	// Validate product
	if err := s.validateProduct(product); err != nil {
		return err
	}

	// Check if product exists
	existingProduct, err := s.productRepo.GetByID(ctx, product.ID)
	if err != nil {
		return err
	}

	// Check if category exists if changed
	if product.CategoryID != nil && (existingProduct.CategoryID == nil || *product.CategoryID != *existingProduct.CategoryID) {
		_, err := s.productRepo.GetCategoryByID(ctx, *product.CategoryID)
		if err != nil {
			if err == domain.ErrNotFound {
				return fmt.Errorf("category not found")
			}
			return fmt.Errorf("check category: %w", err)
		}
	}

	return s.productRepo.Update(ctx, product)
}

// DeleteProduct deletes a product
func (s *productService) DeleteProduct(ctx context.Context, id int) error {
	// Check if product exists
	_, err := s.productRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	return s.productRepo.Delete(ctx, id)
}

// ListProducts retrieves a list of products with filtering
func (s *productService) ListProducts(ctx context.Context, filter domain.ProductFilter) ([]*domain.Product, int64, error) {
	// Set default values
	if filter.Limit <= 0 {
		filter.Limit = 20
	}
	if filter.Limit > 100 {
		filter.Limit = 100 // Max limit
	}

	// Default to showing only active products for public listing
	if filter.IsActive == nil {
		active := true
		filter.IsActive = &active
	}

	return s.productRepo.List(ctx, filter)
}

// ListProductsWithCategories retrieves products with category names
func (s *productService) ListProductsWithCategories(ctx context.Context, filter domain.ProductFilter) ([]*domain.ProductWithCategory, int64, error) {
	// Set default values
	if filter.Limit <= 0 {
		filter.Limit = 20
	}
	if filter.Limit > 100 {
		filter.Limit = 100 // Max limit
	}

	// Default to showing only active products for public listing
	if filter.IsActive == nil {
		active := true
		filter.IsActive = &active
	}

	return s.productRepo.ListWithCategories(ctx, filter)
}

// SearchProducts performs full-text search on products
func (s *productService) SearchProducts(ctx context.Context, query string, limit, offset int) ([]*domain.Product, int64, error) {
	if query == "" {
		return nil, 0, fmt.Errorf("search query cannot be empty")
	}

	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	return s.productRepo.Search(ctx, query, limit, offset)
}

// CreateCategory creates a new category
func (s *productService) CreateCategory(ctx context.Context, category *domain.Category) error {
	// Validate category
	if category.Name == "" {
		return fmt.Errorf("category name is required")
	}

	// Check if category with same name already exists
	existingCategory, err := s.productRepo.GetCategoryByName(ctx, category.Name)
	if err != nil && err != domain.ErrNotFound {
		return fmt.Errorf("check existing category: %w", err)
	}
	if existingCategory != nil {
		return domain.ErrAlreadyExists
	}

	// Check if parent category exists if provided
	if category.ParentID != nil {
		_, err := s.productRepo.GetCategoryByID(ctx, *category.ParentID)
		if err != nil {
			if err == domain.ErrNotFound {
				return fmt.Errorf("parent category not found")
			}
			return fmt.Errorf("check parent category: %w", err)
		}
	}

	return s.productRepo.CreateCategory(ctx, category)
}

// GetCategory retrieves a category by ID
func (s *productService) GetCategory(ctx context.Context, id int) (*domain.Category, error) {
	return s.productRepo.GetCategoryByID(ctx, id)
}

// GetCategoryByName retrieves a category by name
func (s *productService) GetCategoryByName(ctx context.Context, name string) (*domain.Category, error) {
	return s.productRepo.GetCategoryByName(ctx, name)
}

// ListCategories retrieves all categories
func (s *productService) ListCategories(ctx context.Context) ([]*domain.Category, error) {
	return s.productRepo.ListCategories(ctx)
}

// UpdateCategory updates a category
func (s *productService) UpdateCategory(ctx context.Context, category *domain.Category) error {
	// Validate category
	if category.Name == "" {
		return fmt.Errorf("category name is required")
	}

	// Check if category exists
	_, err := s.productRepo.GetCategoryByID(ctx, category.ID)
	if err != nil {
		return err
	}

	// Check if parent category exists if changed
	if category.ParentID != nil {
		// Prevent self-reference
		if *category.ParentID == category.ID {
			return fmt.Errorf("category cannot be its own parent")
		}

		_, err := s.productRepo.GetCategoryByID(ctx, *category.ParentID)
		if err != nil {
			if err == domain.ErrNotFound {
				return fmt.Errorf("parent category not found")
			}
			return fmt.Errorf("check parent category: %w", err)
		}
	}

	return s.productRepo.UpdateCategory(ctx, category)
}

// DeleteCategory deletes a category
func (s *productService) DeleteCategory(ctx context.Context, id int) error {
	// Check if category exists
	_, err := s.productRepo.GetCategoryByID(ctx, id)
	if err != nil {
		return err
	}

	// TODO: Check if category has products or subcategories
	// For now, database CASCADE will handle it

	return s.productRepo.DeleteCategory(ctx, id)
}

// GetProductStatistics retrieves statistics for a product
func (s *productService) GetProductStatistics(ctx context.Context, productID int) (*domain.ProductStatistics, error) {
	// Check if product exists
	_, err := s.productRepo.GetByID(ctx, productID)
	if err != nil {
		return nil, err
	}

	return s.productRepo.GetProductStatistics(ctx, productID)
}

// RefreshStatistics refreshes the product statistics materialized view
func (s *productService) RefreshStatistics(ctx context.Context) error {
	return s.productRepo.RefreshProductStatistics(ctx)
}

// UpdateStock updates product stock
func (s *productService) UpdateStock(ctx context.Context, productID int, quantity int) error {
	product, err := s.productRepo.GetByID(ctx, productID)
	if err != nil {
		return err
	}

	newStock := product.Stock + quantity
	if newStock < 0 {
		return fmt.Errorf("insufficient stock")
	}

	product.Stock = newStock
	return s.productRepo.Update(ctx, product)
}

// CheckStock checks if sufficient stock is available
func (s *productService) CheckStock(ctx context.Context, productID int, quantity int) (bool, error) {
	product, err := s.productRepo.GetByID(ctx, productID)
	if err != nil {
		return false, err
	}

	return product.Stock >= quantity, nil
}

// validateProduct validates product data
func (s *productService) validateProduct(product *domain.Product) error {
	if product.Name == "" {
		return fmt.Errorf("product name is required")
	}

	if product.Price < 0 {
		return fmt.Errorf("product price cannot be negative")
	}

	if product.Stock < 0 {
		return fmt.Errorf("product stock cannot be negative")
	}

	return nil
}
