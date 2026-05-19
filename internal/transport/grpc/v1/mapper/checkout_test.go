package mapper

import (
	"testing"

	orderdomain "github.com/martketplace-vkr/order/domain"
	domainpb "github.com/martketplace-vkr/order/pkg/api/grpc/v1/domain"
)

func TestPaymentFromProtoDefaultsWhenNil(t *testing.T) {
	t.Parallel()

	payment := PaymentFromProto(nil)
	if payment.Type != orderdomain.OnlineByCrypto {
		t.Fatalf("unexpected payment type: %v", payment.Type)
	}
	if payment.Status != orderdomain.PendingPaymentStatus {
		t.Fatalf("unexpected payment status: %v", payment.Status)
	}
}

func TestDeliveryFromProtoDefaultsWhenNil(t *testing.T) {
	t.Parallel()

	delivery := DeliveryFromProto(nil)
	if delivery.Type != orderdomain.PickUp {
		t.Fatalf("unexpected delivery type: %v", delivery.Type)
	}
	if delivery.EntityID() != 0 {
		t.Fatalf("expected empty entity id, got %d", delivery.EntityID())
	}
}

func TestDeliveryFromProtoUsesProvidedIDs(t *testing.T) {
	t.Parallel()

	pickUpPointID := int64(42)
	delivery := DeliveryFromProto(&domainpb.Delivery{
		Type:          int64(orderdomain.PickUp),
		PickUpPointId: &pickUpPointID,
	})
	if delivery.EntityID() != pickUpPointID {
		t.Fatalf("unexpected entity id: %d", delivery.EntityID())
	}
}
