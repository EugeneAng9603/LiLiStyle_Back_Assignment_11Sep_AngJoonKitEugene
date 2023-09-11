package models

import (
	"time"
)

type Product struct {
	ID                uint64    `json:"id"`
	ShopID            int64     `json:"shop_id"`
	Name              string    `json:"name"`
	Description       string    `json:"description"`
	ThumbnailURL      string    `json:"thumbnail_url"`
	OriginPrice       int64     `json:"origin_price"`
	DiscountedPrice   int64     `json:"discounted_price"`
	DiscountedRate    float64   `json:"discounted_rate"`
	Status            string    `json:"status"`
	InStock           bool      `json:"in_stock"`
	IsPreorder        bool      `json:"is_preorder"`
	IsPurchasable     bool      `json:"is_purchasable"`
	DeliveryCondition string    `json:"delivery_condition"`
	DeliveryDisplay   string    `json:"delivery_display"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}
