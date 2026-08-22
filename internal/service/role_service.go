package service

import (
    "context"
    "go-repair-center/internal/repository"
)

type RoleService struct {
    *ResourceService
}

func NewRoleService(repo repository.CRUDRepository) *RoleService {
    return &RoleService{ResourceService: NewResourceService(repo, "role", "code","name")}
}

func (s *RoleService) Search(ctx context.Context, keyword string, page, pageSize int) (repository.Page, error) {
    return s.List(ctx, page, pageSize, keyword, "")
}

func (s *RoleService) Active(ctx context.Context, page, pageSize int) (repository.Page, error) {
    return s.List(ctx, page, pageSize, "", "ACTIVE")
}

func (s *RoleService) Archive(ctx context.Context, id int64) error {
    return s.Delete(ctx, id)
}
