package service

import (
    "context"
    "go-repair-center/internal/repository"
)

type RepairProcedureService struct {
    *ResourceService
}

func NewRepairProcedureService(repo repository.CRUDRepository) *RepairProcedureService {
    return &RepairProcedureService{ResourceService: NewResourceService(repo, "repair_procedure", "repair_execution_id","name")}
}

func (s *RepairProcedureService) Search(ctx context.Context, keyword string, page, pageSize int) (repository.Page, error) {
    return s.List(ctx, page, pageSize, keyword, "")
}

func (s *RepairProcedureService) Active(ctx context.Context, page, pageSize int) (repository.Page, error) {
    return s.List(ctx, page, pageSize, "", "ACTIVE")
}

func (s *RepairProcedureService) Archive(ctx context.Context, id int64) error {
    return s.Delete(ctx, id)
}
