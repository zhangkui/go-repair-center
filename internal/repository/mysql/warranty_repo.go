package mysql

import (
    "context"
    "database/sql"
    "go-repair-center/internal/repository"
)

type WarrantyRepository struct {
    *CRUD
    db *sql.DB
}

func NewWarrantyRepository(db *sql.DB) *WarrantyRepository {
    return &WarrantyRepository{
        CRUD: NewCRUD(db, "warranties", []string{"repair_order_id"}, []string{"repair_order_id","warranty_type","duration_days","start_date","end_date","status","approved_by"}),
        db: db,
    }
}

func (r *WarrantyRepository) CountByStatus(ctx context.Context, status string) (int64, error) {
    var count int64
    err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM warranties WHERE status=? AND deleted_at IS NULL", status).Scan(&count)
    return count, err
}

func (r *WarrantyRepository) ListRecent(ctx context.Context, limit int) ([]repository.Record, error) {
    if limit < 1 || limit > 100 { limit = 20 }
    rows, err := r.db.QueryContext(ctx, "SELECT * FROM warranties WHERE deleted_at IS NULL ORDER BY created_at DESC LIMIT ?", limit)
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

func (r *WarrantyRepository) Exists(ctx context.Context, id int64) (bool, error) {
    var exists bool
    err := r.db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM warranties WHERE id=? AND deleted_at IS NULL)", id).Scan(&exists)
    return exists, err
}