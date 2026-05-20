package admin

import (
	"context"
	"strings"

	trmsqlx "github.com/avito-tech/go-transaction-manager/sqlx"
	"github.com/jmoiron/sqlx"
	"github.com/martketplace-vkr/order/domain"
)

type repository struct {
	ctxGetter *trmsqlx.CtxGetter
	db        *sqlx.DB
}

func New(db *sqlx.DB, ctxGetter *trmsqlx.CtxGetter) *repository {
	return &repository{
		db:        db,
		ctxGetter: ctxGetter,
	}
}

type ListOrdersFilter struct {
	PaymentStatus     string
	FulfillmentStatus string
	Limit             uint32
	Offset            uint64
}

func (r *repository) ListOrders(ctx context.Context, filter ListOrdersFilter) ([]domain.Order, error) {
	query := `
		select
			id,
			checkout_id,
			user_id,
			vendor_id,
			status,
			payment_status,
			fulfillment_status,
			delivery_address_id,
			product_id,
			product_name,
			product_image_url,
			quantity,
			unit_price,
			total_price,
			created_at,
			updated_at
		from "order"."order"
		where ($1 = '' or payment_status = $1)
			and ($2 = '' or fulfillment_status = $2)
		order by created_at desc, id desc
		limit $3 offset $4
	`

	limit := filter.Limit
	if limit == 0 {
		limit = 100
	}

	orders := make([]domain.Order, 0)
	err := r.ctxGetter.DefaultTrOrDB(ctx, r.db).SelectContext(
		ctx,
		&orders,
		query,
		strings.TrimSpace(filter.PaymentStatus),
		strings.TrimSpace(filter.FulfillmentStatus),
		limit,
		filter.Offset,
	)
	if err != nil {
		return nil, err
	}

	return orders, nil
}

func (r *repository) UpdatePaymentStatus(ctx context.Context, orderID int64, paymentStatus domain.OrderPaymentStatus) (*domain.Order, error) {
	query := `
		update "order"."order"
		set
			payment_status = $2,
			updated_at = now()
		where id = $1
		returning
			id,
			checkout_id,
			user_id,
			vendor_id,
			status,
			payment_status,
			fulfillment_status,
			delivery_address_id,
			product_id,
			product_name,
			product_image_url,
			quantity,
			unit_price,
			total_price,
			created_at,
			updated_at
	`

	order := new(domain.Order)
	err := r.ctxGetter.DefaultTrOrDB(ctx, r.db).GetContext(ctx, order, query, orderID, paymentStatus)
	if err != nil {
		return nil, err
	}

	return order, nil
}
