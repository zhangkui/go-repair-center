package service

import (
    "context"
    "go-repair-center/internal/repository"
)

type SystemConfigService struct {
    *ResourceService
}

func NewSystemConfigService(repo repository.CRUDRepository) *SystemConfigService {
    return &SystemConfigService{ResourceService: NewResourceService(repo, "system_config", "config_key","config_value")}
}

func (s *SystemConfigService) Search(ctx context.Context, keyword string, page, pageSize int) (repository.Page, error) {
    return s.List(ctx, page, pageSize, keyword, "")
}

func (s *SystemConfigService) Active(ctx context.Context, page, pageSize int) (repository.Page, error) {
    return s.List(ctx, page, pageSize, "", "ACTIVE")
}

func (s *SystemConfigService) Archive(ctx context.Context, id int64) error {
    return s.Delete(ctx, id)
}
