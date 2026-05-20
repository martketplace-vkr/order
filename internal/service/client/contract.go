package client

import (
	"context"

	cartorderpb "github.com/martketplace-vkr/cart/pkg/api/grpc/v1/order"
	"github.com/martketplace-vkr/order/domain"
	"google.golang.org/grpc"
)

type (
	repository interface {
		CreateOrders(ctx context.Context, orders []domain.Order) ([]domain.Order, error)
		GetOrdersByCheckout(ctx context.Context, userID int64, checkoutID string) ([]domain.Order, error)
		GetOrder(ctx context.Context, userID int64, orderID int64) (*domain.Order, error)
		GetOrderList(ctx context.Context, userID int64) ([]domain.Order, error)
		HasSuccessfulProductOrder(ctx context.Context, userID int64, productID int64) (bool, int64, error)
		CancelOrder(ctx context.Context, userID int64, orderID int64) (*domain.Order, error)
	}
	cartClient interface {
		ReserveCheckoutItems(ctx context.Context, in *cartorderpb.ReserveCheckoutItemsRequest, opts ...grpc.CallOption) (*cartorderpb.ReserveCheckoutItemsResponse, error)
		CommitCheckout(ctx context.Context, in *cartorderpb.CommitCheckoutRequest, opts ...grpc.CallOption) (*cartorderpb.CommitCheckoutResponse, error)
		ReleaseCheckout(ctx context.Context, in *cartorderpb.ReleaseCheckoutRequest, opts ...grpc.CallOption) (*cartorderpb.ReleaseCheckoutResponse, error)
	}
	outbox interface {
		SendOrderCreate(ctx context.Context, order domain.Order) (err error)
		SendOrderCancelled(ctx context.Context, order domain.Order) error
	}
)
