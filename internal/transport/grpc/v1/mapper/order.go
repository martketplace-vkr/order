package mapper

import (
	"time"

	orderdomain "github.com/martketplace-vkr/order/domain"
	domainpb "github.com/martketplace-vkr/order/pkg/api/grpc/v1/domain"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func OrderToProto(order *orderdomain.Order) *domainpb.Order {
	if order == nil {
		return nil
	}

	fulfillmentStatus := order.FulfillmentStatus
	if fulfillmentStatus == 0 {
		fulfillmentStatus = order.Status
	}

	return &domainpb.Order{
		Id:                order.ID,
		UserId:            order.UserID,
		VendorId:          order.VendorID,
		Status:            fulfillmentStatus.String(),
		PaymentStatus:     string(order.PaymentStatus),
		FulfillmentStatus: fulfillmentStatus.String(),
		DeliveryAddressId: order.DeliveryAddressID,
		CurrencyId:        order.CurrencyID,
		Product: &domainpb.OrderProduct{
			ProductId:   order.ProductID,
			ProductName: order.ProductName,
			ImageUrl:    order.ProductImageURL,
		},
		Quantity:   uint32(order.Quantity),
		UnitPrice:  order.UnitPrice,
		TotalPrice: order.TotalPrice,
		Comment:    order.Comment,
		CreatedAt:  newTimestamp(order.CreatedAt),
		UpdatedAt:  newTimestamp(order.UpdatedAt),
	}
}

func OrdersToProto(orders []orderdomain.Order) []*domainpb.Order {
	resp := make([]*domainpb.Order, 0, len(orders))
	for i := range orders {
		resp = append(resp, OrderToProto(&orders[i]))
	}

	return resp
}

func newTimestamp(value time.Time) *timestamppb.Timestamp {
	if value.IsZero() {
		return nil
	}

	return timestamppb.New(value)
}
