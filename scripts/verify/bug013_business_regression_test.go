package verify_test

import (
	"testing"

	"go-repair-center/internal/service"
)

func TestBug013_BusinessRegression(t *testing.T) {
	if got := service.SingleSessionSummaryKey(7); got != "sessions:user:7" {
		t.Fatalf("session summary cache key = %q, want sessions:user:7", got)
	}
}
