package service

import (
    "context"
    "go-repair-center/internal/repository"
)

type RepairOrderService struct {
    *ResourceService
}

func NewRepairOrderService(repo repository.CRUDRepository) *RepairOrderService {
    return &RepairOrderService{ResourceService: NewResourceService(repo, "repair_order", "customer_id","device_id","fault_description","service_method")}
}

func (s *RepairOrderService) Search(ctx context.Context, keyword string, page, pageSize int) (repository.Page, error) {
    return s.List(ctx, page, pageSize, keyword, "")
}

func (s *RepairOrderService) Active(ctx context.Context, page, pageSize int) (repository.Page, error) {
    return s.List(ctx, page, pageSize, "", "ACTIVE")
}

func (s *RepairOrderService) Archive(ctx context.Context, id int64) error {
    return s.Delete(ctx, id)
}

var repair_orderTransitions = map[string][]string{"PENDING":{"DISPATCHED","CANCELLED"},"DISPATCHED":{"WAITING_VISIT","WAITING_DELIVERY","ACCEPTED","CANCELLED"},"WAITING_VISIT":{"ACCEPTED","CANCELLED"},"WAITING_DELIVERY":{"ACCEPTED","CANCELLED"},"ACCEPTED":{"COMPLETED","CANCELLED"}}

func (s *RepairOrderService) ChangeStatus(ctx context.Context, id int64, current, next, reason string) (repository.Record, error) {
    if err := ValidateTransition(current, next, repair_orderTransitions); err != nil { return nil, err }
    return s.Update(ctx, id, repository.Record{"status": next, "status_reason": reason})
}