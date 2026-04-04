package client

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/martketplace-vkr/order/domain"
	"github.com/martketplace-vkr/order/internal/service/ordererrors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type service struct {
	repository  repository
	cartClient  cartClient
	cartTimeout time.Duration
}

func New(repository repository, cartClient cartClient, cartTimeout time.Duration) *service {
	return &service{
		repository:  repository,
		cartClient:  cartClient,
		cartTimeout: cartTimeout,
	}
}

func (s *service) GetOrder(
	ctx context.Context,
	userID int64,
	orderID int64,
) (*domain.Order, error) {
	if err := validateID("user_id", userID); err != nil {
		return nil, err
	}
	if err := validateID("order_id", orderID); err != nil {
		return nil, err
	}

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
	if err := validateID("user_id", userID); err != nil {
		return nil, err
	}

	orders, err := s.repository.GetOrderList(ctx, userID)
	if err != nil {
		return nil, mapRepositoryError(err)
	}

	return orders, nil
}

func (s *service) CancelOrder(
	ctx context.Context,
	userID int64,
	orderID int64,
) (*domain.Order, error) {
	if err := validateID("user_id", userID); err != nil {
		return nil, err
	}
	if err := validateID("order_id", orderID); err != nil {
		return nil, err
	}

	order, err := s.repository.CancelOrder(ctx, userID, orderID)
	if err != nil {
		return nil, mapRepositoryError(err)
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
	if status.Code(err) == codes.Unknown {
		return status.Error(codes.Internal, err.Error())
	}

	return err
}
