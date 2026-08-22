package mysql

import (
    "context"
    "database/sql"
    "go-repair-center/internal/repository"
)

type DeviceRepository struct {
    *CRUD
    db *sql.DB
}

func NewDeviceRepository(db *sql.DB) *DeviceRepository {
    return &DeviceRepository{
        CRUD: NewCRUD(db, "devices", []string{"serial_number"}, []string{"customer_id","brand","model","serial_number","category","appearance_photos","purchase_date","purchase_channel","purchase_price","warranty_start","warranty_end","warranty_status","warranty_voucher","status"}),
        db: db,
    }
}

func (r *DeviceRepository) CountByStatus(ctx context.Context, status string) (int64, error) {
    var count int64
    err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM devices WHERE status=? AND deleted_at IS NULL", status).Scan(&count)
    return count, err
}

func (r *DeviceRepository) ListRecent(ctx context.Context, limit int) ([]repository.Record, error) {
    if limit < 1 || limit > 100 { limit = 20 }
    rows, err := r.db.QueryContext(ctx, "SELECT * FROM devices WHERE deleted_at IS NULL ORDER BY created_at DESC LIMIT ?", limit)
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

func (r *DeviceRepository) Exists(ctx context.Context, id int64) (bool, error) {
    var exists bool
    err := r.db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM devices WHERE id=? AND deleted_at IS NULL)", id).Scan(&exists)
    return exists, err
}