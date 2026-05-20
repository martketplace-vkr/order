package admin

import (
	"context"

	"github.com/martketplace-vkr/order/domain"
)

type service interface {
	ListOrders(ctx context.Context, paymentStatus, fulfillmentStatus string, limit uint32, offset uint64) ([]domain.Order, error)
	UpdatePaymentStatus(ctx context.Context, orderID int64, paymentStatus string) (*domain.Order, error)
}
