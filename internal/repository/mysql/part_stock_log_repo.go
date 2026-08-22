package mysql

import (
    "context"
    "database/sql"
    "go-repair-center/internal/repository"
)

type PartStockLogRepository struct {
    *CRUD
    db *sql.DB
}

func NewPartStockLogRepository(db *sql.DB) *PartStockLogRepository {
    return &PartStockLogRepository{
        CRUD: NewCRUD(db, "part_stock_logs", []string{"part_id"}, []string{"part_id","change_type","quantity_before","change_quantity","quantity_after","reference_type","reference_id","operator_id","remark"}),
        db: db,
    }
}

func (r *PartStockLogRepository) CountByStatus(ctx context.Context, status string) (int64, error) {
    var count int64
    err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM part_stock_logs WHERE status=? AND deleted_at IS NULL", status).Scan(&count)
    return count, err
}

func (r *PartStockLogRepository) ListRecent(ctx context.Context, limit int) ([]repository.Record, error) {
    if limit < 1 || limit > 100 { limit = 20 }
    rows, err := r.db.QueryContext(ctx, "SELECT * FROM part_stock_logs WHERE deleted_at IS NULL ORDER BY created_at DESC LIMIT ?", limit)
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

func (r *PartStockLogRepository) Exists(ctx context.Context, id int64) (bool, error) {
    var exists bool
    err := r.db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM part_stock_logs WHERE id=? AND deleted_at IS NULL)", id).Scan(&exists)
    return exists, err
}