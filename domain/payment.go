package domain

type PaymentType int64

const (
	OfflineByCash PaymentType = iota + 1
	OfflineByCard
	OnlineByCard
	OnlineByCrypto
)

type PaymentStatus int64

const (
	PendingPaymentStatus PaymentStatus = iota + 1
	SuccessPaymentStatus
	FailedPaymentStatus
)

type Payment struct {
	ID     int64
	Type   PaymentType
	Status PaymentStatus
}
