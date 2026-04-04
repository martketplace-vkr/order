package vendor

import (
	"context"

	"github.com/martketplace-vkr/order/domain"
)

type (
	service interface {
		GetOrder(ctx context.Context, vendorID int64, orderID int64) (*domain.Order, error)
		GetOrderList(ctx context.Context, vendorID int64) ([]domain.Order, error)
		UpdateOrder(ctx context.Context, vendorID int64, orderID int64, status string) (*domain.Order, error)
	}
)
