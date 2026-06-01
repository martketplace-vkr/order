package dto

import "github.com/martketplace-vkr/order/domain"

type CheckoutRequest struct {
	UserID              int64
	CheckoutID          string
	ProductIDs          []int64
	ExpectedCartVersion uint64
	PreferredCurrencyID int64
	Payment             domain.Payment
	Delivery            domain.Delivery
}
