package mapper

import (
	orderdomain "github.com/martketplace-vkr/order/domain"
	domainpb "github.com/martketplace-vkr/order/pkg/api/grpc/v1/domain"
)

func PaymentFromProto(payment *domainpb.Payment) orderdomain.Payment {
	result := orderdomain.Payment{
		Type:   orderdomain.OnlineByCrypto,
		Status: orderdomain.PendingPaymentStatus,
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
		Type: orderdomain.PickUp,
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
