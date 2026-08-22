package service

import (
	"context"
	"database/sql"
	"time"
)

type MaintenanceService struct {
	db    *sql.DB
	cache *CacheService
}

type HealthSnapshot struct {
	CheckedAt             time.Time `json:"checked_at"`
	OpenRepairOrders      int64     `json:"open_repair_orders"`
	OpenExecutions        int64     `json:"open_executions"`
	ExpiredQuotations     int64     `json:"expired_quotations"`
	ExpiredWarranties     int64     `json:"expired_warranties"`
	PendingFeedbacks      int64     `json:"pending_feedbacks"`
	LowStockParts         int64     `json:"low_stock_parts"`
}

func NewMaintenanceService(db *sql.DB, cache *CacheService) *MaintenanceService {
	return &MaintenanceService{db: db, cache: cache}
}

func (s *MaintenanceService) Snapshot(ctx context.Context) (*HealthSnapshot, error) {
	if s.cache != nil {
		cached := &HealthSnapshot{}
		if ok, err := s.cache.GetJSON(ctx, "maintenance:snapshot", cached); err == nil && ok {
			return cached, nil
		}
	}

	snapshot := &HealthSnapshot{CheckedAt: time.Now()}
	var err error
	if snapshot.OpenRepairOrders, err = s.count(ctx, "SELECT COUNT(*) FROM repair_orders WHERE deleted_at IS NULL AND status NOT IN ('COMPLETED','CANCELLED')"); err != nil {
		return nil, err
	}
	if snapshot.OpenExecutions, err = s.count(ctx, "SELECT COUNT(*) FROM repair_executions WHERE deleted_at IS NULL AND status NOT IN ('DELIVERED')"); err != nil {
		return nil, err
	}
	if snapshot.ExpiredQuotations, err = s.count(ctx, "SELECT COUNT(*) FROM quotations WHERE deleted_at IS NULL AND valid_until < NOW() AND status NOT IN ('EXPIRED','CONVERTED')"); err != nil {
		return nil, err
	}
	if snapshot.ExpiredWarranties, err = s.count(ctx, "SELECT COUNT(*) FROM warranties WHERE deleted_at IS NULL AND end_date < CURRENT_DATE() AND status='ACTIVE'"); err != nil {
		return nil, err
	}
	if snapshot.PendingFeedbacks, err = s.count(ctx, "SELECT COUNT(*) FROM feedbacks WHERE deleted_at IS NULL AND status='PENDING'"); err != nil {
		return nil, err
	}
	if snapshot.LowStockParts, err = s.count(ctx, "SELECT COUNT(*) FROM parts WHERE deleted_at IS NULL AND stock_quantity <= safety_stock"); err != nil {
		return nil, err
	}

	if s.cache != nil {
		_ = s.cache.SetJSON(ctx, "maintenance:snapshot", snapshot, time.Minute)
	}
	return snapshot, nil
}

func (s *MaintenanceService) count(ctx context.Context, query string) (int64, error) {
	var value int64
	err := s.db.QueryRowContext(ctx, query).Scan(&value)
	return value, err
}

func (s *MaintenanceService) ClearSnapshot(ctx context.Context) error {
	if s.cache == nil {
		return nil
	}
	return s.cache.Delete(ctx, "maintenance:snapshot")
}
