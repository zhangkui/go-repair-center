package verify_test

import (
	"testing"
	"time"

	"go-repair-center/internal/service"
)

func TestBug010_BusinessRegression(t *testing.T) {
	base := time.Date(2026, 8, 23, 10, 0, 0, 0, time.UTC)
	got := service.NewSchedulerService(3, 7).MergeTaskBatches(
		[]service.ScheduledTask{
			{TaskType: "feedback_follow_up", ReferenceID: 7, RunAt: base.Add(2 * time.Hour)},
			{TaskType: "quotation_expiry_check", ReferenceID: 8, RunAt: base.Add(3 * time.Hour)},
		},
		[]service.ScheduledTask{
			{TaskType: "feedback_follow_up", ReferenceID: 7, RunAt: base.Add(1 * time.Hour)},
			{TaskType: "dispatch", ReferenceID: 9, RunAt: base.Add(4 * time.Hour)},
		},
	)
	if len(got) != 3 {
		t.Fatalf("merged task count = %d, want 3: %#v", len(got), got)
	}
	for _, task := range got {
		if task.TaskType == "feedback_follow_up" && task.ReferenceID == 7 && !task.RunAt.Equal(base.Add(time.Hour)) {
			t.Fatalf("duplicate task kept %s, want earliest %s", task.RunAt, base.Add(time.Hour))
		}
	}
}
