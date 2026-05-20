package admin

import (
	"context"
	"errors"

	"github.com/martketplace-vkr/order/internal/service/ordererrors"
	"github.com/martketplace-vkr/order/internal/transport/grpc/v1/mapper"
	adminpb "github.com/martketplace-vkr/order/pkg/api/grpc/v1/admin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Handler struct {
	service service
	adminpb.UnimplementedOrderAdminServiceServer
}

func New(service service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) ListOrders(ctx context.Context, req *adminpb.ListOrdersRequest) (*adminpb.ListOrdersResponse, error) {
	orders, err := h.service.ListOrders(ctx, req.GetPaymentStatus(), req.GetFulfillmentStatus(), req.GetLimit(), req.GetOffset())
	if err != nil {
		return nil, mapError(err)
	}

	return &adminpb.ListOrdersResponse{Orders: mapper.OrdersToProto(orders)}, nil
}

func (h *Handler) UpdatePaymentStatus(ctx context.Context, req *adminpb.UpdatePaymentStatusRequest) (*adminpb.UpdatePaymentStatusResponse, error) {
	order, err := h.service.UpdatePaymentStatus(ctx, req.GetOrderId(), req.GetPaymentStatus())
	if err != nil {
		return nil, mapError(err)
	}

	return &adminpb.UpdatePaymentStatusResponse{Order: mapper.OrderToProto(order)}, nil
}

func mapError(err error) error {
	switch {
	case errors.Is(err, ordererrors.ErrInvalidArgument):
		return status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, ordererrors.ErrOrderNotFound):
		return status.Error(codes.NotFound, err.Error())
	default:
		return status.Error(codes.Internal, err.Error())
	}
}
