package service

import (
    "context"
    "go-repair-center/internal/repository"
)

type FaultCodeService struct {
    *ResourceService
}

func NewFaultCodeService(repo repository.CRUDRepository) *FaultCodeService {
    return &FaultCodeService{ResourceService: NewResourceService(repo, "fault_code", "code","name")}
}

func (s *FaultCodeService) Search(ctx context.Context, keyword string, page, pageSize int) (repository.Page, error) {
    return s.List(ctx, page, pageSize, keyword, "")
}

func (s *FaultCodeService) Active(ctx context.Context, page, pageSize int) (repository.Page, error) {
    return s.List(ctx, page, pageSize, "", "ACTIVE")
}

func (s *FaultCodeService) Archive(ctx context.Context, id int64) error {
    return s.Delete(ctx, id)
}
