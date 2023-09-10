package models

import (
	"time"
)

// Favorite represents the favorites table in the database.
type Favorite struct {
	ID        uint64    `json:"id"`
	UserID    int64     `json:"user_id"`
	ProductID int64     `json:"product_id"`
	CreatedAt time.Time `json:"created_at"`
}
