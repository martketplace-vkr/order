package client

import (
	"context"
	"errors"

	"github.com/martketplace-vkr/order/internal/service/ordererrors"
	"github.com/martketplace-vkr/order/internal/transport/grpc/v1/mapper"
	clientpb "github.com/martketplace-vkr/order/pkg/api/grpc/v1/client"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Handler struct {
	service service
	clientpb.UnimplementedOrderClientServiceServer
}

func New(service service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) Checkout(ctx context.Context, req *clientpb.CheckoutRequest) (resp *clientpb.CheckoutResponse, err error) {
	orders, err := h.service.Checkout(ctx, req.UserId, req.CheckoutId, req.ProductIds, req.ExpectedCartVersion)
	if err != nil {
		return resp, toStatusError(err)
	}

	return &clientpb.CheckoutResponse{
		Orders: mapper.OrdersToProto(orders),
	}, nil
}

func (h *Handler) GetOrder(ctx context.Context, req *clientpb.GetOrderRequest) (resp *clientpb.GetOrderResponse, err error) {
	order, err := h.service.GetOrder(ctx, req.UserId, req.OrderId)
	if err != nil {
		return resp, toStatusError(err)
	}

	return &clientpb.GetOrderResponse{
		Order: mapper.OrderToProto(order),
	}, nil
}

func (h *Handler) GetOrderList(ctx context.Context, req *clientpb.GetOrderListRequest) (resp *clientpb.GetOrderListResponse, err error) {
	orders, err := h.service.GetOrderList(ctx, req.UserId)
	if err != nil {
		return resp, toStatusError(err)
	}

	return &clientpb.GetOrderListResponse{
		Orders: mapper.OrdersToProto(orders),
	}, nil
}

func (h *Handler) CancellOrder(ctx context.Context, req *clientpb.CancelOrderRequest) (resp *clientpb.CancelOrderResponse, err error) {
	order, err := h.service.CancelOrder(ctx, req.UserId, req.OrderId)
	if err != nil {
		return resp, toStatusError(err)
	}

	return &clientpb.CancelOrderResponse{
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
