package client

import (
	"context"
	"database/sql"
	"errors"

	cart "github.com/martketplace-vkr/cart/pkg/api/grpc/v1"
	"github.com/martketplace-vkr/order/domain"
	"github.com/martketplace-vkr/order/internal/service/ordererrors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type service struct {
	repository repository
	outbox     outbox
	cart       *cart.Connector
}

func New(
	repository repository,
	cartClient *cart.Connector,
	outbox outbox,
) *service {
	return &service{
		repository: repository,
		cart:       cartClient,
		outbox:     outbox,
	}
}

func (s *service) GetOrder(
	ctx context.Context,
	userID int64,
	orderID int64,
) (*domain.Order, error) {
	order, err := s.repository.GetOrder(ctx, userID, orderID)
	if err != nil {
		return nil, mapRepositoryError(err)
	}

	return order, nil
}

func (s *service) GetOrderList(
	ctx context.Context,
	userID int64,
) ([]domain.Order, error) {

	orders, err := s.repository.GetOrderList(ctx, userID)
	if err != nil {
		return nil, mapRepositoryError(err)
	}

	return orders, nil
}

func (s *service) HasSuccessfulProductOrder(
	ctx context.Context,
	userID int64,
	productID int64,
) (bool, int64, error) {
	if userID <= 0 || productID <= 0 {
		return false, 0, ordererrors.ErrInvalidArgument
	}

	hasOrder, vendorID, err := s.repository.HasSuccessfulProductOrder(ctx, userID, productID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, 0, nil
		}

		return false, 0, mapRepositoryError(err)
	}

	return hasOrder, vendorID, nil
}

func (s *service) CancelOrder(
	ctx context.Context,
	userID int64,
	orderID int64,
) (*domain.Order, error) {
	order, err := s.repository.CancelOrder(ctx, userID, orderID)
	if err != nil {
		return nil, mapRepositoryError(err)
	}

	return order, nil
}

func mapRepositoryError(err error) error {
	if errors.Is(err, sql.ErrNoRows) {
		return ordererrors.ErrOrderNotFound
	}
	if status.Code(err) == codes.Unknown {
		return status.Error(codes.Internal, err.Error())
	}

	return err
}
