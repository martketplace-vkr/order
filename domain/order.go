package domain

import (
	"database/sql/driver"
	"fmt"
	"strconv"
	"strings"
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
	value, ok := StatusToString[o]
	if !ok {
		return "unknown"
	}

	return value
}

func (o OrderStatus) Value() (driver.Value, error) {
	value, ok := StatusToString[o]
	if !ok {
		return nil, fmt.Errorf("unknown order status: %d", o)
	}

	return value, nil
}

func (o *OrderStatus) Scan(value any) error {
	switch v := value.(type) {
	case nil:
		*o = 0
		return nil
	case int64:
		*o = OrderStatus(v)
		return nil
	case int32:
		*o = OrderStatus(v)
		return nil
	case int:
		*o = OrderStatus(v)
		return nil
	case []byte:
		return o.scanString(string(v))
	case string:
		return o.scanString(v)
	default:
		return fmt.Errorf("scan order status: unsupported type %T", value)
	}
}

func (o *OrderStatus) scanString(value string) error {
	normalized := strings.TrimSpace(value)
	if normalized == "" {
		*o = 0
		return nil
	}

	if numeric, err := strconv.ParseInt(normalized, 10, 64); err == nil {
		*o = OrderStatus(numeric)
		return nil
	}

	status, ok := StringToStatus[normalized]
	if !ok {
		return fmt.Errorf("unknown order status: %q", value)
	}

	*o = status
	return nil
}

const (
	CreatedString           = "created"
	WaitingForPaymentString = "waiting_for_payment"
	AssemblyString          = "assembly"
	DeliveryToPickUpString  = "delivery_to_pick_up"
	DeliveryToClientString  = "delivery_to_client"
	WaitingPickUpString     = "waiting_pick_up"
	SuccessString           = "success"
	CancelledByClientString = "cancelled_by_client"
	CancelledBySellerString = "cancelled_by_seller"
)

var (
	StatusToString = map[OrderStatus]string{
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

	StringToStatus = map[string]OrderStatus{
		CreatedString:           Created,
		WaitingForPaymentString: WaitingForPayment,
		AssemblyString:          Assembly,
		DeliveryToPickUpString:  DeliveryToPickUp,
		DeliveryToClientString:  DeliveryToClient,
		WaitingPickUpString:     WaitingPickUp,
		SuccessString:           Success,
		CancelledByClientString: CancelledByClient,
		CancelledBySellerString: CancelledBySeller,
		"cancelled":             CancelledByClient,
		"canecelled_by_client":  CancelledByClient,
		"canecelled_by_seller":  CancelledBySeller,
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

func (list OrderList) WithPaymentAndDelivery(payment Payment, delivery Delivery) OrderList {
	for i := range list {
		list[i].Payment = payment
		list[i].Delivery = delivery
	}

	return list
}
