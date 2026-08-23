package service

import (
	"context"
	"fmt"
	"strings"
	"time"
)

type NotificationService struct {
	cache *CacheService
}

type NotificationMessage struct {
	Channel     string            `json:"channel"`
	Title       string            `json:"title"`
	Body        string            `json:"body"`
	Recipient   string            `json:"recipient"`
	ScheduledAt time.Time         `json:"scheduled_at"`
	Metadata    map[string]string `json:"metadata"`
}

func NewNotificationService(cache *CacheService) *NotificationService {
	return &NotificationService{cache: cache}
}

func (s *NotificationService) BuildDispatchNotice(orderNumber string, technicianName string, appointment time.Time) NotificationMessage {
	return NotificationMessage{
		Channel:     "SYSTEM",
		Title:       "Dispatch Created",
		Body:        fmt.Sprintf("Order %s assigned to %s at %s", orderNumber, technicianName, appointment.Format(time.RFC3339)),
		Recipient:   technicianName,
		ScheduledAt: time.Now(),
		Metadata: map[string]string{
			"order_number": orderNumber,
			"type":         "dispatch",
		},
	}
}

func (s *NotificationService) BuildQuotationNotice(orderNumber string, amount float64, validUntil time.Time) NotificationMessage {
	return NotificationMessage{
		Channel:     "SYSTEM",
		Title:       "Quotation Ready",
		Body:        fmt.Sprintf("Order %s quotation amount %.2f valid until %s", orderNumber, amount, validUntil.Format("2006-01-02 15:04")),
		Recipient:   "customer",
		ScheduledAt: time.Now(),
		Metadata: map[string]string{
			"order_number": orderNumber,
			"type":         "quotation",
		},
	}
}

func (s *NotificationService) BuildWarrantyExpiryNotice(orderNumber string, expiryDate time.Time) NotificationMessage {
	return NotificationMessage{
		Channel:     "SYSTEM",
		Title:       "Warranty Expiring",
		Body:        fmt.Sprintf("Order %s warranty expires on %s", orderNumber, expiryDate.Format("2006-01-02")),
		Recipient:   "operator",
		ScheduledAt: time.Now(),
		Metadata: map[string]string{
			"order_number": orderNumber,
			"type":         "warranty",
		},
	}
}

func (s *NotificationService) BuildFeedbackReminder(orderNumber string, scheduledAt time.Time) NotificationMessage {
	return NotificationMessage{
		Channel:     "SYSTEM",
		Title:       "Feedback Reminder",
		Body:        fmt.Sprintf("Feedback for order %s is due at %s", orderNumber, scheduledAt.Format("2006-01-02 15:04")),
		Recipient:   feedbackReminderRecipient(),
		ScheduledAt: scheduledAt,
		Metadata: map[string]string{
			"order_number": orderNumber,
			"type":         "feedback",
		},
	}
}

func (s *NotificationService) Queue(ctx context.Context, message NotificationMessage) error {
	if s.cache == nil {
		return nil
	}
	key := s.queueKey(message)
	return s.cache.SetJSON(ctx, key, message, 24*time.Hour)
}

func (s *NotificationService) QueueMany(ctx context.Context, messages []NotificationMessage) error {
	for _, message := range messages {
		if err := s.Queue(ctx, message); err != nil {
			return err
		}
	}
	return nil
}

func (s *NotificationService) Cancel(ctx context.Context, message NotificationMessage) error {
	if s.cache == nil {
		return nil
	}
	return s.cache.Delete(ctx, s.queueKey(message))
}

func (s *NotificationService) queueKey(message NotificationMessage) string {
	parts := []string{
		"notification",
		strings.ToLower(message.Channel),
		strings.ToLower(message.Recipient),
		message.ScheduledAt.Format("20060102150405"),
	}
	return strings.Join(parts, ":")
}

func (s *NotificationService) RenderPlainText(message NotificationMessage) string {
	lines := []string{
		"[" + message.Channel + "] " + message.Title,
		message.Body,
		"recipient: " + message.Recipient,
		"scheduled_at: " + message.ScheduledAt.Format(time.RFC3339),
	}
	for key, value := range message.Metadata {
		lines = append(lines, key+": "+value)
	}
	return strings.Join(lines, "\n")
}

func (s *NotificationService) BuildBulkDispatchNotices(orderNumbers []string, technicianName string, appointment time.Time) []NotificationMessage {
	items := make([]NotificationMessage, 0, len(orderNumbers))
	for _, orderNumber := range orderNumbers {
		items = append(items, s.BuildDispatchNotice(orderNumber, technicianName, appointment))
	}
	return items
}

func (s *NotificationService) NotificationChannels() []string {
	return []string{"SYSTEM", "SMS", "EMAIL", "WEBHOOK"}
}
