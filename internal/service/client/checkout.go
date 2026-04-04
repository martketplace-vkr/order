package client

import (
	"context"
	"time"

	cartorderpb "github.com/martketplace-vkr/cart/pkg/api/grpc/v1/order"
	"github.com/martketplace-vkr/order/domain"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const orderStatusCreated = "created"

func (s *service) Checkout(
	ctx context.Context,
	userID int64,
	checkoutID string,
	productIDs []int64,
	expectedCartVersion uint64,
) ([]domain.Order, error) {
	if err := validateID("user_id", userID); err != nil {
		return nil, err
	}
	if checkoutID == "" {
		return nil, status.Error(codes.InvalidArgument, "checkout_id is required")
	}
	if len(productIDs) == 0 {
		return nil, status.Error(codes.InvalidArgument, "product_ids are required")
	}
	for _, productID := range productIDs {
		if productID <= 0 {
			return nil, status.Error(codes.InvalidArgument, "product_ids must be greater than zero")
		}
	}

	existingOrders, err := s.repository.GetOrdersByCheckout(ctx, userID, checkoutID)
	if err != nil {
		return nil, mapRepositoryError(err)
	}
	if len(existingOrders) > 0 {
		if err := s.commitCartCheckout(ctx, userID, checkoutID, true); err != nil {
			return nil, err
		}

		return existingOrders, nil
	}

	cartCtx, cancel := s.newCartContext(ctx)
	defer cancel()

	reservationResp, err := s.cartClient.ReserveCheckoutItems(cartCtx, &cartorderpb.ReserveCheckoutItemsRequest{
		UserId:              userID,
		CheckoutId:          checkoutID,
		ProductIds:          productIDs,
		ExpectedCartVersion: expectedCartVersion,
	})
	if err != nil {
		return nil, err
	}

	orders := buildOrdersFromReservation(reservationResp.GetReservation())
	if len(orders) == 0 {
		return nil, status.Error(codes.FailedPrecondition, "checkout reservation is empty")
	}

	createdOrders, err := s.repository.CreateOrders(ctx, orders)
	if err != nil {
		s.releaseCartCheckout(ctx, userID, checkoutID, "order creation failed")

		return nil, mapRepositoryError(err)
	}

	if err := s.commitCartCheckout(ctx, userID, checkoutID, false); err != nil {
		return nil, err
	}

	return createdOrders, nil
}

func (s *service) commitCartCheckout(ctx context.Context, userID int64, checkoutID string, ignoreNotFound bool) error {
	cartCtx, cancel := s.newCartContext(ctx)
	defer cancel()

	_, err := s.cartClient.CommitCheckout(cartCtx, &cartorderpb.CommitCheckoutRequest{
		UserId:     userID,
		CheckoutId: checkoutID,
	})
	if ignoreNotFound && status.Code(err) == codes.NotFound {
		return nil
	}

	return err
}

func (s *service) releaseCartCheckout(ctx context.Context, userID int64, checkoutID string, reason string) {
	cartCtx, cancel := s.newCartContext(ctx)
	defer cancel()

	_, _ = s.cartClient.ReleaseCheckout(cartCtx, &cartorderpb.ReleaseCheckoutRequest{
		UserId:     userID,
		CheckoutId: checkoutID,
		Reason:     reason,
	})
}

func (s *service) newCartContext(ctx context.Context) (context.Context, context.CancelFunc) {
	if s.cartTimeout <= 0 {
		return context.WithCancel(ctx)
	}

	return context.WithTimeout(ctx, s.cartTimeout)
}

func buildOrdersFromReservation(reservation *cartorderpb.CheckoutReservation) []domain.Order {
	if reservation == nil {
		return nil
	}

	orders := make([]domain.Order, 0, len(reservation.GetItems()))
	now := time.Now().UTC()

	for _, item := range reservation.GetItems() {
		if item == nil {
			continue
		}

		orders = append(orders, domain.Order{
			CheckoutID:      reservation.GetCheckoutId(),
			UserID:          reservation.GetUserId(),
			VendorID:        item.GetVendorId(),
			Status:          orderStatusCreated,
			ProductID:       item.GetProductId(),
			ProductName:     item.GetProductName(),
			ProductImageURL: item.GetImageUrl(),
			Quantity:        int64(item.GetQuantity()),
			UnitPrice:       item.GetUnitPrice(),
			TotalPrice:      item.GetTotalPrice(),
			CreatedAt:       now,
			UpdatedAt:       now,
		})
	}

	return orders
}
