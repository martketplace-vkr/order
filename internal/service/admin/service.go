package admin

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/martketplace-vkr/order/domain"
	adminrepo "github.com/martketplace-vkr/order/internal/repository/pg/admin"
	"github.com/martketplace-vkr/order/internal/service/ordererrors"
)

type service struct {
	repository repository
}

func New(repository repository) *service {
	return &service{
		repository: repository,
	}
}

func (s *service) ListOrders(ctx context.Context, paymentStatus, fulfillmentStatus string, limit uint32, offset uint64) ([]domain.Order, error) {
	paymentStatus = strings.TrimSpace(paymentStatus)
	fulfillmentStatus = strings.TrimSpace(fulfillmentStatus)
	if paymentStatus != "" {
		if _, err := domain.ParseOrderPaymentStatus(paymentStatus); err != nil {
			return nil, fmt.Errorf("%w: %s", ordererrors.ErrInvalidArgument, err.Error())
		}
	}
	if fulfillmentStatus != "" {
		if _, err := domain.ParseOrderStatus(fulfillmentStatus); err != nil {
			return nil, fmt.Errorf("%w: %s", ordererrors.ErrInvalidArgument, err.Error())
		}
	}

	return s.repository.ListOrders(ctx, adminrepo.ListOrdersFilter{
		PaymentStatus:     paymentStatus,
		FulfillmentStatus: fulfillmentStatus,
		Limit:             limit,
		Offset:            offset,
	})
}

func (s *service) UpdatePaymentStatus(ctx context.Context, orderID int64, paymentStatus string) (*domain.Order, error) {
	if orderID <= 0 {
		return nil, fmt.Errorf("%w: order_id must be greater than zero", ordererrors.ErrInvalidArgument)
	}

	status, err := domain.ParseOrderPaymentStatus(paymentStatus)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ordererrors.ErrInvalidArgument, err.Error())
	}

	order, err := s.repository.UpdatePaymentStatus(ctx, orderID, status)
	if err != nil {
		return nil, mapRepositoryError(err)
	}

	return order, nil
}

func mapRepositoryError(err error) error {
	if errors.Is(err, sql.ErrNoRows) {
		return ordererrors.ErrOrderNotFound
	}

	return err
}
