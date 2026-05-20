package admin

import (
	"context"

	"github.com/martketplace-vkr/order/domain"
	adminrepo "github.com/martketplace-vkr/order/internal/repository/pg/admin"
)

type repository interface {
	ListOrders(ctx context.Context, filter adminrepo.ListOrdersFilter) ([]domain.Order, error)
	UpdatePaymentStatus(ctx context.Context, orderID int64, paymentStatus domain.OrderPaymentStatus) (*domain.Order, error)
}
