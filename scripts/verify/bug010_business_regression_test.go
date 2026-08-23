package verify_test

import (
	"testing"

	"go-repair-center/internal/service"
)

func TestBug010_BusinessRegression(t *testing.T) {
	got := service.NewSchedulerService(3, 7).MergeTaskBatches(
		[]service.ScheduledTask{{TaskType: "feedback_follow_up"}, {TaskType: "quotation_expiry_check"}},
		[]service.ScheduledTask{{TaskType: "warranty_expiry_notice"}, {TaskType: "dispatch"}},
	)
	want := []string{"feedback_follow_up", "quotation_expiry_check", "warranty_expiry_notice", "dispatch"}
	if len(got) != len(want) {
		t.Fatalf("merged task count = %d, want %d: %#v", len(got), len(want), got)
	}
	for index, expected := range want {
		if got[index].TaskType != expected {
			t.Fatalf("merged task %d = %q, want %q", index, got[index].TaskType, expected)
		}
	}
}
