package client

import (
	"context"
	"errors"
	"fmt"

	"github.com/martketplace-vkr/order/internal/service/client/dto"
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
	if err := validateID("user_id", req.GetUserId()); err != nil {
		return nil, err
	}

	if req.GetCheckoutId() == "" {
		return nil, status.Error(codes.InvalidArgument, "checkout_id is required")
	}

	if len(req.GetProductIds()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "product_ids are required")
	}

	for _, productID := range req.GetProductIds() {
		if productID <= 0 {
			return nil, status.Error(codes.InvalidArgument, "product_ids must be greater than zero")
		}
	}

	orders, err := h.service.Checkout(ctx, dto.CheckoutRequest{
		UserID:              req.GetUserId(),
		CheckoutID:          req.GetCheckoutId(),
		ProductIDs:          req.GetProductIds(),
		ExpectedCartVersion: req.GetExpectedCartVersion(),
		PreferredCurrencyID: req.GetPreferredCurrencyId(),
		Payment:             mapper.PaymentFromProto(req.GetPayment()),
		Delivery:            mapper.DeliveryFromProto(req.GetDelivery()),
	})
	if err != nil {
		return resp, toStatusError(err)
	}

	return &clientpb.CheckoutResponse{
		Orders: mapper.OrdersToProto(orders),
	}, nil
}

func (h *Handler) GetOrder(ctx context.Context, req *clientpb.GetOrderRequest) (resp *clientpb.GetOrderResponse, err error) {
	if err := validateID("user_id", req.GetUserId()); err != nil {
		return nil, err
	}
	if err := validateID("order_id", req.GetOrderId()); err != nil {
		return nil, err
	}

	order, err := h.service.GetOrder(ctx, req.UserId, req.OrderId)
	if err != nil {
		return resp, toStatusError(err)
	}

	return &clientpb.GetOrderResponse{
		Order: mapper.OrderToProto(order),
	}, nil
}

func (h *Handler) GetOrderList(ctx context.Context, req *clientpb.GetOrderListRequest) (resp *clientpb.GetOrderListResponse, err error) {
	if err := validateID("user_id", req.GetUserId()); err != nil {
		return nil, err
	}

	orders, err := h.service.GetOrderList(ctx, req.UserId)
	if err != nil {
		return resp, toStatusError(err)
	}

	return &clientpb.GetOrderListResponse{
		Orders: mapper.OrdersToProto(orders),
	}, nil
}

func (h *Handler) HasSuccessfulProductOrder(ctx context.Context, req *clientpb.HasSuccessfulProductOrderRequest) (resp *clientpb.HasSuccessfulProductOrderResponse, err error) {
	if err := validateID("user_id", req.GetUserId()); err != nil {
		return nil, err
	}
	if err := validateID("product_id", req.GetProductId()); err != nil {
		return nil, err
	}

	hasOrder, vendorID, err := h.service.HasSuccessfulProductOrder(ctx, req.GetUserId(), req.GetProductId())
	if err != nil {
		return resp, toStatusError(err)
	}

	return &clientpb.HasSuccessfulProductOrderResponse{
		HasOrder: hasOrder,
		VendorId: vendorID,
	}, nil
}

func (h *Handler) CancellOrder(ctx context.Context, req *clientpb.CancelOrderRequest) (resp *clientpb.CancelOrderResponse, err error) {
	if err := validateID("user_id", req.GetUserId()); err != nil {
		return nil, err
	}
	if err := validateID("order_id", req.GetOrderId()); err != nil {
		return nil, err
	}

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

func validateID(field string, value int64) error {
	if value <= 0 {
		return fmt.Errorf("%w: %s must be greater than zero", ordererrors.ErrInvalidArgument, field)
	}

	return nil
}
