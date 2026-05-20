package mapper

import (
	orderdomain "github.com/martketplace-vkr/order/domain"
	domainpb "github.com/martketplace-vkr/order/pkg/api/grpc/v1/domain"
	"github.com/martketplace-vkr/pkg/utils/currency"
)

func PaymentFromProto(payment *domainpb.Payment) orderdomain.Payment {
	result := orderdomain.Payment{
		CurrencyID: int64(currency.RUB),
		Type:       orderdomain.OnlineByCard,
		Status:     orderdomain.PendingPaymentStatus,
	}
	if payment == nil {
		return result
	}

	if paymentType := orderdomain.PaymentType(payment.GetType()); paymentType > 0 {
		result.Type = paymentType
	}

	return result
}

func DeliveryFromProto(delivery *domainpb.Delivery) orderdomain.Delivery {
	result := orderdomain.Delivery{
		Type: orderdomain.ClientDelivery,
	}
	if delivery == nil {
		return result
	}

	if deliveryType := orderdomain.DeliveryType(delivery.GetType()); deliveryType > 0 {
		result.Type = deliveryType
	}

	if delivery.PickUpPointId != nil {
		pickUpPointID := delivery.GetPickUpPointId()
		result.PickUpPointID = &pickUpPointID
	}

	if delivery.ClintAddressId != nil {
		clientAddressID := delivery.GetClintAddressId()
		result.ClientAddress = &clientAddressID
	}

	return result
}
