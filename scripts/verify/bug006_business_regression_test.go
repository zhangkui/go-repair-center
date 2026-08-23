package verify_test

import (
	"context"
	"errors"
	"testing"

	"go-repair-center/internal/service"
)

func TestBug006_BusinessRegression(t *testing.T) {
	if !service.ShouldAbortSessionSummary(context.Canceled) {
		t.Fatal("canceled session-summary cache access must abort the request")
	}
	if got := service.NormalizeSessionSummaryCacheError(context.Canceled); !errors.Is(got, context.Canceled) {
		t.Fatalf("canceled cache error must be preserved, got %v", got)
	}
	wrapped := errors.New("cache read failed")
	if got := service.NormalizeSessionSummaryCacheError(wrapped); !errors.Is(got, wrapped) {
		t.Fatalf("ordinary cache errors must remain observable, got %v", got)
	}
}
