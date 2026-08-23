package verify_test

import (
	"context"
	"testing"

	"go-repair-center/internal/service"
)

func TestBug006_BusinessRegression(t *testing.T) {
	if !service.ShouldAbortSessionSummary(context.Canceled) {
		t.Fatal("canceled session-summary cache access must abort the request")
	}
}
