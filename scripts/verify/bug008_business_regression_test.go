package verify_test

import (
	"context"
	"testing"

	"go-repair-center/internal/service"
)

func TestBug008_BusinessRegression(t *testing.T) {
	result, err := service.NewDispatchService(nil, nil, nil, nil).BatchDispatch(context.Background(), service.BatchDispatchRequest{})
	if err != nil {
		t.Fatalf("empty batch dispatch failed: %v", err)
	}
	if result == nil || result.Dispatched == nil || result.Skipped == nil || result.Conflicts == nil || result.Notifications == nil {
		t.Fatalf("empty batch result collections must be initialized: %#v", result)
	}
}
