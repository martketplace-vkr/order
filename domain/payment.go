package domain

import (
	"database/sql/driver"
	"fmt"
	"strings"
)

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

type OrderPaymentStatus string

const (
	OrderPaymentPendingFunds OrderPaymentStatus = "pending_funds"
	OrderPaymentReserved     OrderPaymentStatus = "reserved"
	OrderPaymentCaptured     OrderPaymentStatus = "captured"
	OrderPaymentReleased     OrderPaymentStatus = "released"
	OrderPaymentExpired      OrderPaymentStatus = "expired"
	OrderPaymentCancelled    OrderPaymentStatus = "cancelled"
	OrderPaymentFailed       OrderPaymentStatus = "failed"
)

var ValidOrderPaymentStatuses = map[OrderPaymentStatus]struct{}{
	OrderPaymentPendingFunds: {},
	OrderPaymentReserved:     {},
	OrderPaymentCaptured:     {},
	OrderPaymentReleased:     {},
	OrderPaymentExpired:      {},
	OrderPaymentCancelled:    {},
	OrderPaymentFailed:       {},
}

func ParseOrderPaymentStatus(value string) (OrderPaymentStatus, error) {
	status := OrderPaymentStatus(strings.TrimSpace(strings.ToLower(value)))
	if _, ok := ValidOrderPaymentStatuses[status]; !ok {
		return "", fmt.Errorf("unknown payment status: %q", value)
	}
	return status, nil
}

func (s OrderPaymentStatus) Value() (driver.Value, error) {
	if _, ok := ValidOrderPaymentStatuses[s]; !ok {
		return nil, fmt.Errorf("unknown payment status: %q", s)
	}
	return string(s), nil
}

func (s *OrderPaymentStatus) Scan(value any) error {
	if value == nil {
		*s = ""
		return nil
	}
	switch v := value.(type) {
	case string:
		parsed, err := ParseOrderPaymentStatus(v)
		if err != nil {
			return err
		}
		*s = parsed
		return nil
	case []byte:
		return s.Scan(string(v))
	default:
		return fmt.Errorf("scan payment status: unsupported type %T", value)
	}
}

type Payment struct {
	ID         int64
	CurrencyID int64
	Type       PaymentType
	Status     PaymentStatus
}
