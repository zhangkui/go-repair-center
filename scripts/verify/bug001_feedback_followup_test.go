package verify_test

import (
	"testing"
	"time"

	"go-repair-center/internal/repository"
	"go-repair-center/internal/service"
)

func TestBug001_FeedbackTaskUsesDeliveredAt(t *testing.T) {
	scheduler := service.NewSchedulerService(3, 7)
	deliveredAt := time.Date(2026, time.August, 20, 10, 0, 0, 0, time.UTC)
	updatedAt := deliveredAt.Add(48 * time.Hour)

	tasks := scheduler.BuildRepairOrderTasks(repository.Record{
		"id":           int64(1001),
		"delivered_at": deliveredAt,
		"updated_at":   updatedAt,
	})

	if len(tasks) != 1 {
		t.Fatalf("expected 1 feedback task after delivery, got %d", len(tasks))
	}

	wantRunAt := deliveredAt.AddDate(0, 0, 3)
	if !tasks[0].RunAt.Equal(wantRunAt) {
		t.Fatalf("expected follow-up at %s based on delivered_at, got %s", wantRunAt.Format(time.RFC3339), tasks[0].RunAt.Format(time.RFC3339))
	}
}
