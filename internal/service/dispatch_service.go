package service

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"go-repair-center/internal/repository"
)

type DispatchService struct {
	db            *sql.DB
	ops           *OpsService
	orderNumbers  *OrderNumberService
	notifications *NotificationService
}

type BatchDispatchRequest struct {
	OrderIDs         []int64   `json:"order_ids"`
	TechnicianID     int64     `json:"technician_id"`
	TechnicianName   string    `json:"technician_name"`
	AppointmentTime  time.Time `json:"appointment_time"`
	AppointmentEnd   time.Time `json:"appointment_end"`
	ServiceMethod    string    `json:"service_method"`
	NotifyTechnician bool      `json:"notify_technician"`
}

type BatchDispatchResult struct {
	Dispatched    []repository.Record `json:"dispatched"`
	Skipped       []repository.Record `json:"skipped"`
	Conflicts     []repository.Record `json:"conflicts"`
	Notifications []string            `json:"notifications"`
}

func NewDispatchService(db *sql.DB, ops *OpsService, orderNumbers *OrderNumberService, notifications *NotificationService) *DispatchService {
	return &DispatchService{
		db:            db,
		ops:           ops,
		orderNumbers:  orderNumbers,
		notifications: notifications,
	}
}

func (s *DispatchService) CheckConflicts(ctx context.Context, technicianID int64, appointmentTime time.Time, appointmentEnd time.Time) ([]repository.Record, error) {
	return s.ops.AppointmentConflicts(ctx, technicianID, appointmentTime, appointmentEnd)
}

func (s *DispatchService) BatchDispatch(ctx context.Context, request BatchDispatchRequest) (*BatchDispatchResult, error) {
	result := newBatchDispatchResult(len(request.OrderIDs))
	if len(request.OrderIDs) == 0 {
		return result, nil
	}

	conflicts, err := s.CheckConflicts(ctx, request.TechnicianID, request.AppointmentTime, request.AppointmentEnd)
	if err != nil {
		return nil, err
	}
	if len(conflicts) > 0 {
		result.addConflicts(conflicts)
		return result, nil
	}

	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	sequenceSeed, err := s.dailySequenceSeed(ctx, tx, request.AppointmentTime)
	if err != nil {
		return nil, err
	}

	for index, orderID := range request.OrderIDs {
		record, dispatchErr := s.dispatchOne(ctx, tx, orderID, sequenceSeed+index+1, request)
		if dispatchErr != nil {
			return nil, dispatchErr
		}
		if fmt.Sprint(record["dispatch_result"]) == "skipped" {
			result.Skipped = append(result.Skipped, record)
			continue
		}
		result.Dispatched = append(result.Dispatched, record)
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	if request.NotifyTechnician && s.notifications != nil {
		notices := make([]NotificationMessage, 0, len(result.Dispatched))
		for _, item := range result.Dispatched {
			notices = append(notices, s.notifications.BuildDispatchNotice(
				fmt.Sprint(item["order_number"]),
				request.TechnicianName,
				request.AppointmentTime,
			))
			result.addNotification(fmt.Sprintf("queued:%v", item["order_number"]))
		}
		_ = s.notifications.QueueMany(ctx, notices)
	}

	return result, nil
}

func (s *DispatchService) ReworkChain(ctx context.Context, orderID int64) ([]repository.Record, error) {
	rows, err := s.db.QueryContext(ctx, `
WITH RECURSIVE ancestors AS (
  SELECT id, order_number, original_order_id, status, customer_id, device_id, created_at
  FROM repair_orders
  WHERE id = ? AND deleted_at IS NULL
  UNION ALL
  SELECT ro.id, ro.order_number, ro.original_order_id, ro.status, ro.customer_id, ro.device_id, ro.created_at
  FROM repair_orders ro
  JOIN ancestors a ON a.original_order_id = ro.id
  WHERE ro.deleted_at IS NULL
),
root_order AS (
  SELECT id
  FROM ancestors
  ORDER BY CASE WHEN original_order_id IS NULL THEN 0 ELSE 1 END, created_at ASC, id ASC
  LIMIT 1
),
chain AS (
  SELECT ro.id, ro.order_number, ro.original_order_id, ro.status, ro.customer_id, ro.device_id, ro.created_at
  FROM repair_orders ro
  JOIN root_order r ON r.id = ro.id
  WHERE ro.deleted_at IS NULL
  UNION ALL
  SELECT ro.id, ro.order_number, ro.original_order_id, ro.status, ro.customer_id, ro.device_id, ro.created_at
  FROM repair_orders ro
  JOIN chain c ON ro.original_order_id = c.id
  WHERE ro.deleted_at IS NULL
)
SELECT id, order_number, original_order_id, status, customer_id, device_id, created_at
FROM chain
ORDER BY created_at ASC, id ASC
`, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]repository.Record, 0, 8)
	for rows.Next() {
		var id int64
		var orderNumber string
		var originalOrderID sql.NullInt64
		var status string
		var customerID int64
		var deviceID int64
		var createdAt time.Time
		if err := rows.Scan(&id, &orderNumber, &originalOrderID, &status, &customerID, &deviceID, &createdAt); err != nil {
			return nil, err
		}
		record := repository.Record{
			"id":                id,
			"order_number":      orderNumber,
			"status":            status,
			"customer_id":       customerID,
			"device_id":         deviceID,
			"created_at":        createdAt,
			"has_original":      originalOrderID.Valid,
			"original_order_id": originalOrderID.Int64,
		}
		items = append(items, record)
	}
	return items, rows.Err()
}

func (s *DispatchService) dispatchOne(ctx context.Context, tx *sql.Tx, orderID int64, sequence int, request BatchDispatchRequest) (repository.Record, error) {
	var orderNumber sql.NullString
	var status string
	err := tx.QueryRowContext(ctx, `
SELECT order_number, status
FROM repair_orders
WHERE id = ? AND deleted_at IS NULL
FOR UPDATE
`, orderID).Scan(&orderNumber, &status)
	if err != nil {
		return nil, err
	}

	if status != "PENDING" {
		return repository.Record{
			"id":              orderID,
			"status":          status,
			"dispatch_result": "skipped",
			"reason":          "status is not pending",
		}, nil
	}

	finalOrderNumber := orderNumber.String
	if !orderNumber.Valid || orderNumber.String == "" {
		finalOrderNumber = s.orderNumbers.GenerateRepairOrderNumberAt(request.AppointmentTime, sequence)
	}

	notificationStatus := "SKIPPED"
	if request.NotifyTechnician {
		notificationStatus = "QUEUED"
	}

	_, err = tx.ExecContext(ctx, `
UPDATE repair_orders
SET order_number = ?,
    technician_id = ?,
    dispatched_at = NOW(),
    appointment_time = ?,
    appointment_end = ?,
    service_method = ?,
    notification_status = ?,
    status = 'DISPATCHED',
    updated_at = NOW()
WHERE id = ?
`, finalOrderNumber, request.TechnicianID, request.AppointmentTime, request.AppointmentEnd, request.ServiceMethod, notificationStatus, orderID)
	if err != nil {
		return nil, err
	}

	return repository.Record{
		"id":                  orderID,
		"order_number":        finalOrderNumber,
		"technician_id":       request.TechnicianID,
		"appointment_time":    request.AppointmentTime,
		"appointment_end":     request.AppointmentEnd,
		"service_method":      request.ServiceMethod,
		"notification_status": notificationStatus,
		"status":              "DISPATCHED",
		"dispatch_result":     "dispatched",
	}, nil
}

func (s *DispatchService) dailySequenceSeed(ctx context.Context, tx *sql.Tx, appointment time.Time) (int, error) {
	var count int
	err := tx.QueryRowContext(ctx, `
SELECT COUNT(*)
FROM repair_orders
WHERE DATE(created_at) = DATE(?) AND deleted_at IS NULL
`, appointment).Scan(&count)
	return count, err
}
