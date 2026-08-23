package verify_test

import (
	"testing"
	"time"

	"go-repair-center/internal/service"
)

func TestBug018_BusinessRegression(t *testing.T) {
	scheduled := time.Date(2026, 8, 23, 10, 0, 0, 0, time.UTC)
	message := service.NewNotificationService(nil).BuildFeedbackReminder("R-1", scheduled)
	if message.Recipient != "service-team" {
		t.Fatalf("feedback reminder recipient = %q, want service-team", message.Recipient)
	}
	if message.Title != "Feedback Reminder" || message.Body == "" || message.Metadata["order_number"] != "R-1" {
		t.Fatalf("feedback reminder content changed: %#v", message)
	}
}
