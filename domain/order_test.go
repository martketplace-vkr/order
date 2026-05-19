package domain

import (
	"database/sql/driver"
	"testing"
)

func TestOrderStatusScanText(t *testing.T) {
	var status OrderStatus

	if err := status.Scan("created"); err != nil {
		t.Fatalf("scan created status: %v", err)
	}

	if status != Created {
		t.Fatalf("unexpected status: got %d want %d", status, Created)
	}
}

func TestOrderStatusScanLegacyCancelledText(t *testing.T) {
	var status OrderStatus

	if err := status.Scan("canecelled_by_client"); err != nil {
		t.Fatalf("scan legacy cancelled status: %v", err)
	}

	if status != CancelledByClient {
		t.Fatalf("unexpected status: got %d want %d", status, CancelledByClient)
	}
}

func TestOrderStatusValue(t *testing.T) {
	value, err := Created.Value()
	if err != nil {
		t.Fatalf("order status value: %v", err)
	}

	if got, want := value, driver.Value(CreatedString); got != want {
		t.Fatalf("unexpected value: got %v want %v", got, want)
	}
}

func TestOrderStatusScanUnknownText(t *testing.T) {
	var status OrderStatus

	if err := status.Scan("unknown-status"); err == nil {
		t.Fatal("expected scan error")
	}
}
