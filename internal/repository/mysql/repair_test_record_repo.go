package mysql

import (
    "context"
    "database/sql"
    "go-repair-center/internal/repository"
)

type RepairTestRecordRepository struct {
    *CRUD
    db *sql.DB
}

func NewRepairTestRecordRepository(db *sql.DB) *RepairTestRecordRepository {
    return &RepairTestRecordRepository{
        CRUD: NewCRUD(db, "repair_test_records", []string{"repair_execution_id"}, []string{"repair_execution_id","test_item","result","tester_id","tested_at","notes"}),
        db: db,
    }
}

func (r *RepairTestRecordRepository) CountByStatus(ctx context.Context, status string) (int64, error) {
    var count int64
    err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM repair_test_records WHERE status=? AND deleted_at IS NULL", status).Scan(&count)
    return count, err
}

func (r *RepairTestRecordRepository) ListRecent(ctx context.Context, limit int) ([]repository.Record, error) {
    if limit < 1 || limit > 100 { limit = 20 }
    rows, err := r.db.QueryContext(ctx, "SELECT * FROM repair_test_records WHERE deleted_at IS NULL ORDER BY created_at DESC LIMIT ?", limit)
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

func (r *RepairTestRecordRepository) Exists(ctx context.Context, id int64) (bool, error) {
    var exists bool
    err := r.db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM repair_test_records WHERE id=? AND deleted_at IS NULL)", id).Scan(&exists)
    return exists, err
}