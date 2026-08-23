package verify_test

import (
	"testing"

	"go-repair-center/internal/service"
)

func TestBug012_BusinessRegression(t *testing.T) {
	decision := service.NewApprovalService(1000, 90).ReviewQuotation(1, 1000, 0, false, "")
	if decision.NeedReview {
		t.Fatal("quotation exactly at threshold must not require review")
	}
	if !service.NewApprovalService(1000, 90).ReviewQuotation(1, 1000.01, 0, false, "").NeedReview {
		t.Fatal("quotation above threshold must require review")
	}
}
