package mysql

import (
    "context"
    "database/sql"
    "go-repair-center/internal/repository"
)

type FeedbackScoreRepository struct {
    *CRUD
    db *sql.DB
}

func NewFeedbackScoreRepository(db *sql.DB) *FeedbackScoreRepository {
    return &FeedbackScoreRepository{
        CRUD: NewCRUD(db, "feedback_scores", []string{"feedback_id"}, []string{"feedback_id","dimension","score","comment"}),
        db: db,
    }
}

func (r *FeedbackScoreRepository) CountByStatus(ctx context.Context, status string) (int64, error) {
    var count int64
    err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM feedback_scores WHERE status=? AND deleted_at IS NULL", status).Scan(&count)
    return count, err
}

func (r *FeedbackScoreRepository) ListRecent(ctx context.Context, limit int) ([]repository.Record, error) {
    if limit < 1 || limit > 100 { limit = 20 }
    rows, err := r.db.QueryContext(ctx, "SELECT * FROM feedback_scores WHERE deleted_at IS NULL ORDER BY created_at DESC LIMIT ?", limit)
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

func (r *FeedbackScoreRepository) Exists(ctx context.Context, id int64) (bool, error) {
    var exists bool
    err := r.db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM feedback_scores WHERE id=? AND deleted_at IS NULL)", id).Scan(&exists)
    return exists, err
}