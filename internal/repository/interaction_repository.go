package repository

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"

	"github.com/PrimeraAizen/e-comm/internal/domain"
	postgres "github.com/PrimeraAizen/e-comm/pkg/adapter"
)

type InteractionRepository interface {
	// View interactions
	RecordView(ctx context.Context, userID, productID string) error
	GetUserViews(ctx context.Context, userID string, limit int) ([]domain.ProductInteraction, error)
	HasViewed(ctx context.Context, userID, productID string) (bool, error)

	// Like interactions
	RecordLike(ctx context.Context, userID, productID string) error
	RemoveLike(ctx context.Context, userID, productID string) error
	GetUserLikes(ctx context.Context, userID string, limit int) ([]domain.ProductInteraction, error)
	HasLiked(ctx context.Context, userID, productID string) (bool, error)

	// Purchase interactions (via orders/order_items)
	RecordPurchase(ctx context.Context, userID, productID string, quantity int, price float64) error
	GetUserPurchases(ctx context.Context, userID string, limit int) ([]domain.ProductInteraction, error)
	HasPurchased(ctx context.Context, userID, productID string) (bool, error)

	// Summary
	GetUserInteractionSummary(ctx context.Context, userID string) (*domain.UserInteractionSummary, error)

	// For recommendations
	GetAllUserViews(ctx context.Context) ([]domain.UserProductView, error)
	GetAllUserLikes(ctx context.Context) ([]domain.UserProductLike, error)
	GetAllUserPurchases(ctx context.Context) ([]domain.UserProductPurchase, error)
}

type interactionRepository struct {
	db *postgres.Postgres
}

func NewInteractionRepository(db *postgres.Postgres) InteractionRepository {
	return &interactionRepository{db: db}
}

func (r *interactionRepository) RecordView(ctx context.Context, userID, productID string) error {
	query, args, err := r.db.Builder.
		Insert("user_product_views").
		Columns("user_id", "product_id").
		Values(userID, productID).
		ToSql()
	if err != nil {
		return err
	}
	_, err = r.db.Pool.Exec(ctx, query, args...)
	return err
}

func (r *interactionRepository) GetUserViews(ctx context.Context, userID string, limit int) ([]domain.ProductInteraction, error) {
	query, args, err := r.db.Builder.
		Select("p.id", "p.name", "COALESCE(p.category_id::text, '')", "p.price", "upv.viewed_at").
		From("user_product_views upv").
		Join("products p ON p.id = upv.product_id").
		Where(squirrel.Eq{"upv.user_id": userID}).
		OrderBy("upv.viewed_at DESC").
		Limit(uint64(limit)).
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("get user views: %w", err)
	}
	defer rows.Close()

	var interactions []domain.ProductInteraction
	for rows.Next() {
		var i domain.ProductInteraction
		if err := rows.Scan(&i.ProductID, &i.ProductName, &i.CategoryID, &i.Price, &i.InteractedAt); err != nil {
			return nil, err
		}
		interactions = append(interactions, i)
	}

	return interactions, rows.Err()
}

func (r *interactionRepository) HasViewed(ctx context.Context, userID, productID string) (bool, error) {
	query, args, err := r.db.Builder.
		Select("COUNT(*)").
		From("user_product_views").
		Where(squirrel.Eq{"user_id": userID, "product_id": productID}).
		ToSql()
	if err != nil {
		return false, err
	}

	var count int
	if err := r.db.Pool.QueryRow(ctx, query, args...).Scan(&count); err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *interactionRepository) RecordLike(ctx context.Context, userID, productID string) error {
	query, args, err := r.db.Builder.
		Insert("user_product_likes").
		Columns("user_id", "product_id").
		Values(userID, productID).
		Suffix("ON CONFLICT DO NOTHING").
		ToSql()
	if err != nil {
		return err
	}
	_, err = r.db.Pool.Exec(ctx, query, args...)
	return err
}

func (r *interactionRepository) RemoveLike(ctx context.Context, userID, productID string) error {
	query, args, err := r.db.Builder.
		Delete("user_product_likes").
		Where(squirrel.Eq{"user_id": userID, "product_id": productID}).
		ToSql()
	if err != nil {
		return err
	}

	rows, err := r.db.Pool.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("remove like: %w", err)
	}
	if rows.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *interactionRepository) GetUserLikes(ctx context.Context, userID string, limit int) ([]domain.ProductInteraction, error) {
	query, args, err := r.db.Builder.
		Select("p.id", "p.name", "COALESCE(p.category_id::text, '')", "p.price", "upl.created_at").
		From("user_product_likes upl").
		Join("products p ON p.id = upl.product_id").
		Where(squirrel.Eq{"upl.user_id": userID}).
		OrderBy("upl.created_at DESC").
		Limit(uint64(limit)).
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("get user likes: %w", err)
	}
	defer rows.Close()

	var interactions []domain.ProductInteraction
	for rows.Next() {
		var i domain.ProductInteraction
		if err := rows.Scan(&i.ProductID, &i.ProductName, &i.CategoryID, &i.Price, &i.InteractedAt); err != nil {
			return nil, err
		}
		interactions = append(interactions, i)
	}

	return interactions, rows.Err()
}

func (r *interactionRepository) HasLiked(ctx context.Context, userID, productID string) (bool, error) {
	query, args, err := r.db.Builder.
		Select("COUNT(*)").
		From("user_product_likes").
		Where(squirrel.Eq{"user_id": userID, "product_id": productID}).
		ToSql()
	if err != nil {
		return false, err
	}

	var count int
	if err := r.db.Pool.QueryRow(ctx, query, args...).Scan(&count); err != nil {
		return false, err
	}
	return count > 0, nil
}

// RecordPurchase creates an order and order item to record a purchase.
func (r *interactionRepository) RecordPurchase(ctx context.Context, userID, productID string, quantity int, price float64) error {
	totalAmount := price * float64(quantity)

	orderQuery, orderArgs, err := r.db.Builder.
		Insert("orders").
		Columns("user_id", "status", "total_amount", "payment_status").
		Values(userID, "completed", totalAmount, "paid").
		Suffix("RETURNING id").
		ToSql()
	if err != nil {
		return err
	}

	var orderID string
	if err := r.db.Pool.QueryRow(ctx, orderQuery, orderArgs...).Scan(&orderID); err != nil {
		return fmt.Errorf("create order: %w", err)
	}

	itemQuery, itemArgs, err := r.db.Builder.
		Insert("order_items").
		Columns("order_id", "product_id", "quantity", "price_at_purchase", "subtotal").
		Values(orderID, productID, quantity, price, totalAmount).
		ToSql()
	if err != nil {
		return err
	}

	_, err = r.db.Pool.Exec(ctx, itemQuery, itemArgs...)
	return err
}

func (r *interactionRepository) GetUserPurchases(ctx context.Context, userID string, limit int) ([]domain.ProductInteraction, error) {
	query, args, err := r.db.Builder.
		Select("p.id", "p.name", "COALESCE(p.category_id::text, '')", "p.price", "oi.created_at").
		From("order_items oi").
		Join("orders o ON o.id = oi.order_id").
		Join("products p ON p.id = oi.product_id").
		Where(squirrel.Eq{"o.user_id": userID}).
		OrderBy("oi.created_at DESC").
		Limit(uint64(limit)).
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("get user purchases: %w", err)
	}
	defer rows.Close()

	var interactions []domain.ProductInteraction
	for rows.Next() {
		var i domain.ProductInteraction
		if err := rows.Scan(&i.ProductID, &i.ProductName, &i.CategoryID, &i.Price, &i.InteractedAt); err != nil {
			return nil, err
		}
		interactions = append(interactions, i)
	}

	return interactions, rows.Err()
}

func (r *interactionRepository) HasPurchased(ctx context.Context, userID, productID string) (bool, error) {
	var exists bool
	err := r.db.Pool.QueryRow(ctx,
		"SELECT EXISTS(SELECT 1 FROM order_items oi JOIN orders o ON o.id = oi.order_id WHERE o.user_id = $1 AND oi.product_id = $2)",
		userID, productID,
	).Scan(&exists)
	return exists, err
}

func (r *interactionRepository) GetUserInteractionSummary(ctx context.Context, userID string) (*domain.UserInteractionSummary, error) {
	views, err := r.GetUserViews(ctx, userID, 50)
	if err != nil {
		return nil, err
	}

	likes, err := r.GetUserLikes(ctx, userID, 50)
	if err != nil {
		return nil, err
	}

	purchases, err := r.GetUserPurchases(ctx, userID, 50)
	if err != nil {
		return nil, err
	}

	var totalViews int64
	r.db.Pool.QueryRow(ctx, "SELECT COUNT(*) FROM user_product_views WHERE user_id = $1", userID).Scan(&totalViews)

	var totalLikes int64
	r.db.Pool.QueryRow(ctx, "SELECT COUNT(*) FROM user_product_likes WHERE user_id = $1", userID).Scan(&totalLikes)

	var totalPurchases int64
	r.db.Pool.QueryRow(ctx, "SELECT COUNT(oi.id) FROM order_items oi JOIN orders o ON o.id = oi.order_id WHERE o.user_id = $1", userID).Scan(&totalPurchases)

	return &domain.UserInteractionSummary{
		UserID:            userID,
		ViewedProducts:    views,
		LikedProducts:     likes,
		PurchasedProducts: purchases,
		TotalViews:        totalViews,
		TotalLikes:        totalLikes,
		TotalPurchases:    totalPurchases,
	}, nil
}

func (r *interactionRepository) GetAllUserViews(ctx context.Context) ([]domain.UserProductView, error) {
	query, args, err := r.db.Builder.
		Select("user_id", "product_id", "viewed_at").
		From("user_product_views").
		OrderBy("viewed_at DESC").
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("get all views: %w", err)
	}
	defer rows.Close()

	var views []domain.UserProductView
	for rows.Next() {
		var v domain.UserProductView
		if err := rows.Scan(&v.UserID, &v.ProductID, &v.ViewedAt); err != nil {
			return nil, err
		}
		views = append(views, v)
	}

	return views, rows.Err()
}

func (r *interactionRepository) GetAllUserLikes(ctx context.Context) ([]domain.UserProductLike, error) {
	query, args, err := r.db.Builder.
		Select("user_id", "product_id", "created_at").
		From("user_product_likes").
		OrderBy("created_at DESC").
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("get all likes: %w", err)
	}
	defer rows.Close()

	var likes []domain.UserProductLike
	for rows.Next() {
		var l domain.UserProductLike
		if err := rows.Scan(&l.UserID, &l.ProductID, &l.LikedAt); err != nil {
			return nil, err
		}
		likes = append(likes, l)
	}

	return likes, rows.Err()
}

func (r *interactionRepository) GetAllUserPurchases(ctx context.Context) ([]domain.UserProductPurchase, error) {
	rows, err := r.db.Pool.Query(ctx,
		"SELECT o.user_id, oi.product_id, oi.quantity, oi.price_at_purchase, oi.created_at FROM order_items oi JOIN orders o ON o.id = oi.order_id ORDER BY oi.created_at DESC",
	)
	if err != nil {
		return nil, fmt.Errorf("get all purchases: %w", err)
	}
	defer rows.Close()

	var purchases []domain.UserProductPurchase
	for rows.Next() {
		var p domain.UserProductPurchase
		if err := rows.Scan(&p.UserID, &p.ProductID, &p.Quantity, &p.PriceAtPurchase, &p.PurchasedAt); err != nil {
			return nil, err
		}
		purchases = append(purchases, p)
	}

	return purchases, rows.Err()
}
