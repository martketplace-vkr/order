package domain

import "time"

type Order struct {
	ID              int64     `db:"id"`
	CheckoutID      string    `db:"checkout_id"`
	UserID          int64     `db:"user_id"`
	VendorID        int64     `db:"vendor_id"`
	Status          string    `db:"status"`
	ProductID       int64     `db:"product_id"`
	ProductName     string    `db:"product_name"`
	ProductImageURL string    `db:"product_image_url"`
	Quantity        int64     `db:"quantity"`
	UnitPrice       string    `db:"unit_price"`
	TotalPrice      string    `db:"total_price"`
	Comment         string    `db:"-"`
	CreatedAt       time.Time `db:"created_at"`
	UpdatedAt       time.Time `db:"updated_at"`
}
