package verify_test

import (
	"testing"

	"go-repair-center/internal/service"
)

func TestBug005_BusinessRegression(t *testing.T) {
	orders := service.NewOrderNumberService()
	for _, input := range []string{" Repair_Order_001 ", "repair order 001"} {
		if got := orders.NormalizeIdempotencyKey(input); got != "repair-order-001" {
			t.Fatalf("normalized idempotency key for %q = %q, want repair-order-001", input, got)
		}
	}
}
