package verify_test

import (
	"testing"

	"go-repair-center/internal/service"
)

func TestBug016_BusinessRegression(t *testing.T) {
	approval := service.NewApprovalService(1000, 90)
	items := approval.FilterPending([]service.ApprovalDecision{
		{ResourceType: "quotation", ResourceID: 1, NeedReview: true, Approved: true},
		{ResourceType: "quotation", ResourceID: 2, NeedReview: true, Approved: false},
		{ResourceType: "quotation", ResourceID: 3, NeedReview: false, Approved: true},
	})
	if len(items) != 1 || items[0].ResourceID != 2 {
		t.Fatalf("pending decisions = %#v, want only unapproved review %d", items, 2)
	}
}
