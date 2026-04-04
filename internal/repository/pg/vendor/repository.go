package vendor

import (
	"context"

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

func (r *repository) GetOrder(ctx context.Context, vendorID int64, orderID int64) (*domain.Order, error) {
	query := `
		select
			id,
			user_id,
			vendor_id,
			status,
			product_id,
			product_name,
			product_image_url,
			quantity,
			unit_price,
			total_price,
			created_at,
			updated_at
		from orders
		where vendor_id = $1
			and id = $2
	`

	order := new(domain.Order)

	err := r.ctxGetter.DefaultTrOrDB(ctx, r.db).GetContext(
		ctx,
		order,
		query,
		vendorID,
		orderID,
	)
	if err != nil {
		return nil, err
	}

	return order, nil
}

func (r *repository) GetOrderList(ctx context.Context, vendorID int64) ([]domain.Order, error) {
	query := `
		select
			id,
			user_id,
			vendor_id,
			status,
			product_id,
			product_name,
			product_image_url,
			quantity,
			unit_price,
			total_price,
			created_at,
			updated_at
		from orders
		where vendor_id = $1
		order by created_at desc, id desc
	`

	orders := make([]domain.Order, 0)

	err := r.ctxGetter.DefaultTrOrDB(ctx, r.db).SelectContext(
		ctx,
		&orders,
		query,
		vendorID,
	)
	if err != nil {
		return nil, err
	}

	return orders, nil
}

func (r *repository) UpdateOrder(
	ctx context.Context,
	vendorID int64,
	orderID int64,
	status string,
) (*domain.Order, error) {
	query := `
		update orders
		set
			status = $3,
			updated_at = now()
		where vendor_id = $1
			and id = $2
		returning
			id,
			user_id,
			vendor_id,
			status,
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

	err := r.ctxGetter.DefaultTrOrDB(ctx, r.db).GetContext(
		ctx,
		order,
		query,
		vendorID,
		orderID,
		status,
	)
	if err != nil {
		return nil, err
	}

	return order, nil
}
