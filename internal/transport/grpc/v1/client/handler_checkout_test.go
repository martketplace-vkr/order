package client

import (
	"context"
	"testing"

	"github.com/martketplace-vkr/order/domain"
	"github.com/martketplace-vkr/order/internal/service/client/dto"
	clientpb "github.com/martketplace-vkr/order/pkg/api/grpc/v1/client"
)

type checkoutServiceStub struct {
	req dto.CheckoutRequest
}

func (s *checkoutServiceStub) Checkout(_ context.Context, req dto.CheckoutRequest) ([]domain.Order, error) {
	s.req = req
	return nil, nil
}

func (checkoutServiceStub) GetOrder(context.Context, int64, int64) (*domain.Order, error) {
	return nil, nil
}

func (checkoutServiceStub) GetOrderList(context.Context, int64) ([]domain.Order, error) {
	return nil, nil
}

func (checkoutServiceStub) HasSuccessfulProductOrder(context.Context, int64, int64) (bool, int64, error) {
	return false, 0, nil
}

func (checkoutServiceStub) CancelOrder(context.Context, int64, int64) (*domain.Order, error) {
	return nil, nil
}

func TestCheckoutDoesNotPanicWithoutPaymentAndDelivery(t *testing.T) {
	t.Parallel()

	svc := &checkoutServiceStub{}
	h := New(svc)

	_, err := h.Checkout(context.Background(), &clientpb.CheckoutRequest{
		UserId:     1,
		CheckoutId: "chk-1",
		ProductIds: []int64{1024},
	})
	if err != nil {
		t.Fatalf("checkout returned error: %v", err)
	}

	if svc.req.Payment.Type != domain.OnlineByCrypto {
		t.Fatalf("unexpected payment type: %v", svc.req.Payment.Type)
	}
	if svc.req.Payment.Status != domain.PendingPaymentStatus {
		t.Fatalf("unexpected payment status: %v", svc.req.Payment.Status)
	}
	if svc.req.Delivery.Type != domain.PickUp {
		t.Fatalf("unexpected delivery type: %v", svc.req.Delivery.Type)
	}
}
