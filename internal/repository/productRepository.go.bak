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

type ProductRepository interface {
	// Product CRUD
	Create(ctx context.Context, product *domain.Product) error
	GetByID(ctx context.Context, id int) (*domain.Product, error)
	GetByIDWithCategory(ctx context.Context, id int) (*domain.ProductWithCategory, error)
	Update(ctx context.Context, product *domain.Product) error
	Delete(ctx context.Context, id int) error

	// Product listing and search
	List(ctx context.Context, filter domain.ProductFilter) ([]*domain.Product, int64, error)
	ListWithCategories(ctx context.Context, filter domain.ProductFilter) ([]*domain.ProductWithCategory, int64, error)
	Search(ctx context.Context, query string, limit, offset int) ([]*domain.Product, int64, error)

	// Category CRUD
	CreateCategory(ctx context.Context, category *domain.Category) error
	GetCategoryByID(ctx context.Context, id int) (*domain.Category, error)
	GetCategoryByName(ctx context.Context, name string) (*domain.Category, error)
	ListCategories(ctx context.Context) ([]*domain.Category, error)
	UpdateCategory(ctx context.Context, category *domain.Category) error
	DeleteCategory(ctx context.Context, id int) error

	// Product statistics
	GetProductStatistics(ctx context.Context, productID int) (*domain.ProductStatistics, error)
	RefreshProductStatistics(ctx context.Context) error
}

type productRepository struct {
	db *mongodb.MongoDB
}

func NewProductRepository(db *mongodb.MongoDB) ProductRepository {
	return &productRepository{db: db}
}

// Create creates a new product
func (r *productRepository) Create(ctx context.Context, product *domain.Product) error {
	// Get next ID
	nextID, err := r.getNextProductID(ctx)
	if err != nil {
		return fmt.Errorf("get next ID: %w", err)
	}
	product.ID = nextID
	product.CreatedAt = time.Now()
	product.UpdatedAt = time.Now()
	product.IsActive = true

	collection := r.db.Collection("products")
	_, err = collection.InsertOne(ctx, product)
	if err != nil {
		return fmt.Errorf("create product: %w", err)
	}

	return nil
}

// GetByID retrieves a product by ID
func (r *productRepository) GetByID(ctx context.Context, id int) (*domain.Product, error) {
	collection := r.db.Collection("products")

	var product domain.Product
	err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&product)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("get product by id: %w", err)
	}

	return &product, nil
}

// GetByIDWithCategory retrieves a product with category information
func (r *productRepository) GetByIDWithCategory(ctx context.Context, id int) (*domain.ProductWithCategory, error) {
	collection := r.db.Collection("products")

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{"_id": id}}},
		{{Key: "$lookup", Value: bson.M{
			"from":         "categories",
			"localField":   "category_id",
			"foreignField": "_id",
			"as":           "category",
		}}},
		{{Key: "$unwind", Value: bson.M{
			"path":                       "$category",
			"preserveNullAndEmptyArrays": true,
		}}},
		{{Key: "$addFields", Value: bson.M{
			"category_name": "$category.name",
		}}},
		{{Key: "$project", Value: bson.M{
			"category": 0,
		}}},
	}

	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, fmt.Errorf("aggregate product with category: %w", err)
	}
	defer cursor.Close(ctx)

	var result domain.ProductWithCategory
	if cursor.Next(ctx) {
		if err := cursor.Decode(&result); err != nil {
			return nil, fmt.Errorf("decode product with category: %w", err)
		}
		return &result, nil
	}

	return nil, domain.ErrNotFound
}

// Update updates a product
func (r *productRepository) Update(ctx context.Context, product *domain.Product) error {
	collection := r.db.Collection("products")

	product.UpdatedAt = time.Now()

	update := bson.M{
		"$set": bson.M{
			"name":        product.Name,
			"description": product.Description,
			"category_id": product.CategoryID,
			"price":       product.Price,
			"stock":       product.Stock,
			"image_url":   product.ImageURL,
			"is_active":   product.IsActive,
			"updated_at":  product.UpdatedAt,
		},
	}

	result, err := collection.UpdateOne(ctx, bson.M{"_id": product.ID}, update)
	if err != nil {
		return fmt.Errorf("update product: %w", err)
	}

	if result.MatchedCount == 0 {
		return domain.ErrNotFound
	}

	return nil
}

// Delete deletes a product
func (r *productRepository) Delete(ctx context.Context, id int) error {
	collection := r.db.Collection("products")

	result, err := collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return fmt.Errorf("delete product: %w", err)
	}

	if result.DeletedCount == 0 {
		return domain.ErrNotFound
	}

	return nil
}

// List retrieves products with filtering and pagination
func (r *productRepository) List(ctx context.Context, filter domain.ProductFilter) ([]*domain.Product, int64, error) {
	collection := r.db.Collection("products")

	// Build filter
	mongoFilter := bson.M{}

	if filter.CategoryID != nil {
		mongoFilter["category_id"] = *filter.CategoryID
	}

	if filter.MinPrice != nil {
		if _, ok := mongoFilter["price"]; !ok {
			mongoFilter["price"] = bson.M{}
		}
		mongoFilter["price"].(bson.M)["$gte"] = *filter.MinPrice
	}

	if filter.MaxPrice != nil {
		if _, ok := mongoFilter["price"]; !ok {
			mongoFilter["price"] = bson.M{}
		}
		mongoFilter["price"].(bson.M)["$lte"] = *filter.MaxPrice
	}

	if filter.IsActive != nil {
		mongoFilter["is_active"] = *filter.IsActive
	}

	if filter.SearchQuery != "" {
		mongoFilter["$text"] = bson.M{"$search": filter.SearchQuery}
	}

	// Count total
	total, err := collection.CountDocuments(ctx, mongoFilter)
	if err != nil {
		return nil, 0, fmt.Errorf("count products: %w", err)
	}

	// Build options
	opts := options.Find()

	// Sort
	sortField := "created_at"
	if filter.SortBy != "" {
		sortField = filter.SortBy
	}
	sortOrder := -1 // desc by default
	if filter.SortOrder == "asc" {
		sortOrder = 1
	}
	opts.SetSort(bson.M{sortField: sortOrder})

	// Pagination
	if filter.Limit > 0 {
		opts.SetLimit(int64(filter.Limit))
	}
	if filter.Offset > 0 {
		opts.SetSkip(int64(filter.Offset))
	}

	// Execute query
	cursor, err := collection.Find(ctx, mongoFilter, opts)
	if err != nil {
		return nil, 0, fmt.Errorf("find products: %w", err)
	}
	defer cursor.Close(ctx)

	var products []*domain.Product
	if err := cursor.All(ctx, &products); err != nil {
		return nil, 0, fmt.Errorf("decode products: %w", err)
	}

	return products, total, nil
}

// ListWithCategories retrieves products with category names
func (r *productRepository) ListWithCategories(ctx context.Context, filter domain.ProductFilter) ([]*domain.ProductWithCategory, int64, error) {
	collection := r.db.Collection("products")

	// Build match stage
	matchStage := bson.M{}

	if filter.CategoryID != nil {
		matchStage["category_id"] = *filter.CategoryID
	}

	if filter.MinPrice != nil {
		if _, ok := matchStage["price"]; !ok {
			matchStage["price"] = bson.M{}
		}
		matchStage["price"].(bson.M)["$gte"] = *filter.MinPrice
	}

	if filter.MaxPrice != nil {
		if _, ok := matchStage["price"]; !ok {
			matchStage["price"] = bson.M{}
		}
		matchStage["price"].(bson.M)["$lte"] = *filter.MaxPrice
	}

	if filter.IsActive != nil {
		matchStage["is_active"] = *filter.IsActive
	}

	if filter.SearchQuery != "" {
		matchStage["$text"] = bson.M{"$search": filter.SearchQuery}
	}

	// Build pipeline
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: matchStage}},
		{{Key: "$lookup", Value: bson.M{
			"from":         "categories",
			"localField":   "category_id",
			"foreignField": "_id",
			"as":           "category",
		}}},
		{{Key: "$unwind", Value: bson.M{
			"path":                       "$category",
			"preserveNullAndEmptyArrays": true,
		}}},
		{{Key: "$addFields", Value: bson.M{
			"category_name": "$category.name",
		}}},
		{{Key: "$project", Value: bson.M{
			"category": 0,
		}}},
	}

	// Count total
	countPipeline := append(pipeline, bson.D{{Key: "$count", Value: "total"}})
	countCursor, err := collection.Aggregate(ctx, countPipeline)
	if err != nil {
		return nil, 0, fmt.Errorf("count products: %w", err)
	}
	defer countCursor.Close(ctx)

	var countResult []struct {
		Total int64 `bson:"total"`
	}
	if err := countCursor.All(ctx, &countResult); err != nil {
		return nil, 0, fmt.Errorf("decode count: %w", err)
	}

	total := int64(0)
	if len(countResult) > 0 {
		total = countResult[0].Total
	}

	// Sort
	sortField := "created_at"
	if filter.SortBy != "" {
		sortField = filter.SortBy
	}
	sortOrder := -1
	if filter.SortOrder == "asc" {
		sortOrder = 1
	}
	pipeline = append(pipeline, bson.D{{Key: "$sort", Value: bson.M{sortField: sortOrder}}})

	// Pagination
	if filter.Offset > 0 {
		pipeline = append(pipeline, bson.D{{Key: "$skip", Value: filter.Offset}})
	}
	if filter.Limit > 0 {
		pipeline = append(pipeline, bson.D{{Key: "$limit", Value: filter.Limit}})
	}

	// Execute query
	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, 0, fmt.Errorf("aggregate products: %w", err)
	}
	defer cursor.Close(ctx)

	var products []*domain.ProductWithCategory
	for cursor.Next(ctx) {
		var rawDoc bson.M
		if err := cursor.Decode(&rawDoc); err != nil {
			return nil, 0, fmt.Errorf("decode raw doc: %w", err)
		}

		// Convert to bytes and back to properly handle UUID conversion
		rawBytes, err := bson.Marshal(rawDoc)
		if err != nil {
			return nil, 0, fmt.Errorf("marshal doc: %w", err)
		}

		var product domain.ProductWithCategory
		if err := bson.Unmarshal(rawBytes, &product); err != nil {
			return nil, 0, fmt.Errorf("unmarshal product: %w", err)
		}

		products = append(products, &product)
	}

	if err := cursor.Err(); err != nil {
		return nil, 0, fmt.Errorf("cursor error: %w", err)
	}

	return products, total, nil
}

// Search searches for products (alias for List with search query)
func (r *productRepository) Search(ctx context.Context, query string, limit, offset int) ([]*domain.Product, int64, error) {
	filter := domain.ProductFilter{
		SearchQuery: query,
		Limit:       limit,
		Offset:      offset,
	}
	return r.List(ctx, filter)
}

// CreateCategory creates a new category
func (r *productRepository) CreateCategory(ctx context.Context, category *domain.Category) error {
	// Get next ID
	nextID, err := r.getNextCategoryID(ctx)
	if err != nil {
		return fmt.Errorf("get next ID: %w", err)
	}
	category.ID = nextID
	category.CreatedAt = time.Now()
	category.UpdatedAt = time.Now()

	collection := r.db.Collection("categories")
	_, err = collection.InsertOne(ctx, category)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return fmt.Errorf("category with this name already exists: %w", err)
		}
		return fmt.Errorf("create category: %w", err)
	}

	return nil
}

// GetCategoryByID retrieves a category by ID
func (r *productRepository) GetCategoryByID(ctx context.Context, id int) (*domain.Category, error) {
	collection := r.db.Collection("categories")

	var category domain.Category
	err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&category)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("get category by id: %w", err)
	}

	return &category, nil
}

// GetCategoryByName retrieves a category by name
func (r *productRepository) GetCategoryByName(ctx context.Context, name string) (*domain.Category, error) {
	collection := r.db.Collection("categories")

	var category domain.Category
	err := collection.FindOne(ctx, bson.M{"name": name}).Decode(&category)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("get category by name: %w", err)
	}

	return &category, nil
}

// ListCategories retrieves all categories
func (r *productRepository) ListCategories(ctx context.Context) ([]*domain.Category, error) {
	collection := r.db.Collection("categories")

	cursor, err := collection.Find(ctx, bson.M{}, options.Find().SetSort(bson.M{"name": 1}))
	if err != nil {
		return nil, fmt.Errorf("find categories: %w", err)
	}
	defer cursor.Close(ctx)

	var categories []*domain.Category
	if err := cursor.All(ctx, &categories); err != nil {
		return nil, fmt.Errorf("decode categories: %w", err)
	}

	return categories, nil
}

// UpdateCategory updates a category
func (r *productRepository) UpdateCategory(ctx context.Context, category *domain.Category) error {
	collection := r.db.Collection("categories")

	category.UpdatedAt = time.Now()

	update := bson.M{
		"$set": bson.M{
			"name":        category.Name,
			"description": category.Description,
			"parent_id":   category.ParentID,
			"updated_at":  category.UpdatedAt,
		},
	}

	result, err := collection.UpdateOne(ctx, bson.M{"_id": category.ID}, update)
	if err != nil {
		return fmt.Errorf("update category: %w", err)
	}

	if result.MatchedCount == 0 {
		return domain.ErrNotFound
	}

	return nil
}

// DeleteCategory deletes a category
func (r *productRepository) DeleteCategory(ctx context.Context, id int) error {
	collection := r.db.Collection("categories")

	result, err := collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return fmt.Errorf("delete category: %w", err)
	}

	if result.DeletedCount == 0 {
		return domain.ErrNotFound
	}

	return nil
}

// GetProductStatistics retrieves statistics for a product
func (r *productRepository) GetProductStatistics(ctx context.Context, productID int) (*domain.ProductStatistics, error) {
	product, err := r.GetByID(ctx, productID)
	if err != nil {
		return nil, err
	}

	// Count views
	viewsCollection := r.db.Collection("user_product_views")
	viewCount, err := viewsCollection.CountDocuments(ctx, bson.M{"product_id": productID})
	if err != nil {
		viewCount = 0
	}

	// Count likes
	likesCollection := r.db.Collection("user_product_likes")
	likeCount, err := likesCollection.CountDocuments(ctx, bson.M{"product_id": productID})
	if err != nil {
		likeCount = 0
	}

	// Count purchases from order_items
	ordersCollection := r.db.Collection("order_items")
	purchaseCount, err := ordersCollection.CountDocuments(ctx, bson.M{"product_id": productID})
	if err != nil {
		purchaseCount = 0
	}

	stats := &domain.ProductStatistics{
		ProductID:     productID,
		ProductName:   product.Name,
		ViewCount:     viewCount,
		LikeCount:     likeCount,
		PurchaseCount: purchaseCount,
		AverageRating: 0,
		ReviewCount:   0,
	}

	return stats, nil
}

// RefreshProductStatistics is a no-op for MongoDB (no materialized views)
func (r *productRepository) RefreshProductStatistics(ctx context.Context) error {
	// MongoDB doesn't use materialized views, statistics are calculated on-demand
	return nil
}

// getNextProductID gets the next auto-increment ID for products
func (r *productRepository) getNextProductID(ctx context.Context) (int, error) {
	collection := r.db.Collection("products")

	opts := options.Find().SetSort(bson.M{"_id": -1}).SetLimit(1)
	cursor, err := collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return 0, err
	}
	defer cursor.Close(ctx)

	if cursor.Next(ctx) {
		var result domain.Product
		if err := cursor.Decode(&result); err != nil {
			return 0, err
		}
		return result.ID + 1, nil
	}

	return 1, nil
}

// getNextCategoryID gets the next auto-increment ID for categories
func (r *productRepository) getNextCategoryID(ctx context.Context) (int, error) {
	collection := r.db.Collection("categories")

	opts := options.Find().SetSort(bson.M{"_id": -1}).SetLimit(1)
	cursor, err := collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return 0, err
	}
	defer cursor.Close(ctx)

	if cursor.Next(ctx) {
		var result domain.Category
		if err := cursor.Decode(&result); err != nil {
			return 0, err
		}
		return result.ID + 1, nil
	}

	return 1, nil
}
