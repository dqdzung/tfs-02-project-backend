package model

import "time"

type OrderDetail struct {
	Id            int     `json:"id"`
	ProductName   string  `json: "product_name"`
	VariantName   string  `json:"variant_name"`
	Price         float64 `json:"price"`
	OriginalPrice float64 `json:"original_price"`
	Quantity      int     `json:"quantity"`
	Weight        string  `json:"weight"`

	VariantId int `json:"variant_id"`
	OrderId   int `json: "order_id"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
