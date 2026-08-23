package verify_test

import (
	"strings"
	"testing"

	"go-repair-center/internal/service"
)

func TestBug020_BusinessRegression(t *testing.T) {
	message := service.NotificationMessage{
		Channel: "SYSTEM", Title: "T", Body: "B", Recipient: "operator",
		Metadata: map[string]string{"order_number": "R-1", "type": "feedback"},
	}
	got := service.NewNotificationService(nil).RenderPlainText(message)
	for _, expected := range []string{"[SYSTEM] T", "B", "recipient: operator", "order_number: R-1", "type: feedback"} {
		if !strings.Contains(got, expected) {
			t.Fatalf("rendered notification missing %q: %q", expected, got)
		}
	}
}
