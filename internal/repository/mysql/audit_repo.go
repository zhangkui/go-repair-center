package mysql

import (
    "context"
    "database/sql"
    "go-repair-center/internal/repository"
)

type AuditRepository struct {
    *CRUD
    db *sql.DB
}

func NewAuditRepository(db *sql.DB) *AuditRepository {
    return &AuditRepository{
        CRUD: NewCRUD(db, "audit_logs", []string{"user_id"}, []string{"user_id","action","resource_type","resource_id","before_value","after_value","ip_address","user_agent","request_id"}),
        db: db,
    }
}

func (r *AuditRepository) CountByStatus(ctx context.Context, status string) (int64, error) {
    var count int64
    err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM audit_logs WHERE status=? AND deleted_at IS NULL", status).Scan(&count)
    return count, err
}

func (r *AuditRepository) ListRecent(ctx context.Context, limit int) ([]repository.Record, error) {
    if limit < 1 || limit > 100 { limit = 20 }
    rows, err := r.db.QueryContext(ctx, "SELECT * FROM audit_logs WHERE deleted_at IS NULL ORDER BY created_at DESC LIMIT ?", limit)
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

func (r *AuditRepository) Exists(ctx context.Context, id int64) (bool, error) {
    var exists bool
    err := r.db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM audit_logs WHERE id=? AND deleted_at IS NULL)", id).Scan(&exists)
    return exists, err
}