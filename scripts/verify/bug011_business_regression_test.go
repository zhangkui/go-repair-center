package verify_test

import (
	"testing"

	"go-repair-center/internal/service"
)

func TestBug011_BusinessRegression(t *testing.T) {
	if got := service.ParseInventoryQuantity([]byte("12")); got != 12 {
		t.Fatalf("database quantity bytes parsed as %d, want 12", got)
	}
	if got := service.ParseInventoryQuantity([]byte("0")); got != 0 {
		t.Fatalf("zero quantity parsed as %d, want 0", got)
	}
}
