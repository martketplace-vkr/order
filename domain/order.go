package domain

import (
	"time"

	cartorderpb "github.com/martketplace-vkr/cart/pkg/api/grpc/v1/order"
)

type (
	OrderStatus int64
)

const (
	Created OrderStatus = iota + 1
	WaitingForPayment
	Assembly
	DeliveryToPickUp
	DeliveryToClient
	WaitingPickUp
	Success
	CancelledByClient
	CancelledBySeller
)

func (o OrderStatus) String() string {
	return StatusToString[o]
}

const (
	CreatedString           = "created"
	WaitingForPaymentString = "waiting_for_payment"
	AssemblyString          = "assembly"
	DeliveryToPickUpString  = "delivery_to_pick_up"
	DeliveryToClientString  = "delivery_to_client"
	WaitingPickUpString     = "waiting_pick_up"
	SuccessString           = "success"
	CancelledByClientString = "canecelled_by_client"
	CancelledBySellerString = "canecelled_by_seller"
)

var (
	StatusToString map[OrderStatus]string = map[OrderStatus]string{
		Created:           CreatedString,
		WaitingForPayment: WaitingForPaymentString,
		Assembly:          AssemblyString,
		DeliveryToPickUp:  DeliveryToPickUpString,
		DeliveryToClient:  DeliveryToClientString,
		WaitingPickUp:     WaitingPickUpString,
		Success:           SuccessString,
		CancelledByClient: CancelledByClientString,
		CancelledBySeller: CancelledBySellerString,
	}
)

type Order struct {
	ID              int64       `db:"id"`
	CheckoutID      string      `db:"checkout_id"`
	UserID          int64       `db:"user_id"`
	VendorID        int64       `db:"vendor_id"`
	Status          OrderStatus `db:"status"`
	ProductID       int64       `db:"product_id"`
	ProductName     string      `db:"product_name"`
	ProductImageURL string      `db:"product_image_url"`
	Quantity        int64       `db:"quantity"`
	UnitPrice       string      `db:"unit_price"`
	TotalPrice      string      `db:"total_price"`
	Comment         string      `db:"-"`
	CreatedAt       time.Time   `db:"created_at"`
	UpdatedAt       time.Time   `db:"updated_at"`
	Payment         Payment     `db:"-"`
	Delivery        Delivery    `db:"-"`
}

type OrderList []Order

func OrderListFromReservation(reservation *cartorderpb.CheckoutReservation) OrderList {
	if reservation == nil {
		return nil
	}

	orders := make(OrderList, 0, len(reservation.GetItems()))
	now := time.Now().UTC()

	for _, item := range reservation.GetItems() {
		if item == nil {
			continue
		}

		orders = append(orders, Order{
			CheckoutID:      reservation.GetCheckoutId(),
			UserID:          reservation.GetUserId(),
			VendorID:        item.GetVendorId(),
			Status:          Created,
			ProductID:       item.GetProductId(),
			ProductName:     item.GetProductName(),
			ProductImageURL: item.GetImageUrl(),
			Quantity:        int64(item.GetQuantity()),
			UnitPrice:       item.GetUnitPrice(),
			TotalPrice:      item.GetTotalPrice(),
			CreatedAt:       now,
			UpdatedAt:       now,
		})
	}

	return orders
}
