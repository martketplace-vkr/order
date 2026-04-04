package client

import (
	"context"
	"errors"
	"testing"
	"time"

	cartorderpb "github.com/martketplace-vkr/cart/pkg/api/grpc/v1/order"
	"github.com/martketplace-vkr/order/domain"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type checkoutRepositoryStub struct {
	existingOrders []domain.Order
	createdOrders  []domain.Order
	createErr      error
	createCalls    int
}

func (r *checkoutRepositoryStub) CreateOrders(_ context.Context, orders []domain.Order) ([]domain.Order, error) {
	r.createCalls++
	if r.createErr != nil {
		return nil, r.createErr
	}
	if r.createdOrders != nil {
		return r.createdOrders, nil
	}

	return orders, nil
}

func (r *checkoutRepositoryStub) GetOrdersByCheckout(_ context.Context, _ int64, _ string) ([]domain.Order, error) {
	return r.existingOrders, nil
}

func (r *checkoutRepositoryStub) GetOrder(_ context.Context, _ int64, _ int64) (*domain.Order, error) {
	return nil, nil
}

func (r *checkoutRepositoryStub) GetOrderList(_ context.Context, _ int64) ([]domain.Order, error) {
	return nil, nil
}

func (r *checkoutRepositoryStub) CancelOrder(_ context.Context, _ int64, _ int64) (*domain.Order, error) {
	return nil, nil
}

type cartClientStub struct {
	reserveResp  *cartorderpb.ReserveCheckoutItemsResponse
	reserveErr   error
	commitErr    error
	releaseErr   error
	reserveCalls int
	commitCalls  int
	releaseCalls int
	lastRelease  *cartorderpb.ReleaseCheckoutRequest
	lastReserve  *cartorderpb.ReserveCheckoutItemsRequest
	lastCommit   *cartorderpb.CommitCheckoutRequest
}

func (c *cartClientStub) ReserveCheckoutItems(ctx context.Context, in *cartorderpb.ReserveCheckoutItemsRequest, _ ...grpc.CallOption) (*cartorderpb.ReserveCheckoutItemsResponse, error) {
	c.reserveCalls++
	c.lastReserve = in

	return c.reserveResp, c.reserveErr
}

func (c *cartClientStub) CommitCheckout(ctx context.Context, in *cartorderpb.CommitCheckoutRequest, _ ...grpc.CallOption) (*cartorderpb.CommitCheckoutResponse, error) {
	c.commitCalls++
	c.lastCommit = in

	return &cartorderpb.CommitCheckoutResponse{}, c.commitErr
}

func (c *cartClientStub) ReleaseCheckout(ctx context.Context, in *cartorderpb.ReleaseCheckoutRequest, _ ...grpc.CallOption) (*cartorderpb.ReleaseCheckoutResponse, error) {
	c.releaseCalls++
	c.lastRelease = in

	return &cartorderpb.ReleaseCheckoutResponse{}, c.releaseErr
}

func TestCheckoutCreatesOrdersAndCommitsCart(t *testing.T) {
	repo := &checkoutRepositoryStub{}
	cart := &cartClientStub{
		reserveResp: &cartorderpb.ReserveCheckoutItemsResponse{
			Reservation: &cartorderpb.CheckoutReservation{
				CheckoutId: "chk-1",
				UserId:     42,
				Items: []*cartorderpb.CheckoutCartItem{
					{
						ProductId:   1001,
						VendorId:    7,
						ProductName: "Phone",
						ImageUrl:    "image.png",
						Quantity:    2,
						UnitPrice:   "10",
						TotalPrice:  "20",
					},
				},
				CreatedAt: timestamppb.Now(),
			},
		},
	}

	svc := New(repo, cart, time.Second)

	orders, err := svc.Checkout(context.Background(), 42, "chk-1", []int64{1001}, 11)
	if err != nil {
		t.Fatalf("Checkout returned error: %v", err)
	}

	if len(orders) != 1 {
		t.Fatalf("unexpected orders count: %d", len(orders))
	}
	if orders[0].CheckoutID != "chk-1" || orders[0].ProductName != "Phone" {
		t.Fatalf("unexpected order: %+v", orders[0])
	}
	if repo.createCalls != 1 {
		t.Fatalf("expected one create call, got %d", repo.createCalls)
	}
	if cart.reserveCalls != 1 || cart.commitCalls != 1 {
		t.Fatalf("unexpected cart calls reserve=%d commit=%d", cart.reserveCalls, cart.commitCalls)
	}
	if cart.lastReserve.GetExpectedCartVersion() != 11 {
		t.Fatalf("unexpected expected_cart_version: %d", cart.lastReserve.GetExpectedCartVersion())
	}
}

func TestCheckoutReturnsExistingOrdersForRetry(t *testing.T) {
	repo := &checkoutRepositoryStub{
		existingOrders: []domain.Order{
			{
				ID:         1,
				CheckoutID: "chk-1",
				UserID:     42,
				ProductID:  1001,
			},
		},
	}
	cart := &cartClientStub{
		commitErr: status.Error(codes.NotFound, "already committed"),
	}

	svc := New(repo, cart, time.Second)

	orders, err := svc.Checkout(context.Background(), 42, "chk-1", []int64{1001}, 0)
	if err != nil {
		t.Fatalf("Checkout returned error: %v", err)
	}

	if len(orders) != 1 || orders[0].ID != 1 {
		t.Fatalf("unexpected orders: %+v", orders)
	}
	if cart.reserveCalls != 0 {
		t.Fatalf("expected no reserve call, got %d", cart.reserveCalls)
	}
	if cart.commitCalls != 1 {
		t.Fatalf("expected one commit call, got %d", cart.commitCalls)
	}
}

func TestCheckoutReleasesReservationOnCreateError(t *testing.T) {
	repo := &checkoutRepositoryStub{
		createErr: errors.New("insert failed"),
	}
	cart := &cartClientStub{
		reserveResp: &cartorderpb.ReserveCheckoutItemsResponse{
			Reservation: &cartorderpb.CheckoutReservation{
				CheckoutId: "chk-1",
				UserId:     42,
				Items: []*cartorderpb.CheckoutCartItem{
					{
						ProductId:   1001,
						VendorId:    7,
						ProductName: "Phone",
						Quantity:    1,
						UnitPrice:   "10",
						TotalPrice:  "10",
					},
				},
			},
		},
	}

	svc := New(repo, cart, time.Second)

	_, err := svc.Checkout(context.Background(), 42, "chk-1", []int64{1001}, 0)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if cart.releaseCalls != 1 {
		t.Fatalf("expected one release call, got %d", cart.releaseCalls)
	}
	if cart.lastRelease.GetCheckoutId() != "chk-1" {
		t.Fatalf("unexpected release request: %+v", cart.lastRelease)
	}
}
