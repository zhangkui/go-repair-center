package service

import (
    "context"
    "go-repair-center/internal/repository"
)

type WarrantyService struct {
    *ResourceService
}

func NewWarrantyService(repo repository.CRUDRepository) *WarrantyService {
    return &WarrantyService{ResourceService: NewResourceService(repo, "warranty", "repair_order_id","warranty_type")}
}

func (s *WarrantyService) Search(ctx context.Context, keyword string, page, pageSize int) (repository.Page, error) {
    return s.List(ctx, page, pageSize, keyword, "")
}

func (s *WarrantyService) Active(ctx context.Context, page, pageSize int) (repository.Page, error) {
    return s.List(ctx, page, pageSize, "", "ACTIVE")
}

func (s *WarrantyService) Archive(ctx context.Context, id int64) error {
    return s.Delete(ctx, id)
}

var warrantyTransitions = map[string][]string{"ACTIVE":{"EXPIRED","VOID"},"EXPIRED":{},"VOID":{}}

func (s *WarrantyService) ChangeStatus(ctx context.Context, id int64, current, next, reason string) (repository.Record, error) {
    if err := ValidateTransition(current, next, warrantyTransitions); err != nil { return nil, err }
    return s.Update(ctx, id, repository.Record{"status": next, "status_reason": reason})
}