package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type AnalyticsRepo struct {
	db *sql.DB
}

type DailyCount struct {
	Day   string `json:"day"`
	Count int64  `json:"count"`
}

type SatisfactionStat struct {
	Dimension string  `json:"dimension"`
	Average   float64 `json:"average"`
}

type TechnicianStat struct {
	TechnicianID int64   `json:"technician_id"`
	Username     string  `json:"username"`
	OrderCount   int64   `json:"order_count"`
	ReworkCount  int64   `json:"rework_count"`
	AvgScore     float64 `json:"avg_score"`
}

func NewAnalyticsRepo(db *sql.DB) *AnalyticsRepo {
	return &AnalyticsRepo{db: db}
}

func (r *AnalyticsRepo) CountByStatus(ctx context.Context, table string, status string) (int64, error) {
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE status=? AND deleted_at IS NULL", table)
	var count int64
	err := r.db.QueryRowContext(ctx, query, status).Scan(&count)
	return count, err
}

func (r *AnalyticsRepo) TotalCount(ctx context.Context, table string) (int64, error) {
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE deleted_at IS NULL", table)
	var count int64
	err := r.db.QueryRowContext(ctx, query).Scan(&count)
	return count, err
}

func (r *AnalyticsRepo) DailyRepairOrders(ctx context.Context, days int) ([]DailyCount, error) {
	if days <= 0 {
		days = 7
	}
	rows, err := r.db.QueryContext(ctx, `
SELECT DATE_FORMAT(created_at, '%Y-%m-%d') AS day, COUNT(*) AS total
FROM repair_orders
WHERE deleted_at IS NULL AND created_at >= DATE_SUB(NOW(), INTERVAL ? DAY)
GROUP BY DATE_FORMAT(created_at, '%Y-%m-%d')
ORDER BY day ASC
`, days)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]DailyCount, 0, days)
	for rows.Next() {
		var item DailyCount
		if err := rows.Scan(&item.Day, &item.Count); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *AnalyticsRepo) PendingFeedbackCount(ctx context.Context) (int64, error) {
	return r.CountByStatus(ctx, "feedbacks", "PENDING")
}

func (r *AnalyticsRepo) InProgressExecutionCount(ctx context.Context) (int64, error) {
	return r.CountByStatus(ctx, "repair_executions", "IN_PROGRESS")
}

func (r *AnalyticsRepo) WaitingPickupCount(ctx context.Context) (int64, error) {
	return r.CountByStatus(ctx, "repair_executions", "WAITING_PICKUP")
}

func (r *AnalyticsRepo) CompletedToday(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.QueryRowContext(ctx, `
SELECT COUNT(*)
FROM repair_orders
WHERE deleted_at IS NULL
  AND status='COMPLETED'
  AND DATE(updated_at)=CURRENT_DATE()
`).Scan(&count)
	return count, err
}

func (r *AnalyticsRepo) CreatedToday(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.QueryRowContext(ctx, `
SELECT COUNT(*)
FROM repair_orders
WHERE deleted_at IS NULL
  AND DATE(created_at)=CURRENT_DATE()
`).Scan(&count)
	return count, err
}

func (r *AnalyticsRepo) SatisfactionByDimension(ctx context.Context) ([]SatisfactionStat, error) {
	rows, err := r.db.QueryContext(ctx, `
SELECT dimension, ROUND(AVG(score), 2) AS avg_score
FROM feedback_scores
WHERE deleted_at IS NULL
GROUP BY dimension
ORDER BY dimension ASC
`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]SatisfactionStat, 0, 8)
	for rows.Next() {
		var item SatisfactionStat
		if err := rows.Scan(&item.Dimension, &item.Average); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *AnalyticsRepo) TechnicianPerformance(ctx context.Context, since time.Time, limit int) ([]TechnicianStat, error) {
	if limit <= 0 || limit > 100 {
		limit = 10
	}
	rows, err := r.db.QueryContext(ctx, `
SELECT
  u.id,
  u.username,
  COUNT(DISTINCT re.id) AS order_count,
  COUNT(DISTINCT CASE WHEN f.rework_order_id IS NOT NULL THEN f.id END) AS rework_count,
  COALESCE(ROUND(AVG(fs.score), 2), 0) AS avg_score
FROM users u
LEFT JOIN repair_executions re ON re.technician_id=u.id AND re.deleted_at IS NULL AND re.created_at >= ?
LEFT JOIN feedbacks f ON f.repair_order_id=re.repair_order_id AND f.deleted_at IS NULL
LEFT JOIN feedback_scores fs ON fs.feedback_id=f.id AND fs.deleted_at IS NULL
WHERE u.deleted_at IS NULL
GROUP BY u.id, u.username
ORDER BY order_count DESC, avg_score DESC
LIMIT ?
`, since, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	stats := make([]TechnicianStat, 0, limit)
	for rows.Next() {
		var item TechnicianStat
		if err := rows.Scan(&item.TechnicianID, &item.Username, &item.OrderCount, &item.ReworkCount, &item.AvgScore); err != nil {
			return nil, err
		}
		stats = append(stats, item)
	}
	return stats, rows.Err()
}

func (r *AnalyticsRepo) RevenueSummary(ctx context.Context, since time.Time) (map[string]float64, error) {
	row := r.db.QueryRowContext(ctx, `
SELECT
  COALESCE(SUM(amount), 0) AS total_amount,
  COALESCE(SUM(CASE WHEN paid_at >= ? THEN amount ELSE 0 END), 0) AS period_amount
FROM payment_records
WHERE deleted_at IS NULL
`, since)

	var totalAmount float64
	var periodAmount float64
	if err := row.Scan(&totalAmount, &periodAmount); err != nil {
		return nil, err
	}
	return map[string]float64{
		"total_amount":  totalAmount,
		"period_amount": periodAmount,
	}, nil
}
