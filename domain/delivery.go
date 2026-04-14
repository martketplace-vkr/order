package domain

type DeliveryType int64

const (
	PickUp DeliveryType = iota + 1
	ClientDelivery
)

type DeliveryStatus int64

const ()

type (
	Delivery struct {
		Type          DeliveryType
		PickUpPointID *int64
		ClientAddress *int64
		Status        DeliveryStatus
	}
)

func (d *Delivery) EntityID() int64 {
	switch d.Type {
	case PickUp:
		return *d.PickUpPointID
	case ClientDelivery:
		return *d.ClientAddress
	}

	return 0
}
