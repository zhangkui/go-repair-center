package mysql

import (
    "context"
    "database/sql"
    "go-repair-center/internal/repository"
)

type RepairExecutionRepository struct {
    *CRUD
    db *sql.DB
}

func NewRepairExecutionRepository(db *sql.DB) *RepairExecutionRepository {
    return &RepairExecutionRepository{
        CRUD: NewCRUD(db, "repair_executions", []string{"repair_order_id"}, []string{"repair_order_id","technician_id","fault_code_id","diagnosis","started_at","completed_at","status"}),
        db: db,
    }
}

func (r *RepairExecutionRepository) CountByStatus(ctx context.Context, status string) (int64, error) {
    var count int64
    err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM repair_executions WHERE status=? AND deleted_at IS NULL", status).Scan(&count)
    return count, err
}

func (r *RepairExecutionRepository) ListRecent(ctx context.Context, limit int) ([]repository.Record, error) {
    if limit < 1 || limit > 100 { limit = 20 }
    rows, err := r.db.QueryContext(ctx, "SELECT * FROM repair_executions WHERE deleted_at IS NULL ORDER BY created_at DESC LIMIT ?", limit)
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

func (r *RepairExecutionRepository) Exists(ctx context.Context, id int64) (bool, error) {
    var exists bool
    err := r.db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM repair_executions WHERE id=? AND deleted_at IS NULL)", id).Scan(&exists)
    return exists, err
}