package client

import (
	"context"

	"github.com/martketplace-vkr/order/domain"
)

type (
	service interface {
		Checkout(ctx context.Context, userID int64, checkoutID string, productIDs []int64, expectedCartVersion uint64) ([]domain.Order, error)
		GetOrder(ctx context.Context, userID int64, orderID int64) (*domain.Order, error)
		GetOrderList(ctx context.Context, userID int64) ([]domain.Order, error)
		CancelOrder(ctx context.Context, userID int64, orderID int64) (*domain.Order, error)
	}
)
