package mysql

import (
    "context"
    "database/sql"
    "go-repair-center/internal/repository"
)

type FeedbackRepository struct {
    *CRUD
    db *sql.DB
}

func NewFeedbackRepository(db *sql.DB) *FeedbackRepository {
    return &FeedbackRepository{
        CRUD: NewCRUD(db, "feedbacks", []string{"complaint"}, []string{"repair_order_id","scheduled_at","completed_at","method","complaint_type","complaint","resolution","rework_order_id","rework_reason","rework_cost","status"}),
        db: db,
    }
}

func (r *FeedbackRepository) CountByStatus(ctx context.Context, status string) (int64, error) {
    var count int64
    err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM feedbacks WHERE status=? AND deleted_at IS NULL", status).Scan(&count)
    return count, err
}

func (r *FeedbackRepository) ListRecent(ctx context.Context, limit int) ([]repository.Record, error) {
    if limit < 1 || limit > 100 { limit = 20 }
    rows, err := r.db.QueryContext(ctx, "SELECT * FROM feedbacks WHERE deleted_at IS NULL ORDER BY created_at DESC LIMIT ?", limit)
    if err != nil { return nil, err }
    defer rows.Close()
    items := make([]repository.Record, 0, limit)
    for rows.Next() {
        item, scanErr := scanRecord(rows)
        if scanErr != nil { return nil, scanErr }
        items = append(items, item)
    }
    return items, rows.Err()
}

func (r *FeedbackRepository) Exists(ctx context.Context, id int64) (bool, error) {
    var exists bool
    err := r.db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM feedbacks WHERE id=? AND deleted_at IS NULL)", id).Scan(&exists)
    return exists, err
}