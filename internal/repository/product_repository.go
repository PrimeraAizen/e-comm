package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"

	"github.com/PrimeraAizen/e-comm/internal/domain"
	postgres "github.com/PrimeraAizen/e-comm/pkg/adapter"
)

type ProductRepository interface {
	// Product CRUD
	Create(ctx context.Context, product *domain.Product) error
	GetByID(ctx context.Context, id string) (*domain.Product, error)
	GetByIDWithCategory(ctx context.Context, id string) (*domain.ProductWithCategory, error)
	Update(ctx context.Context, product *domain.Product) error
	Delete(ctx context.Context, id string) error

	// Product listing and search
	List(ctx context.Context, filter domain.ProductFilter) ([]*domain.Product, int64, error)
	ListWithCategories(ctx context.Context, filter domain.ProductFilter) ([]*domain.ProductWithCategory, int64, error)
	Search(ctx context.Context, query string, limit, offset int) ([]*domain.Product, int64, error)

	// Category CRUD
	CreateCategory(ctx context.Context, category *domain.Category) error
	GetCategoryByID(ctx context.Context, id string) (*domain.Category, error)
	GetCategoryByName(ctx context.Context, name string) (*domain.Category, error)
	ListCategories(ctx context.Context) ([]*domain.Category, error)
	UpdateCategory(ctx context.Context, category *domain.Category) error
	DeleteCategory(ctx context.Context, id string) error

	// Product statistics
	GetProductStatistics(ctx context.Context, productID string) (*domain.ProductStatistics, error)
	RefreshProductStatistics(ctx context.Context) error
}

type productRepository struct {
	db *postgres.Postgres
}

func NewProductRepository(db *postgres.Postgres) ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) Create(ctx context.Context, product *domain.Product) error {
	now := time.Now()
	product.IsActive = true
	query, args, err := r.db.Builder.
		Insert("products").
		Columns("name", "description", "category_id", "price", "stock", "image_url", "is_active", "created_at", "updated_at").
		Values(product.Name, product.Description, product.CategoryID, product.Price, product.Stock, product.ImageURL, product.IsActive, now, now).
		Suffix("RETURNING id").
		ToSql()
	if err != nil {
		return fmt.Errorf("create product: %w", err)
	}
	return r.db.Pool.QueryRow(ctx, query, args...).Scan(&product.ID)
}

func (r *productRepository) GetByID(ctx context.Context, id string) (*domain.Product, error) {
	query, args, err := r.db.Builder.
		Select("id", "name", "description", "category_id", "price", "stock", "image_url", "is_active", "created_at", "updated_at").
		From("products").
		Where(squirrel.Eq{"id": id}).
		ToSql()
	if err != nil {
		return nil, err
	}

	var p domain.Product
	err = r.db.Pool.QueryRow(ctx, query, args...).Scan(
		&p.ID, &p.Name, &p.Description, &p.CategoryID, &p.Price, &p.Stock, &p.ImageURL, &p.IsActive, &p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("get product by id: %w", err)
	}

	return &p, nil
}

func (r *productRepository) GetByIDWithCategory(ctx context.Context, id string) (*domain.ProductWithCategory, error) {
	query, args, err := r.db.Builder.
		Select("p.id", "p.name", "p.description", "p.category_id", "p.price", "p.stock", "p.image_url", "p.is_active", "p.created_at", "p.updated_at", "COALESCE(c.name, '')").
		From("products p").
		LeftJoin("categories c ON c.id = p.category_id").
		Where(squirrel.Eq{"p.id": id}).
		ToSql()
	if err != nil {
		return nil, err
	}

	var p domain.ProductWithCategory
	err = r.db.Pool.QueryRow(ctx, query, args...).Scan(
		&p.ID, &p.Name, &p.Description, &p.CategoryID, &p.Price, &p.Stock, &p.ImageURL, &p.IsActive, &p.CreatedAt, &p.UpdatedAt, &p.CategoryName,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("get product with category: %w", err)
	}

	return &p, nil
}

func (r *productRepository) Update(ctx context.Context, product *domain.Product) error {
	now := time.Now()
	query, args, err := r.db.Builder.
		Update("products").
		Set("name", product.Name).
		Set("description", product.Description).
		Set("category_id", product.CategoryID).
		Set("price", product.Price).
		Set("stock", product.Stock).
		Set("image_url", product.ImageURL).
		Set("is_active", product.IsActive).
		Set("updated_at", now).
		Where(squirrel.Eq{"id": product.ID}).
		ToSql()
	if err != nil {
		return err
	}

	rows, err := r.db.Pool.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("update product: %w", err)
	}
	if rows.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *productRepository) Delete(ctx context.Context, id string) error {
	query, args, err := r.db.Builder.
		Delete("products").
		Where(squirrel.Eq{"id": id}).
		ToSql()
	if err != nil {
		return err
	}

	rows, err := r.db.Pool.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("delete product: %w", err)
	}
	if rows.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *productRepository) applyProductFilter(sb squirrel.SelectBuilder, filter domain.ProductFilter, prefix string) squirrel.SelectBuilder {
	if prefix != "" {
		prefix += "."
	}
	if filter.CategoryID != nil {
		sb = sb.Where(squirrel.Eq{prefix + "category_id": *filter.CategoryID})
	}
	if filter.MinPrice != nil {
		sb = sb.Where(squirrel.GtOrEq{prefix + "price": *filter.MinPrice})
	}
	if filter.MaxPrice != nil {
		sb = sb.Where(squirrel.LtOrEq{prefix + "price": *filter.MaxPrice})
	}
	if filter.IsActive != nil {
		sb = sb.Where(squirrel.Eq{prefix + "is_active": *filter.IsActive})
	}
	if filter.SearchQuery != "" {
		sb = sb.Where("("+prefix+"name ILIKE ? OR "+prefix+"description ILIKE ?)", "%"+filter.SearchQuery+"%", "%"+filter.SearchQuery+"%")
	}
	return sb
}

func (r *productRepository) List(ctx context.Context, filter domain.ProductFilter) ([]*domain.Product, int64, error) {
	countSb := r.applyProductFilter(r.db.Builder.Select("COUNT(*)").From("products"), filter, "")
	countQuery, countArgs, err := countSb.ToSql()
	if err != nil {
		return nil, 0, err
	}
	var total int64
	if err := r.db.Pool.QueryRow(ctx, countQuery, countArgs...).Scan(&total); err != nil {
		return nil, 0, err
	}

	sortField := "created_at"
	if filter.SortBy != "" {
		sortField = filter.SortBy
	}
	sortOrder := "DESC"
	if filter.SortOrder == "asc" {
		sortOrder = "ASC"
	}

	sb := r.applyProductFilter(
		r.db.Builder.Select("id", "name", "description", "category_id", "price", "stock", "image_url", "is_active", "created_at", "updated_at").From("products"),
		filter, "",
	).OrderBy(sortField + " " + sortOrder)

	if filter.Limit > 0 {
		sb = sb.Limit(uint64(filter.Limit))
	}
	if filter.Offset > 0 {
		sb = sb.Offset(uint64(filter.Offset))
	}

	query, args, err := sb.ToSql()
	if err != nil {
		return nil, 0, err
	}

	rows, err := r.db.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("list products: %w", err)
	}
	defer rows.Close()

	var products []*domain.Product
	for rows.Next() {
		var p domain.Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.CategoryID, &p.Price, &p.Stock, &p.ImageURL, &p.IsActive, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, 0, err
		}
		products = append(products, &p)
	}

	return products, total, rows.Err()
}

func (r *productRepository) ListWithCategories(ctx context.Context, filter domain.ProductFilter) ([]*domain.ProductWithCategory, int64, error) {
	countSb := r.applyProductFilter(r.db.Builder.Select("COUNT(*)").From("products p"), filter, "p")
	countQuery, countArgs, err := countSb.ToSql()
	if err != nil {
		return nil, 0, err
	}
	var total int64
	if err := r.db.Pool.QueryRow(ctx, countQuery, countArgs...).Scan(&total); err != nil {
		return nil, 0, err
	}

	sortField := "p.created_at"
	if filter.SortBy != "" {
		sortField = "p." + filter.SortBy
	}
	sortOrder := "DESC"
	if filter.SortOrder == "asc" {
		sortOrder = "ASC"
	}

	sb := r.applyProductFilter(
		r.db.Builder.
			Select("p.id", "p.name", "p.description", "p.category_id", "p.price", "p.stock", "p.image_url", "p.is_active", "p.created_at", "p.updated_at", "COALESCE(c.name, '')").
			From("products p").
			LeftJoin("categories c ON c.id = p.category_id"),
		filter, "p",
	).OrderBy(sortField + " " + sortOrder)

	if filter.Limit > 0 {
		sb = sb.Limit(uint64(filter.Limit))
	}
	if filter.Offset > 0 {
		sb = sb.Offset(uint64(filter.Offset))
	}

	query, args, err := sb.ToSql()
	if err != nil {
		return nil, 0, err
	}

	rows, err := r.db.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("list products with categories: %w", err)
	}
	defer rows.Close()

	var products []*domain.ProductWithCategory
	for rows.Next() {
		var p domain.ProductWithCategory
		if err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.CategoryID, &p.Price, &p.Stock, &p.ImageURL, &p.IsActive, &p.CreatedAt, &p.UpdatedAt, &p.CategoryName); err != nil {
			return nil, 0, err
		}
		products = append(products, &p)
	}

	return products, total, rows.Err()
}

func (r *productRepository) Search(ctx context.Context, query string, limit, offset int) ([]*domain.Product, int64, error) {
	return r.List(ctx, domain.ProductFilter{
		SearchQuery: query,
		Limit:       limit,
		Offset:      offset,
	})
}

func (r *productRepository) CreateCategory(ctx context.Context, category *domain.Category) error {
	now := time.Now()
	query, args, err := r.db.Builder.
		Insert("categories").
		Columns("name", "description", "parent_id", "created_at", "updated_at").
		Values(category.Name, category.Description, category.ParentID, now, now).
		Suffix("RETURNING id").
		ToSql()
	if err != nil {
		return err
	}
	return r.db.Pool.QueryRow(ctx, query, args...).Scan(&category.ID)
}

func (r *productRepository) GetCategoryByID(ctx context.Context, id string) (*domain.Category, error) {
	query, args, err := r.db.Builder.
		Select("id", "name", "description", "parent_id", "created_at", "updated_at").
		From("categories").
		Where(squirrel.Eq{"id": id}).
		ToSql()
	if err != nil {
		return nil, err
	}

	var c domain.Category
	err = r.db.Pool.QueryRow(ctx, query, args...).Scan(&c.ID, &c.Name, &c.Description, &c.ParentID, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("get category by id: %w", err)
	}

	return &c, nil
}

func (r *productRepository) GetCategoryByName(ctx context.Context, name string) (*domain.Category, error) {
	query, args, err := r.db.Builder.
		Select("id", "name", "description", "parent_id", "created_at", "updated_at").
		From("categories").
		Where(squirrel.Eq{"name": name}).
		ToSql()
	if err != nil {
		return nil, err
	}

	var c domain.Category
	err = r.db.Pool.QueryRow(ctx, query, args...).Scan(&c.ID, &c.Name, &c.Description, &c.ParentID, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("get category by name: %w", err)
	}

	return &c, nil
}

func (r *productRepository) ListCategories(ctx context.Context) ([]*domain.Category, error) {
	query, args, err := r.db.Builder.
		Select("id", "name", "description", "parent_id", "created_at", "updated_at").
		From("categories").
		OrderBy("name ASC").
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("list categories: %w", err)
	}
	defer rows.Close()

	var categories []*domain.Category
	for rows.Next() {
		var c domain.Category
		if err := rows.Scan(&c.ID, &c.Name, &c.Description, &c.ParentID, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		categories = append(categories, &c)
	}

	return categories, rows.Err()
}

func (r *productRepository) UpdateCategory(ctx context.Context, category *domain.Category) error {
	now := time.Now()
	query, args, err := r.db.Builder.
		Update("categories").
		Set("name", category.Name).
		Set("description", category.Description).
		Set("parent_id", category.ParentID).
		Set("updated_at", now).
		Where(squirrel.Eq{"id": category.ID}).
		ToSql()
	if err != nil {
		return err
	}

	rows, err := r.db.Pool.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("update category: %w", err)
	}
	if rows.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *productRepository) DeleteCategory(ctx context.Context, id string) error {
	query, args, err := r.db.Builder.
		Delete("categories").
		Where(squirrel.Eq{"id": id}).
		ToSql()
	if err != nil {
		return err
	}

	rows, err := r.db.Pool.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("delete category: %w", err)
	}
	if rows.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *productRepository) GetProductStatistics(ctx context.Context, productID string) (*domain.ProductStatistics, error) {
	query, args, err := r.db.Builder.
		Select("product_id", "product_name", "view_count", "like_count", "purchase_count", "average_rating", "review_count").
		From("product_statistics").
		Where(squirrel.Eq{"product_id": productID}).
		ToSql()
	if err != nil {
		return nil, err
	}

	var s domain.ProductStatistics
	err = r.db.Pool.QueryRow(ctx, query, args...).Scan(
		&s.ProductID, &s.ProductName, &s.ViewCount, &s.LikeCount, &s.PurchaseCount, &s.AverageRating, &s.ReviewCount,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("get product statistics: %w", err)
	}

	return &s, nil
}

func (r *productRepository) RefreshProductStatistics(ctx context.Context) error {
	_, err := r.db.Pool.Exec(ctx, "REFRESH MATERIALIZED VIEW CONCURRENTLY product_statistics")
	if err != nil {
		return fmt.Errorf("refresh product statistics: %w", err)
	}
	return nil
}
