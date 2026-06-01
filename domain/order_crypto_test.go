package domain

import (
	"testing"

	cartorderpb "github.com/martketplace-vkr/cart/pkg/api/grpc/v1/order"
	"github.com/martketplace-vkr/pkg/utils/currency"
)

func TestOrderListFromReservationPreservesItemCurrency(t *testing.T) {
	orders := OrderListFromReservation(&cartorderpb.CheckoutReservation{
		CheckoutId: "checkout-1",
		UserId:     7,
		Items: []*cartorderpb.CheckoutCartItem{
			{ProductId: 1, CurrencyId: int64(currency.USDTinTRC), UnitPrice: "10", TotalPrice: "10"},
			{ProductId: 2, CurrencyId: int64(currency.RUB), UnitPrice: "900", TotalPrice: "900"},
		},
	}).WithPaymentAndDelivery(Payment{}, Delivery{})

	if len(orders) != 2 {
		t.Fatalf("unexpected order count: %d", len(orders))
	}
	if orders[0].Payment.CurrencyID != int64(currency.USDTinTRC) || orders[0].Payment.Type != OnlineByCrypto {
		t.Fatalf("unexpected USDT payment: %+v", orders[0].Payment)
	}
	if orders[1].Payment.CurrencyID != int64(currency.RUB) || orders[1].Payment.Type != OnlineByCard {
		t.Fatalf("unexpected RUB payment: %+v", orders[1].Payment)
	}
}
