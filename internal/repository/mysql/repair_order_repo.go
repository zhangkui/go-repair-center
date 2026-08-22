package mysql

import (
    "context"
    "database/sql"
    "go-repair-center/internal/repository"
)

type RepairOrderRepository struct {
    *CRUD
    db *sql.DB
}

func NewRepairOrderRepository(db *sql.DB) *RepairOrderRepository {
    return &RepairOrderRepository{
        CRUD: NewCRUD(db, "repair_orders", []string{"order_number","fault_description"}, []string{"order_number","customer_id","device_id","original_order_id","fault_description","urgency","service_method","expected_completed_at","appointment_time","appointment_end","technician_id","dispatched_at","notification_status","status"}),
        db: db,
    }
}

func (r *RepairOrderRepository) CountByStatus(ctx context.Context, status string) (int64, error) {
    var count int64
    err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM repair_orders WHERE status=? AND deleted_at IS NULL", status).Scan(&count)
    return count, err
}

func (r *RepairOrderRepository) ListRecent(ctx context.Context, limit int) ([]repository.Record, error) {
    if limit < 1 || limit > 100 { limit = 20 }
    rows, err := r.db.QueryContext(ctx, "SELECT * FROM repair_orders WHERE deleted_at IS NULL ORDER BY created_at DESC LIMIT ?", limit)
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

func (r *RepairOrderRepository) Exists(ctx context.Context, id int64) (bool, error) {
    var exists bool
    err := r.db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM repair_orders WHERE id=? AND deleted_at IS NULL)", id).Scan(&exists)
    return exists, err
}