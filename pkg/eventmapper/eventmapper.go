package eventmapper

import (
	"context"

	"github.com/martketplace-vkr/pkg/inbox/dto"
)

type service interface {
	HandleOrderPaid(context.Context, dto.Event) error
	HandlePaymentCaptured(context.Context, dto.Event) error
	HandlePaymentReleased(context.Context, dto.Event) error
	HandlePaymentExpired(context.Context, dto.Event) error
}

func GetEventMapper(svc service) map[string]func(context.Context, dto.Event) error {
	return map[string]func(context.Context, dto.Event) error{
		"order_paid":       svc.HandleOrderPaid,
		"payment_captured": svc.HandlePaymentCaptured,
		"payment_released": svc.HandlePaymentReleased,
		"payment_expired":  svc.HandlePaymentExpired,
	}
}
