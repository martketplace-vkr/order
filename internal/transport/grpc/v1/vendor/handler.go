package vendor

import (
	"context"
	"errors"

	"github.com/martketplace-vkr/order/internal/service/ordererrors"
	"github.com/martketplace-vkr/order/internal/transport/grpc/v1/mapper"
	vendorpb "github.com/martketplace-vkr/order/pkg/api/grpc/v1/vendor"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Handler struct {
	service service
	vendorpb.UnimplementedOrderVendorServiceServer
}

func New(service service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) GetOrder(ctx context.Context, req *vendorpb.GetOrderRequest) (resp *vendorpb.GetOrderResponse, err error) {
	order, err := h.service.GetOrder(ctx, req.VendorId, req.OrderId)
	if err != nil {
		return resp, toStatusError(err)
	}

	return &vendorpb.GetOrderResponse{
		Order: mapper.OrderToProto(order),
	}, nil
}

func (h *Handler) GetOrderList(ctx context.Context, req *vendorpb.GetOrderListRequest) (resp *vendorpb.GetOrderListResponse, err error) {
	orders, err := h.service.GetOrderList(ctx, req.VendorId)
	if err != nil {
		return resp, toStatusError(err)
	}

	return &vendorpb.GetOrderListResponse{
		Orders: mapper.OrdersToProto(orders),
	}, nil
}

func (h *Handler) UpdateOrder(ctx context.Context, req *vendorpb.UpdateOrderRequest) (resp *vendorpb.UpdateOrderResponse, err error) {
	order, err := h.service.UpdateOrder(ctx, req.VendorId, req.OrderId, req.Status)
	if err != nil {
		return resp, toStatusError(err)
	}

	return &vendorpb.UpdateOrderResponse{
		Order: mapper.OrderToProto(order),
	}, nil
}

func toStatusError(err error) error {
	switch {
	case errors.Is(err, ordererrors.ErrInvalidArgument):
		return status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, ordererrors.ErrOrderNotFound):
		return status.Error(codes.NotFound, err.Error())
	default:
		return err
	}
}
