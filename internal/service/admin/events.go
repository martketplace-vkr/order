package admin

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"

	"github.com/martketplace-vkr/order/domain"
	"github.com/martketplace-vkr/pkg/inbox/dto"
)

type paymentEvent struct {
	OrderID int64  `json:"order_id"`
	UserID  int64  `json:"user_id"`
	Status  string `json:"status"`
	Reason  string `json:"reason"`
}

func (s *service) HandleOrderPaid(ctx context.Context, event dto.Event) error {
	return s.updatePaymentStatusFromEvent(ctx, event, domain.OrderPaymentReserved)
}

func (s *service) HandlePaymentCaptured(ctx context.Context, event dto.Event) error {
	return s.updatePaymentStatusFromEvent(ctx, event, domain.OrderPaymentCaptured)
}

func (s *service) HandlePaymentReleased(ctx context.Context, event dto.Event) error {
	return s.updatePaymentStatusFromEvent(ctx, event, domain.OrderPaymentReleased)
}

func (s *service) HandlePaymentExpired(ctx context.Context, event dto.Event) error {
	return s.updatePaymentStatusFromEvent(ctx, event, domain.OrderPaymentExpired)
}

func (s *service) updatePaymentStatusFromEvent(ctx context.Context, event dto.Event, status domain.OrderPaymentStatus) error {
	var payload paymentEvent
	if err := json.Unmarshal(event.Payload, &payload); err != nil {
		return err
	}
	if payload.OrderID <= 0 {
		return nil
	}

	_, err := s.repository.UpdatePaymentStatus(ctx, payload.OrderID, status)
	if errors.Is(err, sql.ErrNoRows) {
		return nil
	}
	return err
}
