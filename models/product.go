// models/user.go

package models

import (
	"time"
)

// Product represents the product table in the database.
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

// // NewProduct creates a new Product instance.
// func NewProduct(shopID int64, name, description, thumbnailURL string, originPrice, discountedPrice int64, discountedRate float64, status string, inStock, isPreorder, isPurchasable bool, deliveryCondition, deliveryDisplay string) *Product {
//     return &Product{
//         ShopID:           shopID,
//         Name:             name,
//         Description:      description,
//         ThumbnailURL:     thumbnailURL,
//         OriginPrice:      originPrice,
//         DiscountedPrice:  discountedPrice,
//         DiscountedRate:   discountedRate,
//         Status:           status,
//         InStock:          inStock,
//         IsPreorder:       isPreorder,
//         IsPurchasable:    isPurchasable,
//         DeliveryCondition: deliveryCondition,
//         DeliveryDisplay:  deliveryDisplay,
//         CreatedAt:        time.Now(),
//         UpdatedAt:        time.Now(),
//     }
// }
