package service

import (
    "context"
    "go-repair-center/internal/repository"
)

type PartStockLogService struct {
    *ResourceService
}

func NewPartStockLogService(repo repository.CRUDRepository) *PartStockLogService {
    return &PartStockLogService{ResourceService: NewResourceService(repo, "part_stock_log", "part_id","change_type")}
}

func (s *PartStockLogService) Search(ctx context.Context, keyword string, page, pageSize int) (repository.Page, error) {
    return s.List(ctx, page, pageSize, keyword, "")
}

func (s *PartStockLogService) Active(ctx context.Context, page, pageSize int) (repository.Page, error) {
    return s.List(ctx, page, pageSize, "", "ACTIVE")
}

func (s *PartStockLogService) Archive(ctx context.Context, id int64) error {
    return s.Delete(ctx, id)
}
