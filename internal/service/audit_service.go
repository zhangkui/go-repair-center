package service

import (
    "context"
    "go-repair-center/internal/repository"
)

type AuditService struct {
    *ResourceService
}

func NewAuditService(repo repository.CRUDRepository) *AuditService {
    return &AuditService{ResourceService: NewResourceService(repo, "audit", "action","resource_type")}
}

func (s *AuditService) Search(ctx context.Context, keyword string, page, pageSize int) (repository.Page, error) {
    return s.List(ctx, page, pageSize, keyword, "")
}

func (s *AuditService) Active(ctx context.Context, page, pageSize int) (repository.Page, error) {
    return s.List(ctx, page, pageSize, "", "ACTIVE")
}

func (s *AuditService) Archive(ctx context.Context, id int64) error {
    return s.Delete(ctx, id)
}
