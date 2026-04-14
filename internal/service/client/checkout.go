package client

import (
	"context"

	cartorderpb "github.com/martketplace-vkr/cart/pkg/api/grpc/v1/order"
	"github.com/martketplace-vkr/order/domain"
	"github.com/martketplace-vkr/order/internal/service/client/dto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *service) Checkout(
	ctx context.Context,
	req dto.CheckoutRequest,
) ([]domain.Order, error) {

	existingOrders, err := s.repository.GetOrdersByCheckout(ctx, req.UserID, req.CheckoutID)
	if err != nil {
		return nil, mapRepositoryError(err)
	}

	if len(existingOrders) > 0 {
		if err := s.commitCartCheckout(ctx, req.UserID, req.CheckoutID, true); err != nil {
			return nil, err
		}

		return existingOrders, nil
	}

	reservationResp, err := s.cart.Order.ReserveCheckoutItems(ctx, &cartorderpb.ReserveCheckoutItemsRequest{
		UserId:              req.UserID,
		CheckoutId:          req.CheckoutID,
		ProductIds:          req.ProductIDs,
		ExpectedCartVersion: req.ExpectedCartVersion,
	})
	if err != nil {
		return nil, err
	}

	orders := domain.OrderListFromReservation(reservationResp.GetReservation())
	if len(orders) == 0 {
		return nil, status.Error(codes.FailedPrecondition, "checkout reservation is empty")
	}

	createdOrders, err := s.repository.CreateOrders(ctx, orders)
	if err != nil {
		s.releaseCartCheckout(ctx, req.UserID, req.CheckoutID, "order creation failed")

		return nil, mapRepositoryError(err)
	}

	for _, order := range createdOrders {
		err = s.outbox.SendOrderCreate(ctx, order)
		if err != nil {
			return nil, err
		}
	}

	if err := s.commitCartCheckout(ctx, req.UserID, req.CheckoutID, false); err != nil {
		return nil, err
	}

	return createdOrders, nil
}

func (s *service) commitCartCheckout(ctx context.Context, userID int64, checkoutID string, ignoreNotFound bool) error {
	_, err := s.cart.Order.CommitCheckout(ctx, &cartorderpb.CommitCheckoutRequest{
		UserId:     userID,
		CheckoutId: checkoutID,
	})
	if ignoreNotFound && status.Code(err) == codes.NotFound {
		return nil
	}

	return err
}

func (s *service) releaseCartCheckout(ctx context.Context, userID int64, checkoutID string, reason string) {
	_, _ = s.cart.Order.ReleaseCheckout(ctx, &cartorderpb.ReleaseCheckoutRequest{
		UserId:     userID,
		CheckoutId: checkoutID,
		Reason:     reason,
	})
}
