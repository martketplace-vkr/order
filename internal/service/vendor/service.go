package vendor

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/martketplace-vkr/order/domain"
	"github.com/martketplace-vkr/order/internal/service/ordererrors"
)

type service struct {
	repository repository
	outbox     outbox
}

func New(repository repository, outbox outbox) *service {
	return &service{
		repository: repository,
		outbox:     outbox,
	}
}

func (s *service) GetOrder(
	ctx context.Context,
	vendorID int64,
	orderID int64,
) (*domain.Order, error) {
	if err := validateID("vendor_id", vendorID); err != nil {
		return nil, err
	}
	if err := validateID("order_id", orderID); err != nil {
		return nil, err
	}

	order, err := s.repository.GetOrder(ctx, vendorID, orderID)
	if err != nil {
		return nil, mapRepositoryError(err)
	}

	return order, nil
}

func (s *service) GetOrderList(
	ctx context.Context,
	vendorID int64,
) ([]domain.Order, error) {
	if err := validateID("vendor_id", vendorID); err != nil {
		return nil, err
	}

	orders, err := s.repository.GetOrderList(ctx, vendorID)
	if err != nil {
		return nil, mapRepositoryError(err)
	}

	return orders, nil
}

func (s *service) UpdateOrder(
	ctx context.Context,
	vendorID int64,
	orderID int64,
	status string,
) (*domain.Order, error) {
	if err := validateID("vendor_id", vendorID); err != nil {
		return nil, err
	}
	if err := validateID("order_id", orderID); err != nil {
		return nil, err
	}

	status = strings.TrimSpace(status)
	if status == "" {
		return nil, fmt.Errorf("%w: status must not be empty", ordererrors.ErrInvalidArgument)
	}
	fulfillmentStatus, err := domain.ParseOrderStatus(status)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ordererrors.ErrInvalidArgument, err.Error())
	}
	if _, ok := domain.ValidVendorFulfillmentStatuses[fulfillmentStatus]; !ok {
		return nil, fmt.Errorf("%w: fulfillment status is not allowed for vendor", ordererrors.ErrInvalidArgument)
	}

	order, err := s.repository.UpdateOrder(ctx, vendorID, orderID, status)
	if err != nil {
		return nil, mapRepositoryError(err)
	}

	if s.outbox != nil {
		switch fulfillmentStatus {
		case domain.Success:
			if err := s.outbox.SendOrderPickedUp(ctx, *order); err != nil {
				return nil, err
			}
		case domain.CancelledBySeller:
			if err := s.outbox.SendOrderCancelled(ctx, *order); err != nil {
				return nil, err
			}
		}
	}

	return order, nil
}

func validateID(field string, value int64) error {
	if value <= 0 {
		return fmt.Errorf("%w: %s must be greater than zero", ordererrors.ErrInvalidArgument, field)
	}

	return nil
}

func mapRepositoryError(err error) error {
	if errors.Is(err, sql.ErrNoRows) {
		return ordererrors.ErrOrderNotFound
	}

	return err
}
