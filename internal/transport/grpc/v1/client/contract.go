package client

import (
	"context"

	"github.com/martketplace-vkr/order/domain"
	"github.com/martketplace-vkr/order/internal/service/client/dto"
)

type (
	service interface {
		Checkout(
			ctx context.Context,
			req dto.CheckoutRequest,
		) ([]domain.Order, error)
		GetOrder(ctx context.Context, userID int64, orderID int64) (*domain.Order, error)
		GetOrderList(ctx context.Context, userID int64) ([]domain.Order, error)
		CancelOrder(ctx context.Context, userID int64, orderID int64) (*domain.Order, error)
	}
)
