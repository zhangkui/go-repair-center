package verify_test

import (
	"testing"

	"go-repair-center/internal/service"
)

func TestBug015_BusinessRegression(t *testing.T) {
	approval := service.NewApprovalService(1000, 90)
	decision := approval.ReviewQuotation(1, 10, 0, false, "")
	if err := approval.Validate(decision); err != nil {
		t.Fatalf("decision that needs no review must validate without reviewer: %v", err)
	}

	needsReview := approval.ReviewQuotation(1, 2000, 8, true, "approved")
	if err := approval.Validate(needsReview); err != nil {
		t.Fatalf("complete reviewed decision should validate: %v", err)
	}
}
