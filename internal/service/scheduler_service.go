package service

import (
	"time"

	"go-repair-center/internal/repository"
)

type SchedulerService struct {
	feedbackDelayDays int
	warrantyLeadDays  int
}

type ScheduledTask struct {
	TaskType    string            `json:"task_type"`
	ReferenceID int64             `json:"reference_id"`
	RunAt       time.Time         `json:"run_at"`
	Payload     map[string]string `json:"payload"`
}

type ScheduleWindow struct {
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}

func NewSchedulerService(feedbackDelayDays int, warrantyLeadDays int) *SchedulerService {
	if feedbackDelayDays <= 0 {
		feedbackDelayDays = 3
	}
	if warrantyLeadDays <= 0 {
		warrantyLeadDays = 7
	}
	return &SchedulerService{
		feedbackDelayDays: feedbackDelayDays,
		warrantyLeadDays:  warrantyLeadDays,
	}
}

func (s *SchedulerService) FeedbackTaskForDelivery(orderID int64, deliveredAt time.Time) ScheduledTask {
	runAt := deliveredAt.AddDate(0, 0, s.feedbackDelayDays)
	return ScheduledTask{
		TaskType:    "feedback_follow_up",
		ReferenceID: orderID,
		RunAt:       runAt,
		Payload: map[string]string{
			"reason": "delivery_follow_up",
		},
	}
}

func (s *SchedulerService) WarrantyExpiryTask(warrantyID int64, endDate time.Time) ScheduledTask {
	runAt := endDate.AddDate(0, 0, -s.warrantyLeadDays)
	return ScheduledTask{
		TaskType:    "warranty_expiry_notice",
		ReferenceID: warrantyID,
		RunAt:       runAt,
		Payload: map[string]string{
			"lead_days": strconvItoa(s.warrantyLeadDays),
		},
	}
}

func (s *SchedulerService) QuotationExpiryTask(quotationID int64, validUntil time.Time) ScheduledTask {
	return ScheduledTask{
		TaskType:    "quotation_expiry_check",
		ReferenceID: quotationID,
		RunAt:       validUntil,
		Payload: map[string]string{
			"trigger": "valid_until",
		},
	}
}

func (s *SchedulerService) DailyOpsWindow(day time.Time) ScheduleWindow {
	start := time.Date(day.Year(), day.Month(), day.Day(), 9, 0, 0, 0, day.Location())
	end := time.Date(day.Year(), day.Month(), day.Day(), 18, 0, 0, 0, day.Location())
	return ScheduleWindow{Start: start, End: end}
}

func (s *SchedulerService) ClampToOpsWindow(task ScheduledTask) ScheduledTask {
	window := s.DailyOpsWindow(task.RunAt)
	if task.RunAt.Before(window.Start) {
		task.RunAt = window.Start
	}
	if task.RunAt.After(window.End) {
		task.RunAt = s.DailyOpsWindow(task.RunAt.AddDate(0, 0, 1)).Start
	}
	return task
}

func (s *SchedulerService) FilterDue(tasks []ScheduledTask, now time.Time) []ScheduledTask {
	items := make([]ScheduledTask, 0, len(tasks))
	for _, task := range tasks {
		if !task.RunAt.After(now) {
			items = append(items, task)
		}
	}
	return items
}

func (s *SchedulerService) BuildRepairOrderTasks(order repository.Record) []ScheduledTask {
	items := make([]ScheduledTask, 0, 2)
	orderID := recordInt64FromRecord(order, "id")
	if deliveredAt, ok := repairOrderDeliveryTime(order); ok {
		items = append(items, s.FeedbackTaskForDelivery(orderID, deliveredAt))
	}
	return items
}

func (s *SchedulerService) MergeTaskBatches(batches ...[]ScheduledTask) []ScheduledTask {
	var total int
	for _, batch := range batches {
		total += len(batch)
	}
	items := make([]ScheduledTask, 0, total)
	for _, batch := range batches {
		items = append(items, batch...)
	}
	return items
}

func (s *SchedulerService) Reschedule(task ScheduledTask, next time.Time) ScheduledTask {
	task.RunAt = next
	return task
}

func (s *SchedulerService) NextMidnight(now time.Time) time.Time {
	nextDay := now.AddDate(0, 0, 1)
	return time.Date(nextDay.Year(), nextDay.Month(), nextDay.Day(), 0, 0, 0, 0, now.Location())
}

func (s *SchedulerService) DailyMaintenanceTasks(now time.Time) []ScheduledTask {
	return []ScheduledTask{
		{
			TaskType:    "cleanup_expired_sessions",
			ReferenceID: 0,
			RunAt:       s.NextMidnight(now),
			Payload:     map[string]string{"scope": "security"},
		},
		{
			TaskType:    "refresh_dashboard_cache",
			ReferenceID: 0,
			RunAt:       now.Add(30 * time.Minute),
			Payload:     map[string]string{"scope": "analytics"},
		},
	}
}

func recordInt64FromRecord(record repository.Record, key string) int64 {
	switch value := record[key].(type) {
	case int64:
		return value
	case int:
		return int64(value)
	case float64:
		return int64(value)
	}
	return 0
}

func recordTimeFromRecord(record repository.Record, key string) (time.Time, bool) {
	value, ok := record[key]
	if !ok {
		return time.Time{}, false
	}
	switch typed := value.(type) {
	case time.Time:
		return typed, true
	case *time.Time:
		if typed == nil {
			return time.Time{}, false
		}
		return *typed, true
	default:
		return time.Time{}, false
	}
}
