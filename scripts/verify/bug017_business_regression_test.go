package verify_test

import (
	"testing"

	"go-repair-center/internal/service"
)

func TestBug017_BusinessRegression(t *testing.T) {
	report := &service.RevenueReport{}
	service.NormalizeReportSlices(report)
	if report.RepairTrend == nil || report.Satisfaction == nil || report.TechnicianRanking == nil {
		t.Fatalf("empty report collections must be non-nil slices: %#v", report)
	}
}
