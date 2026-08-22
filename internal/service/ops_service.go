package service

import (
	"context"
	"database/sql"
	"time"

	"go-repair-center/internal/repository"
)

type OpsService struct {
	db *sql.DB
}

type BacklogSummary struct {
	PendingOrders       int64 `json:"pending_orders"`
	WaitingVisitOrders  int64 `json:"waiting_visit_orders"`
	WaitingPickupOrders int64 `json:"waiting_pickup_orders"`
	PendingFeedbacks    int64 `json:"pending_feedbacks"`
	ActiveWarranties    int64 `json:"active_warranties"`
}

type SLAAlert struct {
	ResourceType string    `json:"resource_type"`
	ResourceID   int64     `json:"resource_id"`
	Status       string    `json:"status"`
	DueAt        time.Time `json:"due_at"`
	Severity     string    `json:"severity"`
}

func NewOpsService(db *sql.DB) *OpsService {
	return &OpsService{db: db}
}

func (s *OpsService) Backlog(ctx context.Context) (*BacklogSummary, error) {
	result := &BacklogSummary{}
	var err error
	if result.PendingOrders, err = s.countByStatus(ctx, "repair_orders", "PENDING"); err != nil {
		return nil, err
	}
	if result.WaitingVisitOrders, err = s.countByStatus(ctx, "repair_orders", "WAITING_VISIT"); err != nil {
		return nil, err
	}
	if result.WaitingPickupOrders, err = s.countByStatus(ctx, "repair_executions", "WAITING_PICKUP"); err != nil {
		return nil, err
	}
	if result.PendingFeedbacks, err = s.countByStatus(ctx, "feedbacks", "PENDING"); err != nil {
		return nil, err
	}
	if result.ActiveWarranties, err = s.countByStatus(ctx, "warranties", "ACTIVE"); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *OpsService) AppointmentConflicts(ctx context.Context, technicianID int64, start time.Time, end time.Time) ([]repository.Record, error) {
	rows, err := s.db.QueryContext(ctx, `
SELECT id, order_number, appointment_time, appointment_end, status
FROM repair_orders
WHERE technician_id = ?
  AND deleted_at IS NULL
  AND status IN ('DISPATCHED', 'WAITING_VISIT', 'WAITING_DELIVERY', 'ACCEPTED')
  AND appointment_time IS NOT NULL
  AND appointment_end IS NOT NULL
  AND appointment_time < ?
  AND appointment_end > ?
ORDER BY appointment_time ASC
`, technicianID, end, start)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]repository.Record, 0, 8)
	for rows.Next() {
		var id int64
		var orderNumber string
		var appointmentTime time.Time
		var appointmentEnd time.Time
		var status string
		if err := rows.Scan(&id, &orderNumber, &appointmentTime, &appointmentEnd, &status); err != nil {
			return nil, err
		}
		items = append(items, repository.Record{
			"id":               id,
			"order_number":     orderNumber,
			"appointment_time": appointmentTime,
			"appointment_end":  appointmentEnd,
			"status":           status,
		})
	}
	return items, rows.Err()
}

func (s *OpsService) SLAAlerts(ctx context.Context, now time.Time) ([]SLAAlert, error) {
	rows, err := s.db.QueryContext(ctx, `
SELECT id, status, expected_completed_at
FROM repair_orders
WHERE deleted_at IS NULL
  AND expected_completed_at IS NOT NULL
  AND status IN ('PENDING', 'DISPATCHED', 'WAITING_VISIT', 'WAITING_DELIVERY', 'ACCEPTED')
  AND expected_completed_at < ?
ORDER BY expected_completed_at ASC
`, now)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]SLAAlert, 0, 16)
	for rows.Next() {
		var id int64
		var status string
		var dueAt time.Time
		if err := rows.Scan(&id, &status, &dueAt); err != nil {
			return nil, err
		}
		severity := "MEDIUM"
		if now.Sub(dueAt) > 24*time.Hour {
			severity = "HIGH"
		}
		items = append(items, SLAAlert{
			ResourceType: "repair_order",
			ResourceID:   id,
			Status:       status,
			DueAt:        dueAt,
			Severity:     severity,
		})
	}
	return items, rows.Err()
}

func (s *OpsService) WarrantyExpiring(ctx context.Context, withinDays int) ([]repository.Record, error) {
	if withinDays <= 0 {
		withinDays = 7
	}
	rows, err := s.db.QueryContext(ctx, `
SELECT id, repair_order_id, end_date, status
FROM warranties
WHERE deleted_at IS NULL
  AND status='ACTIVE'
  AND end_date BETWEEN CURRENT_DATE() AND DATE_ADD(CURRENT_DATE(), INTERVAL ? DAY)
ORDER BY end_date ASC
`, withinDays)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]repository.Record, 0, 16)
	for rows.Next() {
		var id int64
		var repairOrderID int64
		var endDate time.Time
		var status string
		if err := rows.Scan(&id, &repairOrderID, &endDate, &status); err != nil {
			return nil, err
		}
		items = append(items, repository.Record{
			"id":              id,
			"repair_order_id": repairOrderID,
			"end_date":        endDate,
			"status":          status,
		})
	}
	return items, rows.Err()
}

func (s *OpsService) countByStatus(ctx context.Context, table string, status string) (int64, error) {
	query := "SELECT COUNT(*) FROM " + table + " WHERE status = ? AND deleted_at IS NULL"
	var count int64
	err := s.db.QueryRowContext(ctx, query, status).Scan(&count)
	return count, err
}
