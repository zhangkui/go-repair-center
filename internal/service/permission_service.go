package service

import (
    "context"
    "go-repair-center/internal/repository"
)

type PermissionService struct {
    *ResourceService
}

func NewPermissionService(repo repository.CRUDRepository) *PermissionService {
    return &PermissionService{ResourceService: NewResourceService(repo, "permission", "code","name","module","action")}
}

func (s *PermissionService) Search(ctx context.Context, keyword string, page, pageSize int) (repository.Page, error) {
    return s.List(ctx, page, pageSize, keyword, "")
}

func (s *PermissionService) Active(ctx context.Context, page, pageSize int) (repository.Page, error) {
    return s.List(ctx, page, pageSize, "", "ACTIVE")
}

func (s *PermissionService) Archive(ctx context.Context, id int64) error {
    return s.Delete(ctx, id)
}
