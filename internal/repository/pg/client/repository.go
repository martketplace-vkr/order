package client

import (
	"context"
	"time"

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

func (r *repository) GetOrder(ctx context.Context, userID int64, orderID int64) (*domain.Order, error) {
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
			coalesce((select p.currency_id from "order".payment p where p.order_id = "order"."order".id), 1000) as currency_id,
			created_at,
			updated_at
		from "order"."order"
		where user_id = $1
			and id = $2
	`

	order := new(domain.Order)

	err := r.ctxGetter.DefaultTrOrDB(ctx, r.db).GetContext(
		ctx,
		order,
		query,
		userID,
		orderID,
	)
	if err != nil {
		return nil, err
	}

	return order, nil
}

func (r *repository) GetOrderList(ctx context.Context, userID int64) ([]domain.Order, error) {
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
			coalesce((select p.currency_id from "order".payment p where p.order_id = "order"."order".id), 1000) as currency_id,
			created_at,
			updated_at
		from "order"."order"
		where user_id = $1
		order by created_at desc, id desc
	`

	orders := make([]domain.Order, 0)

	err := r.ctxGetter.DefaultTrOrDB(ctx, r.db).SelectContext(
		ctx,
		&orders,
		query,
		userID,
	)
	if err != nil {
		return nil, err
	}

	return orders, nil
}

func (r *repository) HasSuccessfulProductOrder(ctx context.Context, userID int64, productID int64) (bool, int64, error) {
	query := `
		select vendor_id
		from "order"."order"
		where user_id = $1
			and product_id = $2
			and fulfillment_status = $3
		order by created_at desc, id desc
		limit 1
	`

	var vendorID int64
	err := r.ctxGetter.DefaultTrOrDB(ctx, r.db).GetContext(
		ctx,
		&vendorID,
		query,
		userID,
		productID,
		domain.Success,
	)
	if err != nil {
		return false, 0, err
	}

	return true, vendorID, nil
}

func (r *repository) CancelOrder(ctx context.Context, userID int64, orderID int64) (*domain.Order, error) {
	query := `
		update "order"."order"
		set
			status = $3,
			fulfillment_status = $3,
			updated_at = now()
		where user_id = $1
			and id = $2
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
			coalesce((select p.currency_id from "order".payment p where p.order_id = "order"."order".id), 1000) as currency_id,
			created_at,
			updated_at
	`

	order := new(domain.Order)

	err := r.ctxGetter.DefaultTrOrDB(ctx, r.db).GetContext(
		ctx,
		order,
		query,
		userID,
		orderID,
		domain.CancelledByClient,
	)
	if err != nil {
		return nil, err
	}

	return order, nil
}

func (r *repository) GetOrdersByCheckout(ctx context.Context, userID int64, checkoutID string) ([]domain.Order, error) {
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
			coalesce((select p.currency_id from "order".payment p where p.order_id = "order"."order".id), 1000) as currency_id,
			created_at,
			updated_at
		from "order"."order"
		where user_id = $1
			and checkout_id = $2
		order by id asc
	`

	orders := make([]domain.Order, 0)

	err := r.ctxGetter.DefaultTrOrDB(ctx, r.db).SelectContext(
		ctx,
		&orders,
		query,
		userID,
		checkoutID,
	)
	if err != nil {
		return nil, err
	}

	return orders, nil
}

func (r *repository) CreateOrders(ctx context.Context, orders []domain.Order) ([]domain.Order, error) {
	if len(orders) == 0 {
		return nil, nil
	}

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	query := `
		insert into "order"."order" (
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
		)
		values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
		on conflict (user_id, checkout_id, product_id)
		do update set
			checkout_id = excluded.checkout_id
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

	paymentQuery := `
		insert into "order".payment(
			order_id,
			currency_id,
			type,
			status
		) values (
			$1,
			$2,
			$3,
			$4
		)
	`

	deliveryQuery := `
		insert into "order".delivery(
			order_id,
			entity_id,
			type,
			status
		) values (
			$1,
			$2,
			$3,
			$4 
		)
	`

	createdOrders := make([]domain.Order, 0, len(orders))
	for i := range orders {
		order := orders[i]
		now := time.Now().UTC()

		if order.CreatedAt.IsZero() {
			order.CreatedAt = now
		}

		if order.UpdatedAt.IsZero() {
			order.UpdatedAt = now
		}
		if order.FulfillmentStatus == 0 {
			order.FulfillmentStatus = domain.Created
		}
		if order.Status == 0 {
			order.Status = order.FulfillmentStatus
		}
		if order.PaymentStatus == "" {
			order.PaymentStatus = domain.OrderPaymentPendingFunds
		}

		var createdOrder domain.Order
		err = tx.GetContext(
			ctx,
			&createdOrder,
			query,
			order.CheckoutID,
			order.UserID,
			order.VendorID,
			order.Status,
			order.PaymentStatus,
			order.FulfillmentStatus,
			order.DeliveryAddressID,
			order.ProductID,
			order.ProductName,
			order.ProductImageURL,
			order.Quantity,
			order.UnitPrice,
			order.TotalPrice,
			order.CreatedAt,
			order.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		_, err = tx.ExecContext(
			ctx,
			paymentQuery,
			createdOrder.ID,
			order.Payment.CurrencyID,
			order.Payment.Type,
			order.Payment.Status,
		)
		if err != nil {
			return nil, err
		}

		_, err = tx.ExecContext(
			ctx,
			deliveryQuery,
			createdOrder.ID,
			order.Delivery.EntityID(),
			order.Delivery.Type,
			order.Delivery.Status,
		)
		if err != nil {
			return nil, err
		}

		createdOrder.Payment = order.Payment
		createdOrder.CurrencyID = order.Payment.CurrencyID
		createdOrder.Delivery = order.Delivery

		createdOrders = append(createdOrders, createdOrder)
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return createdOrders, nil
}
