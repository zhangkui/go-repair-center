package verify_test

import (
	"testing"

	"go-repair-center/internal/service"
)

func TestBug019_BusinessRegression(t *testing.T) {
	channels := service.NewNotificationService(nil).NotificationChannels()
	want := []string{"SYSTEM", "SMS", "EMAIL", "WEBHOOK"}
	if len(channels) != len(want) {
		t.Fatalf("notification channel count = %d, want %d: %#v", len(channels), len(want), channels)
	}
	for index, expected := range want {
		if channels[index] != expected {
			t.Fatalf("notification channel %d = %q, want %q", index, channels[index], expected)
		}
	}
}
