package db

import (
	"context"
	"database/sql"
	"time"
)

// CheckProductLike checks if the user has liked the product.
func CheckProductLike(ctx context.Context, db *sql.DB, userID, productID uint64) (bool, error) {
	var exists bool
	query := `
        SELECT EXISTS (
            SELECT 1
            FROM favorites
            WHERE user_id = $1 AND product_id = $2
        )
    `
	err := db.QueryRowContext(ctx, query, userID, productID).Scan(&exists)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return exists, nil
}

// UnlikeProduct removes the like for a product.
func UnlikeProduct(ctx context.Context, db *sql.DB, userID, productID uint64) (int64, error) {
	query := `
        DELETE FROM favorites
        WHERE user_id = $1 AND product_id = $2
    `
	result, err := db.ExecContext(ctx, query, userID, productID)
	if err != nil {
		return 0, err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	return rowsAffected, nil
}

type Product struct {
	ID                uint64    `db:"id"`
	ShopID            uint64    `db:"shop_id"`
	Name              string    `db:"name"`
	Description       string    `db:"description"`
	ThumbnailURL      string    `db:"thumbnail_url"`
	OriginPrice       int64     `db:"origin_price"`
	DiscountedPrice   int64     `db:"discounted_price"`
	DiscountedRate    float64   `db:"discounted_rate"`
	Status            string    `db:"status"`
	InStock           bool      `db:"in_stock"`
	IsPreorder        bool      `db:"is_preorder"`
	IsPurchasable     bool      `db:"is_purchasable"`
	DeliveryCondition string    `db:"delivery_condition"`
	DeliveryDisplay   string    `db:"delivery_display"`
	CreatedAt         time.Time `db:"created_at"`
	UpdatedAt         time.Time `db:"updated_at"`
}

// RetrieveLikedProducts retrieves liked products with pagination.
func RetrieveLikedProducts(ctx context.Context, db *sql.DB, userID, page, limit int) ([]Product, int, error) {
	var products []Product
	query := `
        SELECT p.*
        FROM products p
        INNER JOIN (
            SELECT product_id
            FROM favorites
            WHERE user_id = $1
            ORDER BY created_at
            LIMIT $2
            OFFSET ($3 - 1) * $2
        ) AS liked_products ON p.id = liked_products.product_id
    `
	rows, err := db.QueryContext(ctx, query, userID, limit, page)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var product Product
		if err := rows.Scan(
			&product.ID, &product.ShopID, &product.Name, // Add other fields here
		); err != nil {
			return nil, 0, err
		}
		products = append(products, product)
	}

	totalCount, err := RetrieveTotalLikedProductCount(ctx, db, userID)
	if err != nil {
		return nil, 0, err
	}

	return products, totalCount, nil
}

// RetrieveTotalLikedProductCount retrieves the total count of liked products for a user.
func RetrieveTotalLikedProductCount(ctx context.Context, db *sql.DB, userID int) (int, error) {
	var totalCount int
	query := `
        SELECT COUNT(*)
        FROM favorites
        WHERE user_id = $1
    `
	err := db.QueryRowContext(ctx, query, userID).Scan(&totalCount)
	if err != nil {
		return 0, err
	}
	return totalCount, nil
}

// AddProductLike adds a product like (inserts a record into favorites).
func AddProductLike(ctx context.Context, db *sql.DB, userID, productID uint64) (int64, error) {
	query := `
        INSERT INTO favorites (user_id, product_id, created_at)
        VALUES ($1, $2, NOW())
        RETURNING id
    `
	var id int64
	err := db.QueryRowContext(ctx, query, userID, productID).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}
