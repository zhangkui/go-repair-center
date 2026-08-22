package service

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"go-repair-center/internal/repository"
)

type ApprovalCenterService struct {
	db       *sql.DB
	approval *ApprovalService
}

type ApprovalOverview struct {
	PendingQuotationCount int64 `json:"pending_quotation_count"`
	PendingWarrantyCount  int64 `json:"pending_warranty_count"`
	PendingReworkCount    int64 `json:"pending_rework_count"`
}

func NewApprovalCenterService(db *sql.DB, approval *ApprovalService) *ApprovalCenterService {
	return &ApprovalCenterService{db: db, approval: approval}
}

func (s *ApprovalCenterService) Overview(ctx context.Context) (*ApprovalOverview, error) {
	overview := &ApprovalOverview{}
	var err error
	if overview.PendingQuotationCount, err = s.count(ctx, `SELECT COUNT(*) FROM quotations WHERE approval_status='PENDING' AND deleted_at IS NULL`); err != nil {
		return nil, err
	}
	if overview.PendingWarrantyCount, err = s.count(ctx, `SELECT COUNT(*) FROM warranties WHERE approved_by IS NULL AND status='ACTIVE' AND deleted_at IS NULL`); err != nil {
		return nil, err
	}
	if overview.PendingReworkCount, err = s.count(ctx, `SELECT COUNT(*) FROM feedbacks WHERE rework_order_id IS NOT NULL AND status<>'CLOSED' AND deleted_at IS NULL`); err != nil {
		return nil, err
	}
	return overview, nil
}

func (s *ApprovalCenterService) PendingQuotations(ctx context.Context, page int, pageSize int) (repository.Page, error) {
	return s.listRows(ctx, page, pageSize, `
SELECT id, repair_order_id, version, total_amount, customer_status, approval_status, status, valid_until, created_at
FROM quotations
WHERE approval_status='PENDING' AND deleted_at IS NULL
ORDER BY created_at DESC
`, `
SELECT COUNT(*)
FROM quotations
WHERE approval_status='PENDING' AND deleted_at IS NULL
`)
}

func (s *ApprovalCenterService) PendingWarranties(ctx context.Context, page int, pageSize int) (repository.Page, error) {
	return s.listRows(ctx, page, pageSize, `
SELECT id, repair_order_id, warranty_type, duration_days, start_date, end_date, status, approved_by, created_at
FROM warranties
WHERE approved_by IS NULL AND status='ACTIVE' AND deleted_at IS NULL
ORDER BY created_at DESC
`, `
SELECT COUNT(*)
FROM warranties
WHERE approved_by IS NULL AND status='ACTIVE' AND deleted_at IS NULL
`)
}

func (s *ApprovalCenterService) DecideQuotation(ctx context.Context, quotationID int64, reviewerID int64, approve bool, comment string) (repository.Record, error) {
	var totalAmount float64
	err := s.db.QueryRowContext(ctx, `
SELECT total_amount
FROM quotations
WHERE id = ? AND deleted_at IS NULL
`, quotationID).Scan(&totalAmount)
	if err != nil {
		return nil, err
	}

	decision := s.approval.ReviewQuotation(quotationID, totalAmount, reviewerID, approve, comment)
	if err := s.approval.Validate(decision); err != nil {
		return nil, err
	}

	status := "REJECTED"
	if approve {
		status = "APPROVED"
	}

	_, err = s.db.ExecContext(ctx, `
UPDATE quotations
SET approval_status = ?, approved_by = ?, approved_at = ?, updated_at = NOW()
WHERE id = ? AND deleted_at IS NULL
`, status, reviewerID, time.Now(), quotationID)
	if err != nil {
		return nil, err
	}

	return repository.Record{
		"id":              quotationID,
		"approval_status": status,
		"reviewer_id":     reviewerID,
		"comment":         decision.Reason,
	}, nil
}

func (s *ApprovalCenterService) DecideWarranty(ctx context.Context, warrantyID int64, reviewerID int64, approve bool, comment string) (repository.Record, error) {
	var durationDays int
	err := s.db.QueryRowContext(ctx, `
SELECT duration_days
FROM warranties
WHERE id = ? AND deleted_at IS NULL
`, warrantyID).Scan(&durationDays)
	if err != nil {
		return nil, err
	}

	decision := s.approval.ReviewWarranty(warrantyID, durationDays, reviewerID, approve, comment)
	if err := s.approval.Validate(decision); err != nil {
		return nil, err
	}

	status := "ACTIVE"
	if !approve {
		status = "VOID"
	}
	_, err = s.db.ExecContext(ctx, `
UPDATE warranties
SET status = ?, approved_by = ?, updated_at = NOW()
WHERE id = ? AND deleted_at IS NULL
`, status, reviewerID, warrantyID)
	if err != nil {
		return nil, err
	}

	return repository.Record{
		"id":          warrantyID,
		"status":      status,
		"reviewer_id": reviewerID,
		"comment":     decision.Reason,
	}, nil
}

func (s *ApprovalCenterService) count(ctx context.Context, query string) (int64, error) {
	var count int64
	err := s.db.QueryRowContext(ctx, query).Scan(&count)
	return count, err
}

func (s *ApprovalCenterService) listRows(ctx context.Context, page int, pageSize int, query string, countQuery string) (repository.Page, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	total, err := s.count(ctx, countQuery)
	if err != nil {
		return repository.Page{}, err
	}

	rows, err := s.db.QueryContext(ctx, query+` LIMIT ? OFFSET ?`, pageSize, (page-1)*pageSize)
	if err != nil {
		return repository.Page{}, err
	}
	defer rows.Close()

	items := make([]repository.Record, 0, pageSize)
	for rows.Next() {
		item, scanErr := scanApprovalRow(rows)
		if scanErr != nil {
			return repository.Page{}, scanErr
		}
		items = append(items, item)
	}

	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))
	return repository.Page{
		Items:      items,
		Page:       page,
		PageSize:   pageSize,
		Total:      total,
		TotalPages: totalPages,
	}, rows.Err()
}

func scanApprovalRow(rows *sql.Rows) (repository.Record, error) {
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	values := make([]any, len(columns))
	pointers := make([]any, len(columns))
	for index := range values {
		pointers[index] = &values[index]
	}
	if err := rows.Scan(pointers...); err != nil {
		return nil, err
	}
	record := make(repository.Record, len(columns))
	for index, column := range columns {
		value := values[index]
		if raw, ok := value.([]byte); ok {
			record[column] = string(raw)
		} else {
			record[column] = value
		}
	}
	return record, nil
}

func (s *ApprovalCenterService) ExplainDecision(resourceType string, approve bool, comment string) string {
	action := "rejected"
	if approve {
		action = "approved"
	}
	if comment == "" {
		return fmt.Sprintf("%s %s", resourceType, action)
	}
	return fmt.Sprintf("%s %s: %s", resourceType, action, comment)
}
